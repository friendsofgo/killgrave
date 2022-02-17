package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	killgrave "github.com/friendsofgo/killgrave/internal"
)

// Proxy represent reverse proxy server.
type Proxy struct {
	server        *httputil.ReverseProxy
	mode          killgrave.ProxyMode
	url           *url.URL
	impostersPath string

	recorder RecorderHTTP
}

// NewProxy creates new proxy server.
func NewProxy(rawurl, impostersPath string, mode killgrave.ProxyMode, recorder RecorderHTTP) (*Proxy, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	reverseProxy := httputil.NewSingleHostReverseProxy(u)
	return &Proxy{
		server:   reverseProxy,
		mode:     mode,
		url:      u,
		recorder: recorder}, nil
}

// Handler returns handler that sends request to another server.
func (p Proxy) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Host = p.url.Host
		r.URL.Scheme = p.url.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = p.url.Host
		if p.mode == killgrave.ProxyRecord {
			p.server.ModifyResponse = p.recordProxy
		}

		p.server.ServeHTTP(w, r)
	}
}

func (p Proxy) recordProxy(resp *http.Response) error {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyStr := string(b)
	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))

	responseRecorder := ResponseRecorder{
		Headers: resp.Header,
		Status:  resp.StatusCode,
		Body:    bodyStr,
	}

	return p.recorder.Record(resp.Request, responseRecorder)
}
