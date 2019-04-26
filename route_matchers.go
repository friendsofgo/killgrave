package killgrave

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
)

// MatcherBySchema check if the request matching with the schema file
func MatcherBySchema(imposter Imposter) mux.MatcherFunc {
	return func(req *http.Request, rm *mux.RouteMatch) bool {
		err := validateSchema(imposter, req)

		// TODO: inject the logger
		if err != nil {
			log.Println(err)
			return false
		}
		return true
	}
}

func validateSchema(imposter Imposter, req *http.Request) error {
	if imposter.Request.SchemaFile == nil {
		return nil
	}

	var b []byte

	defer func() {
		req.Body.Close()
		req.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	}()

	schemaFile := imposter.CalculateFilePath(*imposter.Request.SchemaFile)
	if _, err := os.Stat(schemaFile); os.IsNotExist(err) {
		return errors.Wrapf(err, "the schema file %s not found", schemaFile)
	}

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return errors.Wrapf(err, "impossible read the request body")
	}

	contentBody := string(b)
	if contentBody == "" {
		return fmt.Errorf("expected an body request and got empy body request")
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
