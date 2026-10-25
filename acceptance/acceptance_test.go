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
	// Verify the built binary before starting the scenarios.
	version, err := extractKillgraveVersion(t)
	require.NoError(t, err, "build Killgrave with make build before running acceptance tests")
	t.Logf("Running acceptance tests with Killgrave version: %s", version)

	// Each archive contains setup files, optional startup settings, and named
	// HTTP steps. Scenarios run concurrently; steps within a scenario run in order.
	tcs := collectTestCases(t)

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// 1. Create a temporary directory with the configuration files.
			opts := readScenario(t, tc.files)
			port := availablePort(t)
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
			rrs := collectRequestResponses(t, tc.files)
			for _, rr := range rrs {
				rr := rr
				t.Run(rr.name, func(t *testing.T) {
					// Override the address
					rr.overrideAddress(address)

					rr.run(t, client, path)
				})
			}
			if opts.UpstreamCalls != nil {
				assert.Equal(t, *opts.UpstreamCalls, int(calls.Load()), "upstream request count")
			}
		})
	}
}

// availablePort asks the OS for a port, then releases it for the child process.
func availablePort(t *testing.T) int {
	t.Helper()
	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	port := listener.Addr().(*net.TCPAddr).Port
	require.NoError(t, listener.Close())
	return port
}

// testCase holds one scenario archive: setup files, startup options, and HTTP steps.
type testCase struct {
	name  string
	files []txtar.File
}

// collectTestCases reads one archive per scenario from the tests directory.
func collectTestCases(t *testing.T) []testCase {
	t.Helper()
	var tcs []testCase

	entries, err := os.ReadDir(testsDir)
	require.NoError(t, err)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txtar") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(testsDir, entry.Name()))
		require.NoError(t, err)
		tcs = append(tcs, testCase{
			name:  strings.TrimSuffix(entry.Name(), ".txtar"),
			files: txtar.Parse(data).Files,
		})
	}

	require.NotEmpty(t, tcs, "no acceptance scenarios found")
	return tcs
}

// reqRes is a data structure that holds the information required to run
// each of the test cases that compose each acceptance test:
// - the name of the test case.
// - the request to be sent to Killgrave, as *http.Request.
// - the expected response from Killgrave, as []byte.
type reqRes struct {
	name       string
	req        *http.Request
	res        []byte
	writes     []txtar.File
	eventually bool
}

// overrideAddress changes the request's URL to use the provided address.
// This is useful to run the tests against different addresses, e.g. different ports,
// which is a requirement to be able to run the acceptance tests concurrently.
func (rr reqRes) overrideAddress(address string) {
	rr.req.URL.Scheme = "http"
	rr.req.URL.Host = address
}

// run applies file changes once, then checks the response. Only stateless GET
// exchanges may opt into retries while waiting for a watcher reload.
func (rr reqRes) run(t *testing.T, client *http.Client, dir string) {
	t.Helper()
	for _, file := range rr.writes {
		require.True(t, filepath.IsLocal(file.Name), "write path must be relative: %s", file.Name)
		target := filepath.Join(dir, file.Name)
		require.NoError(t, os.MkdirAll(filepath.Dir(target), 0o755))
		// Publish a complete file so the watcher cannot load partially written JSON.
		tmp, err := os.CreateTemp(filepath.Dir(target), ".acceptance-*")
		require.NoError(t, err)
		_, writeErr := tmp.Write(file.Data)
		closeErr := tmp.Close()
		require.NoError(t, writeErr)
		require.NoError(t, closeErr)
		require.NoError(t, os.Rename(tmp.Name(), target))
	}
	if rr.eventually {
		require.Equal(t, http.MethodGet, rr.req.Method, "eventually is only for stateless GET requests")
		require.True(t, rr.req.Body == nil || rr.req.Body == http.NoBody, "eventually requests cannot have a body")
		require.EventuallyWithT(t, func(c *assert.CollectT) {
			res, err := client.Do(rr.req.Clone(context.Background()))
			if !assert.NoError(c, err) {
				return
			}
			actual, err := formatResponse(res)
			if assert.NoError(c, err) {
				assert.Equal(c, string(rr.res), actual)
			}
		}, 5*time.Second, 50*time.Millisecond)
		return
	}
	res, err := client.Do(rr.req)
	require.NoError(t, err)
	actual, err := formatResponse(res)
	require.NoError(t, err)
	assert.Equal(t, string(rr.res), actual)
}

// formatResponse consumes and closes the body on every attempt, including retries.
func formatResponse(response *http.Response) (string, error) {
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// Format the status line
	statusLine := fmt.Sprintf("HTTP/%d.%d %s", response.ProtoMajor, response.ProtoMinor, response.Status)

	// Strip dynamic headers (to prevent false negatives)
	response.Header.Del("Date")

	// Format the headers
	var headersBuilder strings.Builder
	err = response.Header.Write(&headersBuilder)
	if err != nil {
		return "", err
	}
	headers := strings.ReplaceAll(headersBuilder.String(), "\r\n", "\n")

	// Format the response
	return fmt.Sprintf("%s\n%s\n\n%s", statusLine, headers, body), nil
}

// collectRequestResponses groups http/<step>/ entries, preserving archive order.
func collectRequestResponses(t *testing.T, files []txtar.File) (rrs []reqRes) {
	t.Helper()
	var order []string
	steps := make(map[string][]txtar.File)
	for _, file := range files {
		rel, ok := strings.CutPrefix(file.Name, "http/")
		if !ok {
			continue
		}
		name, part, ok := strings.Cut(rel, "/")
		require.True(t, ok && name != "", "expected http/<step>/<entry>: %s", file.Name)
		if _, exists := steps[name]; !exists {
			order = append(order, name)
		}
		steps[name] = append(steps[name], txtar.File{Name: part, Data: file.Data})
	}
	require.NotEmpty(t, order, "scenario has no HTTP exchanges")
	for _, name := range order {
		entries := steps[name]
		require.NotEmpty(t, find(entries, "req.http"), "missing request for %s", name)
		require.NotEmpty(t, find(entries, "res.http"), "missing response for %s", name)
		rr := reqRes{
			name: name,
			req:  readRequest(t, find(entries, "req.http")),
			res:  find(entries, "res.http"),
		}
		for _, file := range entries {
			if name, ok := strings.CutPrefix(file.Name, "write/"); ok {
				rr.writes = append(rr.writes, txtar.File{Name: name, Data: file.Data})
			}
			if file.Name == "eventually" {
				eventually, err := strconv.ParseBool(strings.TrimSpace(string(file.Data)))
				require.NoError(t, err)
				rr.eventually = eventually
			}
		}
		rrs = append(rrs, rr)
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

// createTmpCfgDir extracts setup files into the process working directory.
// Scenario metadata and HTTP steps stay in the harness.
func createTmpCfgDir(t *testing.T, tc testCase, replace *strings.Replacer) string {
	t.Helper()
	tmpCfgDir := t.TempDir()

	for _, f := range tc.files {
		if f.Name == "scenario.json" || strings.HasPrefix(f.Name, "http/") {
			continue
		}
		require.True(t, filepath.IsLocal(f.Name), "setup path must be relative: %s", f.Name)
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

func readScenario(t *testing.T, files []txtar.File) scenario {
	t.Helper()
	data := find(files, "scenario.json")
	if data == nil {
		return scenario{}
	}
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
