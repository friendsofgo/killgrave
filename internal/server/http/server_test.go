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
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestServer_Build(t *testing.T) {
	imposterFs := NewImposterFS(afero.NewOsFs())

	var serverData = []struct {
		name   string
		server Server
		err    error
	}{
		{"imposter directory not found", NewServer("failImposterPath", nil, &http.Server{}, &Proxy{}, false, imposterFs), errors.New("hello")},
		{"malformatted json", NewServer("test/testdata/malformatted_imposters", nil, &http.Server{}, &Proxy{}, false, imposterFs), nil},
		{"valid imposter", NewServer("test/testdata/imposters", mux.NewRouter(), &http.Server{}, &Proxy{}, false, imposterFs), nil},
	}

	for _, tt := range serverData {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.server.Build()

			if tt.err != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
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
		assert.Nil(t, err)
		imposterFs := NewImposterFS(afero.NewOsFs())
		server := NewServer("test/testdata/imposters", router, httpServer, proxyServer, false, imposterFs)
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
		assert.Nil(t, err)
		imposterFs := NewImposterFS(afero.NewOsFs())
		server := NewServer("test/testdata/imposters_secure", router, httpServer, proxyServer, true, imposterFs)
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

				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					return false
				}

				return string(body) == tc.body && response.StatusCode == tc.status
			}, 1*time.Second, 50*time.Millisecond)
		})
	}
}
