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
	"gopkg.in/yaml.v2"

	killgrave "github.com/friendsofgo/killgrave/internal"
)

const (
	serverCert = `-----BEGIN CERTIFICATE-----
MIIDHjCCAgYCCQC9DwH2RNHUejANBgkqhkiG9w0BAQsFADBQMRYwFAYDVQQKDA1G
cmllbmRzIG9mIEdvMRIwEAYDVQQDDAlraWxsZ3JhdmUxIjAgBgkqhkiG9w0BCQEW
E2l0QGZyaWVuZHNvZmdvLnRlY2gwIBcNMjEwMzEwMjI0OTM0WhgPMjEyMTAyMTQy
MjQ5MzRaMFAxFjAUBgNVBAoMDUZyaWVuZHMgb2YgR28xEjAQBgNVBAMMCWtpbGxn
cmF2ZTEiMCAGCSqGSIb3DQEJARYTaXRAZnJpZW5kc29mZ28udGVjaDCCASIwDQYJ
KoZIhvcNAQEBBQADggEPADCCAQoCggEBALKWzDXl072p4FV9S16u2gwUStdHRYZd
xyyE5BEusNADQVB5nQ66/4VhxGebQWootaaZStWUl8Et07bvHhk0cpy3Jd6gQaRu
g+FXiwtLsE8xMbdAbX4QvQCwZ39Ay6uNHZ4GRX5qVr46myvnL+GaemgcZpbsGGrE
reAba2WUXNpDIBEQ/aCEe+rwpdN7d1obLMoT3H9iPqJwQAaTaOzTg3WhaSYOkInF
RAG/FOCXoDZUb8pDWSefF8ZF26uRF9iwiTSTu29hR65ld4V0RS9kUvNu010tMew8
dBGRwRcTaqmZyWZKe/5EHaZD9p7MmfrQs7xup0hHeg0EZuC89S/uDkcCAwEAATAN
BgkqhkiG9w0BAQsFAAOCAQEAHPwIBPl25v/kGAtPZsmHsKMl2TgfRYjid0yDkbPm
Nb+BzsCoUCnA5EMpIyBbyapsI6lBe8SObwRXSra9XwOSEWwsSpql3occvSNCCaS/
Ti6NSUpFzq/wIMC9JU4KH7SPXSTHOJ32GfbaI5TECu27hvDUBmML4zyeGSf3h6pa
mpf1bMqqXgxfzWaLn3Q859ejR2whS4eeqMFKl3RxV8QvW5JfUij+2LwCzttFcKKg
i4s0QCjR8eIreIcMvZP/T7zSIcw1dx1lFh8b+hbuVXhOojMY/hHWEPNf/Uu0PZy7
THHLAlBA4Lbms9gg0dF8czJ3AzH+UG9xpRkR4KQ1Md+zJw==
-----END CERTIFICATE-----`

	serverKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAspbMNeXTvangVX1LXq7aDBRK10dFhl3HLITkES6w0ANBUHmd
Drr/hWHEZ5tBaii1pplK1ZSXwS3Ttu8eGTRynLcl3qBBpG6D4VeLC0uwTzExt0Bt
fhC9ALBnf0DLq40dngZFfmpWvjqbK+cv4Zp6aBxmluwYasSt4BtrZZRc2kMgERD9
oIR76vCl03t3WhssyhPcf2I+onBABpNo7NODdaFpJg6QicVEAb8U4JegNlRvykNZ
J58XxkXbq5EX2LCJNJO7b2FHrmV3hXRFL2RS827TXS0x7Dx0EZHBFxNqqZnJZkp7
/kQdpkP2nsyZ+tCzvG6nSEd6DQRm4Lz1L+4ORwIDAQABAoIBABV5HES+xZ7gdiDR
V+aij4U0S2tnHmzxialIsUN/obLhMVFDziafRWn8P2lVuZ/SFUVa2SylGToZEIPG
bJALRlyhiOQj0MC8qQ7HP+izyRc8iwXFsWSfDpqum0Mpv1N5PD5r8p8omhV1ZoL4
4UD3GhC6mXs8GBN+Yom3wkoMdL2pcRaLDqkQTJnZquGo6QndS0X5liqmRWROXug6
Uona653GyEpokP6saepfdO5rXMHdTeFCRZPfloNK9ulw62Evww7tviokNTVABz/f
tbcdHewrNMFL1y8QuF+/6QJtaV8utLwpdwLTHvlcPuSF92jXqi2mSS9ptqYHwB3B
3YDQ/QECgYEA3U8RcHnaxEZyP2xvBJ7pNqMG4b7bMak4W9pqfzvXtt3ZBNWQQewN
IrRKdKgALf5VBjcSDjr5sUR8lRqAFURMkr97oDmhFFMV6ywo/t/DjdHlCIHPIJYD
LjBxKlD9k7f73xtKzLg5hJXR8WBJ9JBfPNv10icEyNNYnGwlfuu3vfECgYEAzpVr
wo5lcJ7VDWi+pp7rnz2BBQKoUPux8G7p7XQ/+Yhe75TRn7qm9SjD6uGEHEhN6PjO
1DYHkIRWVDK6kTIL+qiHLWgYEUM63r1jhgY+n+x7d82Sz798WFEiYjLGL/CTIdFs
vSsz3QMWAhm08E4k9MFE5ATWYKpCVQ/yUg2ut7cCgYBlca4DyceO+t+51OGa06EB
a39nEU52mCP+bsMsaWj7KPwmrCKBJUvsIYqTqMLUUmX1AF9laIE2UbdtvYUCupkD
F4T6sA/3OhKtB0QPeNCx/Imo+Z/RRxJUJN5q0E88XDS3U1JZPwUWknp203VzBo6x
Xf5zg3E9ASv4H9acND64cQKBgHZavukxQca7CN7s0sWNKPsLdp6TPjFfcjuIn/cN
8hUZXyKtxUdY3Yx5dX1dBJ5bgl9mJMEJz12po/gLND45SQmrgf6us5M4TEMOiDVh
4IEpMDecDG9/ilLi8OsHoeoXT4RBgqYCWW1W9kXvym0eqCedjsWAS/4HrYckYrVF
54KTAoGAc1L9pjJhQZuf6uE0rh0WiUIPE9R2q2FKQqfPss0KHv9pENH/DPI3y3jF
1/BE5jythAgP8biHGC126xW4dP/2NyTN1Ek8/N5ibdSMSG02FNFQNmqmdwP/HUXj
vfIUhmkGtpT6HYWjKI2+47PcdpJgZg4X4mTX5frT0eSeiNwfePs=
-----END RSA PRIVATE KEY-----
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
	var imposterConfigCh = make(chan ImposterConfig)
	var done = make(chan bool)

	go func() {
		findImposters(s.impostersPath, imposterConfigCh)
		done <- true
	}()
loop:
	for {
		select {
		case imposterConfig := <-imposterConfigCh:
			var imposters []Imposter
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

func (s *Server) addImposterHandler(imposters []Imposter, imposterConfig ImposterConfig) {
	for _, imposter := range imposters {
		imposter.BasePath = filepath.Dir(imposterConfig.FilePath)
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

func (s *Server) unmarshalImposters(imposterConfig ImposterConfig, imposters *[]Imposter) error {
	imposterFile, _ := os.Open(imposterConfig.FilePath)
	defer imposterFile.Close()

	bytes, _ := ioutil.ReadAll(imposterFile)

	var parseError error

	switch imposterConfig.Type {
	case JSONImposter:
		parseError = json.Unmarshal(bytes, imposters)
	case YAMLImposter:
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
