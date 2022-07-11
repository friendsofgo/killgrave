package http

import (
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"

	killgrave "github.com/friendsofgo/killgrave/internal"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//go:embed cert/server.key
var serverKey []byte

//go:embed cert/server.cert
var serverCert []byte

var (
	defaultCORSMethods        = []string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE", "PATCH", "TRACE", "CONNECT"}
	defaultCORSHeaders        = []string{"X-Requested-With", "Content-Type", "Authorization"}
	defaultCORSExposedHeaders = []string{"Cache-Control", "Content-Language", "Content-Type", "Expires", "Last-Modified", "Pragma"}
)

// ServerOpt function that allow modify the current server
type ServerOpt func(s *Server)

// Server definition of mock server
type Server struct {
	impostersPath string
	router        *mux.Router
	httpServer    *http.Server
	proxy         *Proxy
	secure        bool
	imposterFs    ImposterFs
}

// NewServer initialize the mock server
func NewServer(p string, r *mux.Router, httpServer *http.Server, proxyServer *Proxy, secure bool, fs ImposterFs) Server {
	return Server{
		impostersPath: p,
		router:        r,
		httpServer:    httpServer,
		proxy:         proxyServer,
		secure:        secure,
		imposterFs:    fs,
	}
}

// PrepareAccessControl Return options to initialize the mock server with default access control
func PrepareAccessControl(config killgrave.ConfigCORS) (h []handlers.CORSOption) {
	h = append(h, handlers.AllowedMethods(defaultCORSMethods))
	h = append(h, handlers.AllowedHeaders(defaultCORSHeaders))
	h = append(h, handlers.ExposedHeaders(defaultCORSExposedHeaders))

	if len(config.Methods) > 0 {
		h = append(h, handlers.AllowedMethods(config.Methods))
	}

	if len(config.Origins) > 0 {
		h = append(h, handlers.AllowedOrigins(config.Origins))
	}

	if len(config.Headers) > 0 {
		h = append(h, handlers.AllowedHeaders(config.Headers))
	}

	if len(config.ExposedHeaders) > 0 {
		h = append(h, handlers.ExposedHeaders(config.ExposedHeaders))
	}

	if config.AllowCredentials {
		h = append(h, handlers.AllowCredentials())
	}

	return
}

// Build read all the files on the impostersPath and add different
// handlers for each imposter
func (s *Server) Build() error {
	if s.proxy.mode == killgrave.ProxyAll {
		// not necessary load the imposters if you will use the tool as a proxy
		s.router.PathPrefix("/").HandlerFunc(s.proxy.Handler())
		return nil
	}

	if _, err := os.Stat(s.impostersPath); os.IsNotExist(err) {
		return fmt.Errorf("%w: the directory %s doesn't exists", err, s.impostersPath)
	}
	var impostersCh = make(chan []Imposter)
	var done = make(chan struct{})

	go func() {
		s.imposterFs.FindImposters(s.impostersPath, impostersCh)
		done <- struct{}{}
	}()
loop:
	for {
		select {
		case imposters := <-impostersCh:
			s.addImposterHandler(imposters)
			log.Printf("imposter %s loaded\n", imposters[0].Path)
		case <-done:
			close(impostersCh)
			close(done)
			break loop
		}
	}
	if s.proxy.mode == killgrave.ProxyMissing {
		s.router.NotFoundHandler = s.proxy.Handler()
	}
	return nil
}

// Run run launch a previous configured http server if any error happens while the starting process
// application will be crashed
func (s *Server) Run() {
	go func() {
		var tlsString string
		if s.secure {
			tlsString = "(TLS mode)"
		}
		log.Printf("The fake server is on tap now: %s%s\n", s.httpServer.Addr, tlsString)
		err := s.run(s.secure)
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}

func (s *Server) run(secure bool) error {
	if !secure {
		return s.httpServer.ListenAndServe()
	}

	cert, err := tls.X509KeyPair(serverCert, serverKey)
	if err != nil {
		log.Fatal(err)
	}

	s.httpServer.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	return s.httpServer.ListenAndServeTLS("", "")
}

// Shutdown shutdown the current http server
func (s *Server) Shutdown() error {
	log.Println("stopping server...")
	if err := s.httpServer.Shutdown(context.TODO()); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	return nil
}

func (s *Server) addImposterHandler(imposters []Imposter) {
	for _, imposter := range imposters {
		r := s.router.HandleFunc(imposter.Request.Endpoint, ImposterHandler(imposter)).
			Methods(imposter.Request.Method).
			MatcherFunc(MatcherBySchema(imposter))

		if imposter.Request.Headers != nil {
			for k, v := range *imposter.Request.Headers {
				r.HeadersRegexp(k, v)
			}
		}

		if imposter.Request.Params != nil {
			for k, v := range *imposter.Request.Params {
				r.Queries(k, v)
			}
		}
	}
}

func (s *Server) handleAll(h http.HandlerFunc) {
	s.router.PathPrefix("/").HandlerFunc(h)
}
