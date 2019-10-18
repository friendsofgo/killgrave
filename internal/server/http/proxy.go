package http

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	killgrave "github.com/friendsofgo/killgrave/internal"
)

// Proxy represent reverse proxy server.
type Proxy struct {
	server *httputil.ReverseProxy
	mode   killgrave.ProxyMode
}

// NewProxy creates new proxy server.
func NewProxy(rawurl string, mode killgrave.ProxyMode) (*Proxy, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	reverseProxy := httputil.NewSingleHostReverseProxy(u)
	return &Proxy{reverseProxy, mode}, nil
}

// Handler returns handler that sends request to another server.
func (p *Proxy) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p.server.ServeHTTP(w, r)
	}
}
