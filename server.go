package killgrave

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

// Server definition of mock server
type Server struct {
	impostersPath string
	handler       http.Handler
}

// NewServer initialize the mock server
func NewServer(p string, h http.Handler) *Server {
	return &Server{
		impostersPath: p,
		handler:       h,
	}
}

// Run read all the files on the impostersPath and creates different
// handlers for each imposter
func (s *Server) Run() error {
	if _, err := os.Stat(s.impostersPath); os.IsNotExist(err) {
		return fmt.Errorf("the directory %s doesn't exists", s.impostersPath)
	}
	imposters, err := s.fetchImposters()
	if err != nil {
		return err
	}
	fmt.Println(imposters)
	return nil
}

func (s *Server) fetchImposters() ([]Imposter, error) {
	var imposters []Imposter
	files, err := ioutil.ReadDir(s.impostersPath)
	if err != nil {
		return imposters, errors.Wrapf(err, "an error ocurred while read dir %s", s.impostersPath)
	}

	for _, f := range files {
		var imposter Imposter
		if err := s.buildImposter(f.Name(), &imposter); err != nil || imposter.Request.Endpoint == "" {
			log.Println(err)
			continue
		}

		imposters = append(imposters, imposter)
	}

	return imposters, nil
}

func (s *Server) buildImposter(imposterFileName string, imposter *Imposter) error {
	f := s.impostersPath + "/" + imposterFileName
	imposterFile, err := os.Open(f)
	if err != nil {
		return errors.Wrapf(err, "error reading imposter file: %s", f)
	}
	defer imposterFile.Close()

	bytes, _ := ioutil.ReadAll(imposterFile)
	if err := json.Unmarshal(bytes, imposter); err != nil {
		return errors.Wrapf(err, "error while unmarshall imposter file %s", f)
	}
	return nil
}
