package http

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	killgrave "github.com/friendsofgo/killgrave/internal"
)

// Proxy represent reverse proxy server.
type Proxy struct {
	server *httputil.ReverseProxy
	mode   killgrave.ProxyMode
	url    *url.URL
	impostersPath string
}

var ErrImpostersPathEmpty = errors.New("if you want to record the missing request you will need to indicate an imposters path")

// NewProxy creates new proxy server.
func NewProxy(rawurl, impostersPath string, mode killgrave.ProxyMode) (*Proxy, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	reverseProxy := httputil.NewSingleHostReverseProxy(u)
	if mode == killgrave.ProxyRecord {
		if impostersPath == "" {
			return nil, ErrImpostersPathEmpty
		}
		reverseProxy.ModifyResponse = recordProxy
	}
	return &Proxy{server: reverseProxy, mode: mode, url: u}, nil
}

// Handler returns handler that sends request to another server.
func (p Proxy) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Host = p.url.Host
		r.URL.Scheme = p.url.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = p.url.Host

		p.server.ServeHTTP(w, r)
	}
}

func recordProxy(resp *http.Response) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))
	return nil
}