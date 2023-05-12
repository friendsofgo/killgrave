package acceptance

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/txtar"
)

const addr = "http://localhost:3000"

func Test(t *testing.T) {
	// For every test directory
	tcs := collectTestCases(t)

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Prepare the config directory
			path, clean := createTmpCfgDir(t, tc)
			t.Cleanup(clean)

			// Start the application
			stop := runApplication(t, path)
			t.Cleanup(stop)

			// For every request and response pair
			rrs := collectRequestResponses(t, tc.path)
			for _, rr := range rrs {
				rr := rr
				t.Run(rr.name, func(t *testing.T) {
					// Send the request
					res, err := http.DefaultClient.Do(rr.req)
					require.NoError(t, err)

					// Assert the res
					rr.assertResponse(res)
				})
			}
		})
	}
}

type tc struct {
	name string
	path string
}

func collectTestCases(t *testing.T) (cases []tc) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	testsDir := filepath.Join(cwd, "tests")
	entries, err := os.ReadDir(testsDir)
	require.NoError(t, err)

	for _, entry := range entries {
		if entry.IsDir() {
			cases = append(cases, tc{
				name: entry.Name(),
				path: filepath.Join(testsDir, entry.Name()),
			})
		}
	}

	return
}

type rr struct {
	*testing.T
	name string
	req  *http.Request
	res  []byte
}

func (rr rr) assertResponse(response *http.Response) {
	// Read the response body
	body, err := io.ReadAll(response.Body)
	require.NoError(rr, err)

	// Format the status line
	statusLine := fmt.Sprintf("HTTP/%d.%d %s", response.ProtoMajor, response.ProtoMinor, response.Status)

	// Strip dynamic headers (to prevent false negatives)
	response.Header.Del("Date")

	// Format the headers
	var headersBuilder strings.Builder
	err = response.Header.Write(&headersBuilder)
	require.NoError(rr, err)
	headers := strings.ReplaceAll(headersBuilder.String(), "\r\n", "\n")

	// Format the response
	res := fmt.Sprintf("%s\n%s\n\n%s", statusLine, headers, body)
	assert.Equal(rr, string(rr.res), res)
}

func collectRequestResponses(t *testing.T, path string) (rrs []rr) {
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

		rrs = append(rrs, rr{
			T:    t,
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

func readRequest(t *testing.T, raw []byte) *http.Request {
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(raw)))
	require.NoError(t, err)

	baseURL, err := url.Parse(addr)
	require.NoError(t, err)

	req.RequestURI = ""
	req.URL = baseURL.ResolveReference(req.URL)

	return req
}

func createTmpCfgDir(t *testing.T, tc tc) (string, func()) {
	tmpCfgDir := filepath.Join(os.TempDir(), tc.name)

	cfgFilePath := filepath.Join(tc.path, "config.txtar")

	contents, err := os.ReadFile(cfgFilePath)
	require.NoError(t, err)

	archive := txtar.Parse(contents)

	for _, f := range archive.Files {
		filePath := filepath.Join(tmpCfgDir, f.Name)
		fileDir := filepath.Dir(filePath)

		err := os.MkdirAll(fileDir, os.ModePerm)
		require.NoError(t, err)

		err = os.WriteFile(filePath, f.Data, os.ModePerm)
		require.NoError(t, err)
	}

	return tmpCfgDir, func() {
		err := os.RemoveAll(tmpCfgDir)
		require.NoError(t, err)
	}
}

func runApplication(t *testing.T, from string) func() {
	cmd := exec.Command("killgrave", "--imposters", filepath.Join(from, "imposters"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	require.NoError(t, err)

	// Trick to give time to the app to start
	time.Sleep(1 * time.Second)

	return func() {
		err := cmd.Process.Signal(os.Interrupt)
		require.NoError(t, err)

		err = cmd.Wait()
		require.NoError(t, err)
	}
}
