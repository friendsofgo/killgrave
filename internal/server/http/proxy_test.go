package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	killgrave "github.com/friendsofgo/killgrave/internal"
	"github.com/stretchr/testify/assert"
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
			proxy, err := NewProxy(tc.rawURL, tc.mode)
			if tc.err != nil {
				assert.NotNil(t, err)
				return
			} else {
				assert.Nil(t, err)
			}

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
	assert.Nil(t, err)

	frontend := httptest.NewServer(proxy.Handler())
	defer frontend.Close()

	_, err = http.Get(frontend.URL)
	assert.Nil(t, err)
	assert.True(t, isRequestHandled)

}
