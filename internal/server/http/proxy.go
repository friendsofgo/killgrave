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
	url    *url.URL
}

// NewProxy creates new proxy server.
func NewProxy(rawurl string, mode killgrave.ProxyMode) (*Proxy, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	reverseProxy := httputil.NewSingleHostReverseProxy(u)
	return &Proxy{server: reverseProxy, mode: mode, url: u}, nil
}

// Handler returns handler that sends request to another server.
func (p *Proxy) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Host = p.url.Host
		r.URL.Scheme = p.url.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = p.url.Host

		p.server.ServeHTTP(w, r)
	}
}
