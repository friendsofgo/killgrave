package http

import (
	"context"
	"crypto/tls"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"

	killgrave "github.com/friendsofgo/killgrave/internal"
	"github.com/friendsofgo/killgrave/internal/debugger"
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
	debugger      debugger.Debugger
}

// NewServer initialize the mock server
func NewServer(
	p string,
	r *mux.Router,
	httpServer *http.Server,
	proxyServer *Proxy,
	secure bool,
	debuggerCfg killgrave.ConfigDebugger,
) (Server, error) {
	var (
		dbg debugger.Debugger = debugger.NewNoOp()
		err error
	)

	if debuggerCfg.Enabled {
		dbg, err = debugger.NewWs(debuggerCfg)
		if err != nil {
			return Server{}, err
		}
	}

	return Server{
		impostersPath: p,
		router:        r,
		httpServer:    httpServer,
		proxy:         proxyServer,
		secure:        secure,
		debugger:      dbg,
	}, nil
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
		s.handleAll(s.proxy.Handler())
	}
	if _, err := os.Stat(s.impostersPath); os.IsNotExist(err) {
		return fmt.Errorf("%w: the directory %s doesn't exists", err, s.impostersPath)
	}
	var imposterConfigCh = make(chan killgrave.ImposterConfig)
	var done = make(chan bool)

	go func() {
		killgrave.FindImposters(s.impostersPath, imposterConfigCh)
		done <- true
	}()
loop:
	for {
		select {
		case imposterConfig := <-imposterConfigCh:
			var imposters []killgrave.Imposter
			err := s.unmarshalImposters(imposterConfig, &imposters)
			if err != nil {
				log.Printf("error trying to load %s imposter: %v", imposterConfig.FilePath, err)
			} else {
				s.addImposterHandler(imposters, imposterConfig)
				log.Printf("imposter %s loaded\n", imposterConfig.FilePath)
			}
		case <-done:
			close(imposterConfigCh)
			close(done)
			break loop
		}
	}
	if s.proxy.mode == killgrave.ProxyMissing {
		s.handleAll(s.proxy.Handler())
	}
	return nil
}

// Run launches a previously configured http server if any error happens while the starting process
// application will be crashed
func (s *Server) Run() {
	go func() {
		log.Printf("The fake server is on tap now: %s\n", s.httpServer.Addr)
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
// TODO: Handle ErrServerClosed on Serve methods
// TODO: Set up a maximum time to shut down
func (s *Server) Shutdown() error {
	log.Println("stopping server...")
	if err := s.httpServer.Shutdown(context.TODO()); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	return nil
}

func (s *Server) addImposterHandler(imposters []killgrave.Imposter, imposterConfig killgrave.ImposterConfig) {
	for _, imposter := range imposters {
		imposter.BasePath = filepath.Dir(imposterConfig.FilePath)
		r := s.router.HandleFunc(imposter.Request.Endpoint, ImposterHandler(s.debugger, imposter)).
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

func (s *Server) unmarshalImposters(imposterConfig killgrave.ImposterConfig, imposters *[]killgrave.Imposter) error {
	imposterFile, _ := os.Open(imposterConfig.FilePath)
	defer imposterFile.Close()

	bytes, _ := io.ReadAll(imposterFile)

	var parseError error

	switch imposterConfig.Type {
	case killgrave.JSONImposter:
		parseError = json.Unmarshal(bytes, imposters)
	case killgrave.YAMLImposter:
		parseError = yaml.Unmarshal(bytes, imposters)
	default:
		parseError = fmt.Errorf("Unsupported imposter type %v", imposterConfig.Type)
	}

	if parseError != nil {
		return fmt.Errorf("%w: error while unmarshalling imposter's file %s", parseError, imposterConfig.FilePath)
	}

	return nil
}

func (s *Server) handleAll(h http.HandlerFunc) {
	s.router.PathPrefix("/").HandlerFunc(h)
}
