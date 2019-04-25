package killgrave

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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
	if err := s.buildImposters(); err != nil {
		return err
	}

	return nil
}

func (s *Server) buildImposters() error {
	files, _ := ioutil.ReadDir(s.impostersPath)

	for _, f := range files {
		var imposter Imposter
		if err := s.buildImposter(f.Name(), &imposter); err != nil {
			return err
		}

		if imposter.Request.Endpoint == "" {
			continue
		}
		s.router.HandleFunc(imposter.Request.Endpoint, ImposterHandler(imposter)).
			Methods(imposter.Request.Method).
			MatcherFunc(func(req *http.Request, rm *mux.RouteMatch) bool {
				err := validateSchema(imposter, req)
				return err == nil
			})
	}

	return nil
}

func (s *Server) buildImposter(imposterFileName string, imposter *Imposter) error {
	f := s.impostersPath + "/" + imposterFileName
	imposterFile, _ := os.Open(f)
	defer imposterFile.Close()

	bytes, _ := ioutil.ReadAll(imposterFile)
	if err := json.Unmarshal(bytes, imposter); err != nil {
		return malformattedImposterError(fmt.Sprintf("error while unmarshall imposter file %s", f))
	}
	return nil
}

func validateSchema(imposter Imposter, req *http.Request) error {
	var b []byte

	defer func() {
		req.Body.Close()
		req.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	}()

	if imposter.Request.SchemaFile == nil {
		return nil
	}

	schemaFile := *imposter.Request.SchemaFile
	if _, err := os.Stat(schemaFile); os.IsNotExist(err) {
		return errors.Wrapf(err, "the schema file %s not found", schemaFile)
	}

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return errors.Wrapf(err, "impossible read the request body")
	}

	dir, _ := os.Getwd()
	schemaFilePath := "file://" + dir + "/" + schemaFile
	schema := gojsonschema.NewReferenceLoader(schemaFilePath)
	document := gojsonschema.NewStringLoader(string(b))

	res, err := gojsonschema.Validate(schema, document)
	if err != nil {
		return errors.Wrap(err, "error validating the json schema")
	}

	if !res.Valid() {
		for _, desc := range res.Errors() {
			return errors.New(desc.String())
		}
	}

	return nil
}
