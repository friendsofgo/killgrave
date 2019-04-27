package killgrave

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Server definition of mock server
type Server struct {
	impostersPath string
	router        *mux.Router
}

// NewServer initialize the mock server
func NewServer(p string, r *mux.Router) *Server {
	return &Server{
		impostersPath: p,
		router:        r,
	}
}

// Run read all the files on the impostersPath and creates different
// handlers for each imposter
func (s *Server) Run() error {
	if _, err := os.Stat(s.impostersPath); os.IsNotExist(err) {
		return invalidDirectoryError(fmt.Sprintf("the directory %s doesn't exists", s.impostersPath))
	}
	var imposterFileCh = make(chan string)
	var done = make(chan bool)

	go func() {
		findImposters(s.impostersPath, imposterFileCh)
		done <- true
	}()
	for {
		select {
		case f := <-imposterFileCh:
			var imposter Imposter
			err := s.buildImposter(f, &imposter)
			if err != nil {
				log.Printf("error trying to load %s imposter: %v", f, err)
			} else {
				if err := s.createImposterHandler(imposter); err != nil {
					log.Printf("%v on %s", err, f)
					break
				}
				log.Printf("imposter %s loaded\n", f)
			}
		case <-done:
			close(imposterFileCh)
			close(done)
			return nil
		}
	}
}

func (s *Server) createImposterHandler(imposter Imposter) error {
	if imposter.Request.Endpoint == "" {
		return errors.New("the request.endpoint file is required for an valid imposter")
	}
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
	return nil
}

func (s *Server) buildImposter(imposterFileName string, imposter *Imposter) error {
	imposterFile, _ := os.Open(imposterFileName)
	defer imposterFile.Close()

	bytes, _ := ioutil.ReadAll(imposterFile)
	if err := json.Unmarshal(bytes, imposter); err != nil {
		return malformattedImposterError(fmt.Sprintf("error while unmarshall imposter file %s", imposterFileName))
	}
	imposter.BasePath = filepath.Dir(imposterFileName)

	return nil
}
