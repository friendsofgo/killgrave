package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	killgrave "github.com/friendsofgo/killgrave/internal"
	"github.com/stretchr/testify/assert"
)

func TestNewProxy(t *testing.T) {
	testCases := map[string]struct {
		rawURL  string
		mode    killgrave.ProxyMode
		wantErr bool
	}{
		"valid all":       {"all", killgrave.ProxyAll, false},
		"valid mode none": {"none", killgrave.ProxyNone, false},
		"error rawURL":    {":http!/gogle.com", killgrave.ProxyNone, true},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			proxy, err := NewProxy(tc.rawURL, tc.mode)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.mode, proxy.mode)
		})
	}
}

func TestProxyHandler(t *testing.T) {
	isRequestHandled := false
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isRequestHandled = true
	}))
	defer backend.Close()

	proxy, err := NewProxy(backend.URL, killgrave.ProxyAll)
	assert.NoError(t, err)

	frontend := httptest.NewServer(proxy.Handler())
	defer frontend.Close()

	_, err = http.Get(frontend.URL)
	assert.NoError(t, err)
	assert.True(t, isRequestHandled)

}
