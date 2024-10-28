package http

import (
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"net/http"
	"os"

	killgrave "github.com/friendsofgo/killgrave/internal"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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
	router     *mux.Router
	httpServer *http.Server
	proxy      *Proxy
	secure     bool
	imposterFs ImposterFs
}

// NewServer initialize the mock server
func NewServer(r *mux.Router, httpServer *http.Server, proxyServer *Proxy, secure bool, fs ImposterFs) Server {
	return Server{
		router:     r,
		httpServer: httpServer,
		proxy:      proxyServer,
		secure:     secure,
		imposterFs: fs,
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
		log.Infoln("ProxyAll mode enabled, no imposter will be used")
		s.router.PathPrefix("/").HandlerFunc(s.proxy.Handler())
		return nil
	}

	var impostersCh = make(chan []Imposter)
	var done = make(chan struct{})

	go func() {
		s.imposterFs.FindImposters(impostersCh)
		done <- struct{}{}
	}()
loop:
	for {
		select {
		case imposters := <-impostersCh:
			s.addImposterHandler(imposters)
		case <-done:
			close(impostersCh)
			close(done)
			break loop
		}
	}
	if s.proxy.mode == killgrave.ProxyMissing {
		log.Infof("Proxying missed requests to: %v", s.proxy.url)
		s.router.NotFoundHandler = s.proxy.Handler()
	} else {
		log.Infoln("No proxy has been configured for non-matching requests, defaulting to a 404 response")
		s.router.NotFoundHandler = s.defaultNotFoundHandler()
	}
	return nil
}

// Run launch a previous configured http server if any error happens while the starting process
// application will be crashed
func (s *Server) Run() {
	go func() {
		var tlsString string
		if s.secure {
			tlsString = "(TLS mode)"
		}
		log.Infof("The fake server is on tap now: %s%s\n", s.httpServer.Addr, tlsString)
		err := s.run(s.secure)
		if err != http.ErrServerClosed {
			log.Error(err)
			os.Exit(1)
		}
	}()
}

func (s *Server) run(secure bool) error {
	if !secure {
		return s.httpServer.ListenAndServe()
	}

	cert, err := tls.X509KeyPair(serverCert, serverKey)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	s.httpServer.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	return s.httpServer.ListenAndServeTLS("", "")
}

// Shutdown shutdown the current http server
func (s *Server) Shutdown() error {
	log.Info("stopping server...")
	if err := s.httpServer.Shutdown(context.TODO()); err != nil {
		return fmt.Errorf("Server Shutdown Failed:%+v", err)
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
		log.WithFields(imposter.LogFields()).Debugln("imposter loaded")
	}
}

func (s *Server) defaultNotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(killgrave.LogFieldsFromRequest(r)).Debugf("Request didn't match any imposter, and proxyMode is %v", s.proxy.mode)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 page not found\n"))
	}
}

// not used?
func (s *Server) handleAll(h http.HandlerFunc) {
	s.router.PathPrefix("/").HandlerFunc(h)
}
