package http

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"

	killgrave "github.com/friendsofgo/killgrave/internal"
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
		{"imposter directory not found", NewServer("failImposterPath", nil, &http.Server{}, &Proxy{}, false), errors.New("hello")},
		{"malformatted json", NewServer("test/testdata/malformatted_imposters", nil, &http.Server{}, &Proxy{}, false), nil},
		{"valid imposter", NewServer("test/testdata/imposters", mux.NewRouter(), &http.Server{}, &Proxy{}, false), nil},
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
		httpServer := http.Server{Handler: router}
		proxyServer, err := NewProxy(proxyServer.URL, mode)
		if err != nil {
			t.Fatal("NewProxy failed: ", err)
		}
		server := NewServer("test/testdata/imposters", router, &httpServer, proxyServer, false)
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
func TestBuildTLSMode(t *testing.T) {
	// proxyServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	io.WriteString(w, "Proxied")
	// }))
	// defer proxyServer.Close()
	// makeServer := func() (*Server, func()) {
	// 	router := mux.NewRouter()
	// 	httpServer := http.Server{Handler: router}
	// 	proxyServer, err := NewProxy(proxyServer.URL, killgrave.ProxyNone)
	// 	if err != nil {
	// 		t.Fatal("NewProxy failed: ", err)
	// 	}
	// 	server := NewServer("test/testdata/imposters_tls", router, &httpServer, proxyServer, true)
	// 	return &server, func() {
	// 		httpServer.Close()
	// 	}
	// }
	// testCases := map[string]struct {
	// 	url    string
	// 	body   string
	// 	status int
	// }{
	// 	"ProxyNone_Hit": {
	// 		url:    "/testRequest",
	// 		body:   "Handled",
	// 		status: http.StatusOK,
	// 	},
	// }
	// for name, tc := range testCases {
	// 	t.Run(name, func(t *testing.T) {
	// 		s, cleanUp := makeServer()
	// 		defer cleanUp()
	// 		s.Build()

	// 		client := http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	// 		// url := "https://" + proxyServer.Listener.Addr().String() + tc.url
	// 		response, err := client.Get(proxyServer.URL + tc.url)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		defer response.Body.Close()
	// 		body, _ := ioutil.ReadAll(response.Body)

	// 		if string(body) != tc.body {
	// 			t.Errorf("Expected body: %v, got: %s", tc.body, body)
	// 		}
	// 		if response.StatusCode != tc.status {
	// 			t.Errorf("Expected status code: %v, got: %v", tc.status, response.StatusCode)
	// 		}
	// 	})
	// }
}
