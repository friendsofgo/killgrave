package http

import (
	"context"
	"crypto/tls"
	_ "embed"
	"log"
	"net/http"
	"os"
	"sync"

	killgrave "github.com/friendsofgo/killgrave/internal"
	sc "github.com/friendsofgo/killgrave/internal/serverconfig"
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

// Server definition of mock server
type Server struct {
	router     *mux.Router
	httpServer *http.Server
	proxy      *Proxy
	secure     bool
	imposterFs ImposterFs
	serverCfg  *sc.ServerConfig
	dumpCh     chan *RequestData
	wg         *sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewServer initialize the mock server
func NewServer(r *mux.Router, httpServer *http.Server, proxyServer *Proxy, secure bool, fs ImposterFs, options ...sc.ServerOption) Server {
	ctx, cancel := context.WithCancel(context.Background())
	cfg := &sc.ServerConfig{}
	for _, opt := range options {
		opt(cfg)
	}
	return Server{
		router:     r,
		httpServer: httpServer,
		proxy:      proxyServer,
		secure:     secure,
		imposterFs: fs,
		serverCfg:  cfg,
		wg:         &sync.WaitGroup{},
		ctx:        ctx,
		cancel:     cancel,
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

	// only intantiate the request dump if we need it
	if s.dumpCh == nil && s.serverCfg.LogWriter != nil {
		s.dumpCh = make(chan *RequestData, 1000)
		// Start the RequestWriter goroutine with context
		s.wg.Add(1)
		go RequestWriter(s.ctx, s.wg, s.serverCfg.LogWriter, s.dumpCh)
	}

	// setup the logging handler
	var handler http.Handler = s.router
	if s.serverCfg.LogLevel > 0 || shouldRecordRequest(s) {
		handler = CustomLoggingHandler(log.Writer(), handler, s)
	}
	s.httpServer.Handler = handlers.CORS(s.serverCfg.CORSOptions...)(handler)

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

// Run launch a previous configured http server if any error happens while the starting process
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
	if err := s.httpServer.Shutdown(s.ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	// Cancel the context to stop the RequestWriter goroutine
	s.cancel()

	// wait for all goroutines to finish
	s.wg.Wait()
	if s.serverCfg.LogWriter != nil {
		if f, ok := s.serverCfg.LogWriter.(*os.File); ok {
			f.Close()
		}
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
