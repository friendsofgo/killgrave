//go:build acceptance

package acceptance

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/txtar"

	"github.com/friendsofgo/killgrave/acceptance/utils/network"
)

const (
	// testsDir is the directory, within the `acceptance` folder,
	// where the acceptance tests are located at.
	testsDir = "tests"

	// addr is the address where the Killgrave binary
	addr = "localhost"

	// bin is the path to where the Killgrave binary
	// is expected to run acceptance tests.
	bin = "../bin/killgrave"
)

// Test is the entry point for the acceptance tests.
func Test(t *testing.T) {
	// First of all, we extract the Killgrave version by running `killgrave version`.
	// This is useful, not only to write a log that can serve as metadata for the test results,
	// but also to ensure that the Killgrave binary is available.
	version, err := extractKillgraveVersion(t)
	if err != nil {
		var pathErr *fs.PathError
		if errors.As(err, &pathErr) || strings.Contains(err.Error(), "executable file not found in $PATH") {
			log.Fatalf("Attention! It looks like you haven't compiled Killgrave, the execution of the Killgrave "+
				"binary has failed with: %v", err)
		}
		log.Fatalf("The execution of `killgrave version` has failed with: %v", err)
	}
	t.Logf("Running acceptance tests with Killgrave version: %s", version)

	// Once we now that the Killgrave binary is available, we can proceed with the acceptance tests.
	// The first step is to collect all test cases from the `tests` directory. For each test:
	//
	// 0. The test case self-contain on each directory, which name is used as the test name.
	// 1. Requires a `config.txtar` file at the root level of the test case directory,
	//    it is used to initialize a file system with all the configuration-related files,
	//    which not only includes the Killgrave configuration file but also the imposters.
	// 2. Runs Killgrave in any available port (so we can run multiple test cases at the same
	// 	  time), using the configuration files from the previous step.
	// 3. Requires an `http` directory which contains a set of request and response pairs.
	//    Each pair is defined by two files: `req.http` and `res.http`, where the request
	//    is the HTTP request that will be performed as part of the test, and the response
	//    is the response expected from Killgrave (which will be asserted).
	//    Each pair is considered one of the test cases that compose the acceptance test,
	//    defined by the aforementioned parent directory.
	tcs := collectTestCases(t)

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// 1. Create a temporary directory with the configuration files.
			path := createTmpCfgDir(t, tc)

			// 2. Start the Killgrave process.
			address := runApplication(t, path)

			// 3. Collect the request and response pairs
			// and iterate over them to perform the tests.
			rrs := collectRequestResponses(t, tc.path)
			for _, rr := range rrs {
				rr := rr
				t.Run(rr.name, func(t *testing.T) {
					// Override the address
					rr.overrideAddress(address)

					// Send the request
					res, err := http.DefaultClient.Do(rr.req)
					require.NoError(t, err)

					// Assert the res
					rr.assertResponse(t, res)
				})
			}
		})
	}
}

// testCase represents a test case to be run.
// It is defined by the name of the test case and the path
// to the directory where the testCase files live in.
type testCase struct {
	name string
	path string
}

// collectTestCases walks over the `tests` directory and
// constructs all the testCase's from the directories found.
func collectTestCases(t *testing.T) []testCase {
	var tcs []testCase

	cwd, err := os.Getwd()
	require.NoError(t, err)

	testsDir := filepath.Join(cwd, testsDir)
	entries, err := os.ReadDir(testsDir)
	require.NoError(t, err)

	for _, entry := range entries {
		if entry.IsDir() {
			tcs = append(tcs, testCase{
				name: entry.Name(),
				path: filepath.Join(testsDir, entry.Name()),
			})
		}
	}

	return tcs
}

// reqRes is a data structure that holds the information required to run
// each of the test cases that compose each acceptance test:
// - the name of the test case.
// - the request to be sent to Killgrave, as *http.Request.
// - the expected response from Killgrave, as []byte.
type reqRes struct {
	name string
	req  *http.Request
	res  []byte
}

// overrideAddress changes the request's URL to use the provided address.
// This is useful to run the tests against different addresses, e.g. different ports,
// which is a requirement to be able to run the acceptance tests concurrently.
func (rr reqRes) overrideAddress(address string) {
	rr.req.URL.Scheme = "http"
	rr.req.URL.Host = address
}

// assertResponse is a self-contained function that can be used to assert
// that the response received from Killgrave matches the expected response.
//
// It builds the response string from the response object, and then it
// compares it with the expected response (from the test definition).
func (rr reqRes) assertResponse(t *testing.T, response *http.Response) {
	t.Helper()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	// Format the status line
	statusLine := fmt.Sprintf("HTTP/%d.%d %s", response.ProtoMajor, response.ProtoMinor, response.Status)

	// Strip dynamic headers (to prevent false negatives)
	response.Header.Del("Date")

	// Format the headers
	var headersBuilder strings.Builder
	err = response.Header.Write(&headersBuilder)
	require.NoError(t, err)
	headers := strings.ReplaceAll(headersBuilder.String(), "\r\n", "\n")

	// Format the response
	res := fmt.Sprintf("%s\n%s\n\n%s", statusLine, headers, body)
	assert.Equal(t, string(rr.res), res)
}

// collectRequestResponses walks over the `http` directory of the test case
// and collects all the reqRes pairs. In other words, it collects all the test cases.
func collectRequestResponses(t *testing.T, path string) (rrs []reqRes) {
	httpDir := filepath.Join(path, "http")
	entries, err := os.ReadDir(httpDir)
	require.NoError(t, err)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		rrFilePath := filepath.Join(httpDir, entry.Name())

		contents, err := os.ReadFile(rrFilePath)
		require.NoError(t, err)

		archive := txtar.Parse(contents)

		rrs = append(rrs, reqRes{
			name: entry.Name(),
			req:  readRequest(t, find(archive.Files, "req.http")),
			res:  find(archive.Files, "res.http"),
		})
	}

	return
}

func find(ff []txtar.File, name string) []byte {
	for _, f := range ff {
		if f.Name == name {
			return f.Data
		}
	}
	return nil
}

// readRequest reads a raw HTTP request from a []byte (e.g. read from a file),
// and instantiates the equivalent *http.Request object from it.
func readRequest(t *testing.T, raw []byte) *http.Request {
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(raw)))
	require.NoError(t, err)

	baseURL, err := url.Parse(addr)
	require.NoError(t, err)

	req.RequestURI = ""
	req.URL = baseURL.ResolveReference(req.URL)

	return req
}

// createTmpCfgDir creates a temporary directory with the configuration files defined
// in the `config.txtar` file, replicating the structure to be used by Killgrave, when executed.
//
// We follow this approach because this way we can test the application assuming the binary
// already exists, so these tests can be run with a recently generated binary (e.g. a release).
//
// Additionally, in the future we might explore ways to reuse this setup to run these tests
// as "integration tests", so faking using a fake, likely in-memory, file system but directly
// calling app.Run().
func createTmpCfgDir(t *testing.T, tc testCase) string {
	// First, we read the `config.txtar` file and initialize a txtar.Archive with its contents.
	tmpCfgDir := filepath.Join(os.TempDir(), tc.name)
	cfgFilePath := filepath.Join(tc.path, "config.txtar")
	contents, err := os.ReadFile(cfgFilePath)
	require.NoError(t, err)
	archive := txtar.Parse(contents)

	// Then, we create the temporary directory and write the files.
	for _, f := range archive.Files {
		filePath := filepath.Join(tmpCfgDir, f.Name)
		fileDir := filepath.Dir(filePath)

		err := os.MkdirAll(fileDir, os.ModePerm)
		require.NoError(t, err)

		err = os.WriteFile(filePath, f.Data, os.ModePerm)
		require.NoError(t, err)
	}

	// Tell the testing framework to clean up the temporary directory after the test is done.
	t.Cleanup(func() {
		err := os.RemoveAll(tmpCfgDir)
		require.NoError(t, err)
	})

	return tmpCfgDir
}

// runApplication runs Killgrave assuming the binary already exists.
// It uses the imposters located at `from`, which path must be absolute.
// It runs the application on any available port, so we can run multiple
// tests concurrently. It returns the address as the first return value.
//
// For now, it redirects the application's output (stdout and stderr)
// to the test's output, but in the future, we might want to capture
// the output to assert the logs, and or use it in a smarter way.
func runApplication(t *testing.T, from string) string {
	// Look for any available port.
	port, err := network.AnyAvailablePort()
	address := addr + ":" + strconv.Itoa(port)
	require.NoError(t, err, "failed to find an available port")

	// Prepare the `killgrave` command, and start it.
	cmd := exec.Command(bin, "-P", strconv.Itoa(port), "--imposters", filepath.Join(from, "imposters"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	require.NoError(t, err)

	// Tell the testing framework to stop the process after the test is done.
	t.Cleanup(func() {
		err := cmd.Process.Signal(os.Interrupt)
		require.NoError(t, err)

		err = cmd.Wait()
		require.NoError(t, err)
	})

	// Wait for the application to be ready.
	const (
		maxWaitTime = 2 * time.Second
		checkEvery  = 100 * time.Millisecond
	)
	require.Eventually(t, func() bool {
		res, err := http.Get("http://" + address + "/nonExistingEndpoint")
		return err == nil && res != nil && res.StatusCode == http.StatusNotFound
	}, maxWaitTime, checkEvery)

	return address
}

// extractKillgraveVersion runs the `killgrave version` command and uses a regular expression
// to extract the version from the output. In case there's any error (e.g. the binary is not
// available), it returns the error.
func extractKillgraveVersion(t *testing.T) (string, error) {
	t.Helper()

	// Prepare the `killgrave version` command.
	cmd := exec.Command(bin, "version", "-v")

	// Capture the command's output.
	out := new(bytes.Buffer)
	cmd.Stdout = out

	// Run the command, and check for errors.
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Extract the Killgrave version from the output.
	re := regexp.MustCompile(`Killgrave version:\s*([a-zA-Z0-9\-]+)`)
	match := re.FindStringSubmatch(out.String())
	if len(match) == 2 {
		return match[1], nil
	}

	return "", errors.New("version not found")
}
