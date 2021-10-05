package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	killgrave "github.com/friendsofgo/killgrave/internal"
)

func TestNewProxy(t *testing.T) {
	testCases := map[string]struct {
		rawURL string
		mode   killgrave.ProxyMode
		err    error
	}{
		"valid all":       {"all", killgrave.ProxyAll, nil},
		"valid mode none": {"none", killgrave.ProxyNone, nil},
		"error rawURL":    {":http!/gogle.com", killgrave.ProxyNone, errors.New("error")},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			proxy, err := NewProxy(tc.rawURL, "", tc.mode)
			if err != nil && tc.err == nil {
				t.Fatalf("not expected any erros and got %v", err)
			}

			if err == nil && tc.err != nil {
				t.Fatalf("expected an error and got nil")
			}
			if err != nil {
				return
			}
			if tc.mode != proxy.mode {
				t.Fatalf("expected: %v, got: %v", tc.mode, proxy.mode)
			}
		})
	}
}

func TestProxyHandler(t *testing.T) {
	isRequestHandled := false
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isRequestHandled = true
	}))
	defer backend.Close()

	proxy, err := NewProxy(backend.URL, "", killgrave.ProxyAll)
	if err != nil {
		t.Fatal("NewProxy failed: ", err)
	}

	frontend := httptest.NewServer(proxy.Handler())
	defer frontend.Close()

	_, err = http.Get(frontend.URL)
	if err != nil {
		t.Fatal("Frontend GET method failed: ", err)
	}
	if isRequestHandled != true {
		t.Fatal("Request was not proxied to backend server")
	}

}
