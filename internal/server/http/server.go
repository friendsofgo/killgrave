package http

import (
	"context"
	"crypto/tls"
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

const (
	serverCert = `-----BEGIN CERTIFICATE-----
MIICBzCCAXACCQCKZvfDub3HcjANBgkqhkiG9w0BAQsFADBIMQswCQYDVQQGEwJH
QjEPMA0GA1UEBwwGTG9uZG9uMRQwEgYDVQQLDAtFbmdpbmVlcmluZzESMBAGA1UE
AwwJbG9jYWxob3N0MB4XDTIwMTAwNDE1MjYxMVoXDTIwMTEwMzE1MjYxMVowSDEL
MAkGA1UEBhMCR0IxDzANBgNVBAcMBkxvbmRvbjEUMBIGA1UECwwLRW5naW5lZXJp
bmcxEjAQBgNVBAMMCWxvY2FsaG9zdDCBnzANBgkqhkiG9w0BAQEFAAOBjQAwgYkC
gYEA1oLsjCcJRxzUeSpDWGJetA/v73yoRdF7//Yqd5wbH6ku+0P9HaYOYxpeGJwa
wbfrikDKfYtVAjWfYHe0qhmX342hfPipPsFM+gX0efjb+ULPuS1FV0ts4MHWLuMQ
HCPON9RokXYKEm6dH3T6Ezcw+Ku76/xf3HV2PiEe0lla9m0CAwEAATANBgkqhkiG
9w0BAQsFAAOBgQCuFU6z1p1PcYFb4PCo8QOXX8W70jFgE2uaFwvLZJ06tjR3yh1S
JK/kNAt689epQVLrNXZuWPl6M0Q4wV/0p6uTr6iOFHXWtApJh2cMRJeKScOEF0ex
B5clSLb8bL6cPVN0pCBsmddaK6d2Tyh85+mbYQN1x2EfMg+H7/IHZkXwdQ==
-----END CERTIFICATE-----
`

	serverKey = `-----BEGIN PRIVATE KEY-----
MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBANaC7IwnCUcc1Hkq
Q1hiXrQP7+98qEXRe//2KnecGx+pLvtD/R2mDmMaXhicGsG364pAyn2LVQI1n2B3
tKoZl9+NoXz4qT7BTPoF9Hn42/lCz7ktRVdLbODB1i7jEBwjzjfUaJF2ChJunR90
+hM3MPiru+v8X9x1dj4hHtJZWvZtAgMBAAECgYAgxHET26asKTg/pfgRmT00LjcN
kzI1MBHMALNt//eYt4RIt5MDo2kRNGbpRXdE3i5puQn1cYIzyzMkTkTXsv8izGDq
kl2wC4xMyorYv3AzMzx868+R6UJFOMRo0SCTGDKoeI9rydJoulCbuIpWpsiNXj7x
Y+j4+TpJ4M5cURx04QJBAPLyY4R9Tgtqieee1j+NER0Yjk9RLyqeKRg+ITjSqE9H
WygWRQLJONwLQWyd4NMR4TTmCQZmFUgWNiKb8hUhSGUCQQDiCWuo/xoumbp8S3rn
Y9wD3cbmH8oa0avmAQ/mz08Jv4mpR/v4JTFRElboNvwf3WtRzKoXvzIjpq8AbEr9
zmFpAkBZory9AU5uP9yprJz3zaBmz8yRzy5L1xbqbuHrCS44MeecHrtPj9Z+uVhm
Lsnolkw1LDpgNgHcGvXWRxtGWIVRAkA4HXqa0+oeE5AWd26lr0bZtt9AFjhIfDEe
wri95k2K8AAvBG3rZuBdbh4hPDVPe9q+zf6UMqUx8JmVk0zywZ+xAkAFshMafMou
7iqsFjCsWv56RnHvMVq/VbcCIXqlEYUIa+ygh4h8CepuzHjeiQ0EbqOJ32NAl14E
vQxnjxWH8M9y
-----END PRIVATE KEY-----
`
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
	httpServer    *http.Server
	proxy         *Proxy
	secure        bool
}

// NewServer initialize the mock server
func NewServer(p string, r *mux.Router, httpServer *http.Server, proxyServer *Proxy, secure bool) Server {
	return Server{
		impostersPath: p,
		router:        r,
		httpServer:    httpServer,
		proxy:         proxyServer,
		secure:        secure,
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

	cert, err := tls.X509KeyPair([]byte(serverCert), []byte(serverKey))
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
