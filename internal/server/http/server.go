package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	killgrave "github.com/friendsofgo/killgrave/internal"
)

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
	httpServer    http.Server
	proxy         *Proxy
}

// NewServer initialize the mock server
func NewServer(p string, r *mux.Router, httpServer http.Server, proxyServer *Proxy) Server {
	return Server{
		impostersPath: p,
		router:        r,
		httpServer:    httpServer,
		proxy:         proxyServer,
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
		s.handleAll(s.proxy.Handler())
	}
	if _, err := os.Stat(s.impostersPath); os.IsNotExist(err) {
		return fmt.Errorf("%w: the directory %s doesn't exists", err, s.impostersPath)
	}
	var imposterFileCh = make(chan string)
	var done = make(chan bool)

	go func() {
		findImposters(s.impostersPath, imposterFileCh)
		done <- true
	}()
loop:
	for {
		select {
		case f := <-imposterFileCh:
			var imposters []Imposter
			err := s.unmarshalImposters(f, &imposters)
			if err != nil {
				log.Printf("error trying to load %s imposter: %v", f, err)
			} else {
				s.addImposterHandler(imposters, f)
				log.Printf("imposter %s loaded\n", f)
			}
		case <-done:
			close(imposterFileCh)
			close(done)
			break loop
		}
	}
	if s.proxy.mode == killgrave.ProxyMissing {
		s.handleAll(s.proxy.Handler())
	}
	return nil
}

// Run run launch a previous configured http server if any error happens while the starting process
// application will be crashed
func (s *Server) Run() {
	go func() {
		log.Printf("The fake server is on tap now: %s\n", s.httpServer.Addr)
		err := s.httpServer.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}

// Shutdown shutdown the current http server
func (s *Server) Shutdown() error {
	log.Println("stopping server...")
	if err := s.httpServer.Shutdown(context.TODO()); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	return nil
}

func (s *Server) addImposterHandler(imposters []Imposter, imposterFilePath string) {
	for _, imposter := range imposters {
		imposter.BasePath = filepath.Dir(imposterFilePath)
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

func (s Server) unmarshalImposters(imposterFileName string, imposters *[]Imposter) error {
	imposterFile, _ := os.Open(imposterFileName)
	defer imposterFile.Close()

	bytes, _ := ioutil.ReadAll(imposterFile)
	if err := json.Unmarshal(bytes, imposters); err != nil {
		return fmt.Errorf("%w: error while unmarshall imposter file %s", err, imposterFileName)
	}
	return nil
}

func (s *Server) handleAll(h http.HandlerFunc) {
	s.router.PathPrefix("/").HandlerFunc(h)
}
