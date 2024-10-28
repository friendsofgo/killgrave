package http

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	killgrave "github.com/friendsofgo/killgrave/internal"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
)

// MatcherBySchema check if the request matching with the schema file
func MatcherBySchema(imposter Imposter) mux.MatcherFunc {
	return func(req *http.Request, rm *mux.RouteMatch) bool {
		err := validateSchema(imposter, req)

		if err != nil {
			log.WithFields(killgrave.LogFieldsFromRequest(req)).Warn(err)
			return false
		}
		return true
	}
}

func validateSchema(imposter Imposter, req *http.Request) error {
	if imposter.Request.SchemaFile == nil {
		return nil
	}

	var requestBodyBytes []byte

	defer func() {
		req.Body.Close()
		req.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))
	}()

	schemaFile := imposter.CalculateFilePath(*imposter.Request.SchemaFile)
	if _, err := os.Stat(schemaFile); os.IsNotExist(err) {
		return fmt.Errorf("%w: the schema file %s not found", err, schemaFile)
	}

	requestBodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("%w: impossible read the request body", err)
	}

	contentBody := string(requestBodyBytes)
	if contentBody == "" {
		return fmt.Errorf("unexpected empty body request")
	}

	schemaFilePath, _ := filepath.Abs(schemaFile)
	if err != nil {
		return fmt.Errorf("%w: impossible find the schema", err)
	}

	schemaBytes, err := os.ReadFile(schemaFilePath)
	if err != nil {
		return fmt.Errorf("%w: impossible read the schema", err)
	}

	schema := gojsonschema.NewStringLoader(string(schemaBytes))
	document := gojsonschema.NewStringLoader(string(requestBodyBytes))

	res, err := gojsonschema.Validate(schema, document)
	if err != nil {
		return fmt.Errorf("%w: error validating the json schema", err)
	}

	if !res.Valid() {
		for _, desc := range res.Errors() {
			return errors.New(desc.String())
		}
	}

	return nil
}
