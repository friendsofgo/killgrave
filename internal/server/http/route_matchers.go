package http

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/xeipuuv/gojsonschema"
)

// SchemaValidationMiddleware check if the request matching with the schema file
func SchemaValidationMiddleware(imposter Imposter, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := validateSchema(imposter, r)

		// TODO: inject the logger
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
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
		return fmt.Errorf("%w: the schema file %s not found", err, schemaFile)
	}

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return fmt.Errorf("%w: impossible read the request body", err)
	}

	contentBody := string(b)
	if contentBody == "" {
		return fmt.Errorf("unexpected empty body request")
	}

	dir, _ := os.Getwd()
	schemaFilePath := "file://" + dir + "/" + schemaFile
	schema := gojsonschema.NewReferenceLoader(schemaFilePath)
	document := gojsonschema.NewStringLoader(string(b))

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
