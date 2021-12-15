package http

import (
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	killgrave "github.com/friendsofgo/killgrave/internal"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestServer_Build(t *testing.T) {
	var serverData = []struct {
		name   string
		server Server
		err    error
	}{
		{"imposter directory not found", NewServer("failImposterPath", nil, &http.Server{}, &Proxy{}, false, nil, nil), errors.New("hello")},
		{"malformatted json", NewServer("test/testdata/malformatted_imposters", nil, &http.Server{}, &Proxy{}, false, nil, nil), nil},
		{"valid imposter", NewServer("test/testdata/imposters", mux.NewRouter(), &http.Server{}, &Proxy{}, false, nil, nil), nil},
	}

	for _, tt := range serverData {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.server.Build()

			if err == nil {
				if tt.err != nil {
					t.Fatalf("expected an error and got nil")
				}
			}

			if err != nil {
				if tt.err == nil {
					t.Fatalf("not expected any erros and got %+v", err)
				}
			}
		})
	}
}

func TestBuildProxyMode(t *testing.T) {
	proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Proxied")
	}))
	defer proxyServer.Close()
	makeServer := func(mode killgrave.ProxyMode) (*Server, func()) {
		router := mux.NewRouter()
		httpServer := &http.Server{Handler: router}
		proxyServer, err := NewProxy(proxyServer.URL, mode)
		if err != nil {
			t.Fatal("NewProxy failed: ", err)
		}
		server := NewServer("test/testdata/imposters", router, httpServer, proxyServer, false, nil, nil)
		return &server, func() {
			httpServer.Close()
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
			body, _ := ioutil.ReadAll(response.Body)

			if string(body) != tc.body {
				t.Errorf("Expected body: %v, got: %s", tc.body, body)
			}
			if response.StatusCode != tc.status {
				t.Errorf("Expected status code: %v, got: %v", tc.status, response.StatusCode)
			}
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
		if err != nil {
			t.Fatal("NewProxy failed: ", err)
		}
		server := NewServer("test/testdata/imposters_secure", router, httpServer, proxyServer, true, nil, nil)
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
			if err != nil {
				t.Fatalf("Non expected error trying to build server: %v", err)
			}
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

				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					return false
				}

				return string(body) == tc.body && response.StatusCode == tc.status
			}, 1*time.Second, 50*time.Millisecond)
		})
	}
}
