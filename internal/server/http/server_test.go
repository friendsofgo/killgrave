package http

import (
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	killgrave "github.com/friendsofgo/killgrave/internal"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
	os.Exit(m.Run())
}

func TestServer_Build(t *testing.T) {
	newServer := func(fs ImposterFs) Server {
		return NewServer(mux.NewRouter(), &http.Server{}, &Proxy{}, false, fs, nil, 0, "")
	}

	testCases := map[string]struct {
		impostersPath string
		shouldFail    bool
	}{
		"imposters with malformed json": {impostersPath: "test/testdata/malformed_imposters"},
		"valid imposters":               {impostersPath: "test/testdata/imposters"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			fs, err := NewImposterFS(tc.impostersPath)
			require.NoError(t, err)

			srv := newServer(fs)
			err = srv.Build()

			if tc.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBuildProxyMode(t *testing.T) {
	proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Proxied")
	}))
	defer proxyServer.Close()

	makeServer := func(mode killgrave.ProxyMode) (*Server, func() error) {
		router := mux.NewRouter()
		httpServer := &http.Server{Handler: router}

		proxyServer, err := NewProxy(proxyServer.URL, mode)
		require.NoError(t, err)

		imposterFs, err := NewImposterFS("test/testdata/imposters")
		require.NoError(t, err)

		server := NewServer(router, httpServer, proxyServer, false, imposterFs, nil, 0, "")
		return &server, func() error {
			return httpServer.Close()
		}
	}

	testCases := map[string]struct {
		mode   killgrave.ProxyMode
		url    string
		body   string
		status int
	}{
		"ProxyAll": {
			mode:   killgrave.ProxyAll,
			url:    "/testRequest",
			body:   "Proxied",
			status: http.StatusOK,
		},
		"ProxyMissing_Hit": {
			mode:   killgrave.ProxyMissing,
			url:    "/testRequest",
			body:   "Handled",
			status: http.StatusOK,
		},
		"ProxyMissing_Proxied": {
			mode:   killgrave.ProxyMissing,
			url:    "/NonExistentURL123",
			body:   "Proxied",
			status: http.StatusOK,
		},
		"ProxyNone_Hit": {
			mode:   killgrave.ProxyNone,
			url:    "/testRequest",
			body:   "Handled",
			status: http.StatusOK,
		},
		"ProxyNone_Missing": {
			mode:   killgrave.ProxyNone,
			url:    "/NonExistentURL123",
			body:   "404 page not found\n",
			status: http.StatusNotFound,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			s, cleanUp := makeServer(tc.mode)
			defer cleanUp()
			s.Build()

			req := httptest.NewRequest("GET", tc.url, nil)
			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)
			response := w.Result()
			body, _ := io.ReadAll(response.Body)

			assert.Equal(t, tc.body, string(body))
			assert.Equal(t, tc.status, response.StatusCode)
		})
	}
}

func TestBuildSecureMode(t *testing.T) {
	proxyServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Proxied")
	}))
	defer proxyServer.Close()

	makeServer := func(mode killgrave.ProxyMode) (*Server, func()) {
		router := mux.NewRouter()
		cert, _ := tls.X509KeyPair(serverCert, serverKey)
		httpServer := &http.Server{Handler: router, Addr: ":4430", TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		}}

		proxyServer, err := NewProxy(proxyServer.URL, mode)
		require.NoError(t, err)

		imposterFs, err := NewImposterFS("test/testdata/imposters_secure")
		require.NoError(t, err)

		server := NewServer(router, httpServer, proxyServer, true, imposterFs, nil, 0, "")
		return &server, func() {
			httpServer.Close()
		}
	}

	testCases := map[string]struct {
		mode   killgrave.ProxyMode
		url    string
		body   string
		status int
		server *httptest.Server
	}{
		"ProxyNone_Hit": {
			mode:   killgrave.ProxyNone,
			url:    "https://localhost:4430/testHTTPSRequest",
			body:   "Handled",
			status: http.StatusOK,
			server: proxyServer,
		},
		"ProxyAlways_Hit": {
			mode:   killgrave.ProxyAll,
			url:    proxyServer.URL,
			body:   "Proxied",
			status: http.StatusOK,
			server: proxyServer,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			s, cleanUp := makeServer(tc.mode)
			defer cleanUp()

			err := s.Build()
			assert.Nil(t, err)
			s.Run()

			client := tc.server.Client()
			client.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}

			assert.Eventually(t, func() bool {
				response, err := client.Get(tc.url)
				if err != nil {
					return false
				}

				defer response.Body.Close()

				body, err := io.ReadAll(response.Body)
				if err != nil {
					return false
				}

				return string(body) == tc.body && response.StatusCode == tc.status
			}, 1*time.Second, 50*time.Millisecond)
		})
	}
}

func TestBuildLogRequests(t *testing.T) {
	testCases := map[string]struct {
		method         string
		path           string
		contentType    string
		body           string
		logLevel       int
		expectedLog    string
		expectedStatus int
	}{
		"GET valid imposter request": {
			method:         "GET",
			path:           "/yamlTestDumpRequest",
			contentType:    "text/plain",
			body:           "Dumped",
			logLevel:       1,
			expectedLog:    "GET /yamlTestDumpRequest HTTP/1.1\" 200 17\n",
			expectedStatus: http.StatusOK,
		},
		"GET valid imposter request with body": {
			method:         "GET",
			path:           "/yamlTestDumpRequest",
			contentType:    "text/plain",
			body:           "Dumped",
			logLevel:       2,
			expectedLog:    "GET /yamlTestDumpRequest HTTP/1.1\" 200 17 Dumped\n",
			expectedStatus: http.StatusOK,
		},
		"GET valid imposter binary request": {
			method:         "GET",
			path:           "/yamlTestDumpRequest",
			contentType:    "application/octet-stream",
			body:           "Dumped",
			logLevel:       1,
			expectedLog:    "GET /yamlTestDumpRequest HTTP/1.1\" 200 17\n",
			expectedStatus: http.StatusOK,
		},
		"GET valid imposter binary request with body": {
			method:         "GET",
			path:           "/yamlTestDumpRequest",
			contentType:    "application/octet-stream",
			body:           "Dumped",
			logLevel:       2,
			expectedLog:    "GET /yamlTestDumpRequest HTTP/1.1\" 200 17 RHVtcGVk\n",
			expectedStatus: http.StatusOK,
		},
		"GET valid imposter request no body": {
			method:         "GET",
			path:           "/yamlTestDumpRequest",
			contentType:    "text/plain",
			body:           "",
			logLevel:       2,
			expectedLog:    "GET /yamlTestDumpRequest HTTP/1.1\" 200 17\n",
			expectedStatus: http.StatusOK,
		},
		"GET invalid imposter request": {
			method:         "GET",
			path:           "/doesnotexist",
			contentType:    "text/plain",
			body:           "Dumped",
			logLevel:       1,
			expectedLog:    "GET /doesnotexist HTTP/1.1\" 404 19\n",
			expectedStatus: http.StatusNotFound,
		},
		"GET invalid imposter request with body": {
			method:         "GET",
			path:           "/doesnotexist",
			contentType:    "text/plain",
			body:           "Dumped",
			logLevel:       2,
			expectedLog:    "GET /doesnotexist HTTP/1.1\" 404 19 Dumped\n",
			expectedStatus: http.StatusNotFound,
		},
		"GET invalid imposter binary request with body": {
			method:         "GET",
			path:           "/doesnotexist",
			contentType:    "video/mp4",
			body:           "Dumped",
			logLevel:       2,
			expectedLog:    "GET /doesnotexist HTTP/1.1\" 404 19 RHVtcGVk\n",
			expectedStatus: http.StatusNotFound,
		},
		"GET invalid imposter request no body": {
			method:         "GET",
			path:           "/doesnotexist",
			contentType:    "text/plain",
			body:           "",
			logLevel:       2,
			expectedLog:    "GET /doesnotexist HTTP/1.1\" 404 19\n",
			expectedStatus: http.StatusNotFound,
		},
	}
	for name, tc := range testCases {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer func() {
				log.SetOutput(os.Stderr)
			}()

			imposterFs, err := NewImposterFS("test/testdata/imposters")
			assert.NoError(t, err)
			server := NewServer(mux.NewRouter(), &http.Server{}, &Proxy{}, false, imposterFs, nil, tc.logLevel, "")
			err = server.Build()
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", tc.contentType)

			server.httpServer.Handler.ServeHTTP(w, req)

			response := w.Result()
			assert.Equal(t, tc.expectedStatus, response.StatusCode, "Expected status code: %v, got: %v", tc.expectedStatus, response.StatusCode)

			// verify the request is dumped in the logs
			assert.Contains(t, buf.String(), tc.expectedLog, "Expect request dumped on logs failed")
		})
	}
}

func TestBuildRecordRequests(t *testing.T) {

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	tempDir, err := os.MkdirTemp("", "testdir")
	assert.NoError(t, err, "Failed to create temporary directory")
	defer os.RemoveAll(tempDir)
	dumpFile := filepath.Join(tempDir, "dump_requests.log")

	imposterFs, err := NewImposterFS("test/testdata/imposters")
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	server := NewServer(mux.NewRouter(), &http.Server{}, &Proxy{}, false, imposterFs, nil, 0, dumpFile)
	err = server.Build()
	assert.NoError(t, err)

	expectedBodies := []string{"Dumped1", ""}
	req1 := httptest.NewRequest("GET", "/yamlTestDumpRequest", strings.NewReader(expectedBodies[0]))
	req2 := httptest.NewRequest("GET", "/yamlTestDumpRequest", strings.NewReader(expectedBodies[1]))
	server.httpServer.Handler.ServeHTTP(w, req1)
	server.httpServer.Handler.ServeHTTP(w, req2)

	// wait for channel to print out the requests
	time.Sleep(1 * time.Second)

	// check recoreded request dumps
	reqs, err := getRecordedRequests(dumpFile)
	assert.NoError(t, err, "Failed to read requests from file")
	assert.Equal(t, 2, len(reqs), "Expect 2 requests to be dumped in file failed")
	for i, expectedBody := range expectedBodies {
		assert.Equal(t, expectedBody, reqs[i].Body, "Expect request body to be dumped in file failed")
	}
}
