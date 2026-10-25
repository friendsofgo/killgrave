//go:build acceptance

package acceptance

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
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
	require.NoError(t, err, "build Killgrave with make build before running acceptance tests")
	t.Logf("Running acceptance tests with Killgrave version: %s", version)

	// Once we know the Killgrave binary is available, proceed with the acceptance tests.
	// The first step is to collect all test cases from the `tests` directory. For each test:
	//
	// 0. The test case self-contain on each directory, which name is used as the test name.
	// 1. Requires a `config.txtar` file at the root level of the test case directory,
	//    it is used to initialize a file system with all the configuration-related files,
	//    which not only includes the Killgrave configuration file but also the imposters.
	// 2. Uses optional scenario.json startup settings, or the default imposter flags.
	//    Startup-error scenarios assert an unsuccessful exit instead of sending HTTP.
	//    Each server runs on an available port so scenarios can run concurrently.
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
			opts := readScenario(t, tc.path)
			port, err := network.AnyAvailablePort()
			require.NoError(t, err)
			var calls atomic.Int32
			upstreamURL := ""
			if opts.UpstreamCalls != nil {
				upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					calls.Add(1)
					body, err := io.ReadAll(r.Body)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					w.Header().Set("Content-Type", "text/plain")
					w.Header().Set("X-Upstream", "yes")
					w.WriteHeader(http.StatusAccepted)
					fmt.Fprintf(w, "%s %s\nX-Test: %s\n%s", r.Method, r.URL.RequestURI(), r.Header.Get("X-Test"), body)
				}))
				t.Cleanup(upstream.Close)
				upstreamURL = upstream.URL
				if opts.UpstreamDown {
					upstream.Close()
				}
			}
			replace := strings.NewReplacer("${PORT}", strconv.Itoa(port), "${UPSTREAM}", upstreamURL)
			path := createTmpCfgDir(t, tc, replace)
			if opts.OccupyPort {
				listener, err := net.Listen("tcp", net.JoinHostPort(addr, strconv.Itoa(port)))
				require.NoError(t, err)
				t.Cleanup(func() { listener.Close() })
			}

			// 2. Start the Killgrave process.
			args := []string{"-P", strconv.Itoa(port), "--imposters", filepath.Join(path, "imposters")}
			if opts.Args != nil {
				args = make([]string, len(opts.Args))
				for i, arg := range opts.Args {
					args[i] = replace.Replace(arg)
				}
			}
			if opts.Error != "" {
				assertStartupError(t, path, args, opts.Error)
				return
			}
			address := net.JoinHostPort(addr, strconv.Itoa(port))
			runApplication(t, path, address, args)
			client := &http.Client{Timeout: 3 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			}}

			// 3. Collect the request and response pairs
			// and iterate over them to perform the tests.
			rrs := collectRequestResponses(t, tc.path)
			for _, rr := range rrs {
				rr := rr
				t.Run(rr.name, func(t *testing.T) {
					// Override the address
					rr.overrideAddress(address)

					// Send the request
					res, err := client.Do(rr.req)
					require.NoError(t, err)
					t.Cleanup(func() { res.Body.Close() })

					// Assert the res
					rr.assertResponse(t, res)
				})
			}
			if opts.UpstreamCalls != nil {
				assert.Equal(t, *opts.UpstreamCalls, int(calls.Load()), "upstream request count")
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
func createTmpCfgDir(t *testing.T, tc testCase, replace *strings.Replacer) string {
	// First, we read the `config.txtar` file and initialize a txtar.Archive with its contents.
	tmpCfgDir := t.TempDir()
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

		err = os.WriteFile(filePath, []byte(replace.Replace(string(f.Data))), 0o600)
		require.NoError(t, err)
	}

	return tmpCfgDir
}

// scenario contains optional startup settings. Args are passed literally, with only
// ${PORT} and ${UPSTREAM} substituted; no shell is involved.
type scenario struct {
	Args          []string `json:"args"`
	Error         string   `json:"error"`
	OccupyPort    bool     `json:"occupyPort"`
	UpstreamCalls *int     `json:"upstreamCalls"`
	UpstreamDown  bool     `json:"upstreamDown"`
}

func readScenario(t *testing.T, path string) scenario {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(path, "scenario.json"))
	if errors.Is(err, os.ErrNotExist) {
		return scenario{}
	}
	require.NoError(t, err)
	var opts scenario
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	require.NoError(t, decoder.Decode(&opts))
	return opts
}

func applicationCommand(t *testing.T, dir string, args []string) *exec.Cmd {
	t.Helper()
	binary, err := filepath.Abs(bin)
	require.NoError(t, err)
	cmd := exec.Command(binary, args...)
	cmd.Dir = dir
	return cmd
}

func assertStartupError(t *testing.T, dir string, args []string, message string) {
	t.Helper()
	cmd := applicationCommand(t, dir, args)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd = exec.CommandContext(ctx, cmd.Path, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	require.NoError(t, ctx.Err(), "startup timed out: %s", output)
	var exitErr *exec.ExitError
	require.ErrorAs(t, err, &exitErr, "expected unsuccessful exit: %s", output)
	assert.Contains(t, string(output), message)
}

// runApplication captures output in a file so concurrent writes and diagnostics
// cannot race. TCP readiness does not consume responses or hit a proxy upstream.
func runApplication(t *testing.T, dir, address string, args []string) {
	t.Helper()
	cmd := applicationCommand(t, dir, args)
	outputPath := filepath.Join(dir, "process.log")
	output, err := os.Create(outputPath)
	require.NoError(t, err)
	t.Cleanup(func() {
		output.Close()
		if t.Failed() {
			data, _ := os.ReadFile(outputPath)
			t.Logf("Killgrave output:\n%s", data)
		}
	})
	cmd.Stdout, cmd.Stderr = output, output
	require.NoError(t, cmd.Start())
	done := make(chan struct{})
	var waitErr error
	go func() {
		waitErr = cmd.Wait()
		close(done)
	}()
	t.Cleanup(func() {
		select {
		case <-done:
		default:
			_ = cmd.Process.Signal(os.Interrupt)
			select {
			case <-done:
			case <-time.After(3 * time.Second):
				_ = cmd.Process.Kill()
				<-done
				t.Error("Killgrave did not stop within 3 seconds")
			}
		}
		assert.NoError(t, waitErr, "Killgrave process exit")
	})
	deadline := time.NewTimer(5 * time.Second)
	defer deadline.Stop()
	tick := time.NewTicker(20 * time.Millisecond)
	defer tick.Stop()
	for {
		select {
		case <-done:
			t.Fatalf("Killgrave exited before becoming ready: %v", waitErr)
		case <-deadline.C:
			t.Fatal("Killgrave did not become ready within 5 seconds")
		case <-tick.C:
			conn, err := net.DialTimeout("tcp", address, 100*time.Millisecond)
			if err == nil {
				conn.Close()
				return
			}
		}
	}
}

// extractKillgraveVersion runs `killgrave --version` and uses a regular expression
// to extract the version from the output. In case there's any error (e.g. the binary is not
// available), it returns the error.
func extractKillgraveVersion(t *testing.T) (string, error) {
	t.Helper()

	// Bound even the initial binary check so a broken command cannot hang the suite.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, "--version")

	// Capture the command's output.
	out := new(bytes.Buffer)
	cmd.Stdout = out

	// Run the command, and check for errors.
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Extract the Killgrave version from the output.
	re := regexp.MustCompile(`Killgrave version:\s*(\S+)`)
	match := re.FindStringSubmatch(out.String())
	if len(match) == 2 {
		return match[1], nil
	}

	return "", errors.New("version not found")
}
