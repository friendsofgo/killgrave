package http

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/gorilla/mux"
	"github.com/jbussdieker/golibxml"
	"github.com/krolaw/xsd"
	"github.com/xeipuuv/gojsonschema"
)

// MatcherBySchema check if the request matching with the schema file
func MatcherBySchema(imposter Imposter) mux.MatcherFunc {
	return func(req *http.Request, rm *mux.RouteMatch) bool {
		if imposter.Request.SchemaFile == nil {
			return true
		}

		var err error
		defer func() {
			if err != nil {
				log.Println(err)
			}
		}()

		bodyBytes, err := readBodyBytes(req)
		if err != nil {
			return false
		}

		schemaBytes, err := readSchemaBytes(imposter)
		if err != nil {
			return false
		}

		switch filepath.Ext(*imposter.Request.SchemaFile) {
		case ".json":
			err = validateJSONSchema(bodyBytes, schemaBytes)
		case ".xml", ".xsd":
			err = validateXMLSchema(bodyBytes, schemaBytes)
		default:
			err = errors.New("unknown schema file extension")
		}

		return err == nil
	}
}

func readBodyBytes(req *http.Request) (bodyBytes []byte, err error) {
	defer func() {
		req.Body.Close()
		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}()

	bodyBytes, err = ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: error reading the request body", err)
	}

	contentBody := string(bodyBytes)
	if contentBody == "" {
		return nil, fmt.Errorf("unexpected empty body request")
	}

	return
}

func readSchemaBytes(imposter Imposter) ([]byte, error) {
	schemaFile := imposter.CalculateFilePath(*imposter.Request.SchemaFile)
	if _, err := os.Stat(schemaFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: schema file not found", err)
	}

	schemaFilePath, err := filepath.Abs(schemaFile)
	if err != nil {
		return nil, fmt.Errorf("%w: error finding the schema file", err)
	}

	schemaBytes, err := ioutil.ReadFile(schemaFilePath)
	if err != nil {
		return nil, fmt.Errorf("%w: error reading the schema file", err)
	}

	return schemaBytes, nil
}

func validateJSONSchema(bodyBytes, schemaBytes []byte) error {
	schema := gojsonschema.NewStringLoader(string(schemaBytes))
	document := gojsonschema.NewStringLoader(string(bodyBytes))

	res, err := gojsonschema.Validate(schema, document)
	if err != nil {
		return fmt.Errorf("%w: error validating the json schema", err)
	}

	if !res.Valid() {
		for _, desc := range res.Errors() {
			return fmt.Errorf("%s: error validating the json schema", desc.String())
		}
	}

	return nil
}

func validateXMLSchema(bodyBytes, schemaBytes []byte) error {
	// Following the example from the official documentation:
	// https://godoc.org/github.com/krolaw/xsd
	schema, err := xsd.ParseSchema(schemaBytes)
	if err != nil {
		return fmt.Errorf("%w: error parsing the xml schema", err)
	}

	document := golibxml.ParseDoc(string(bodyBytes))
	defer document.Free()

	if err := schema.Validate(xsd.DocPtr(unsafe.Pointer(document.Ptr))); err != nil {
		return fmt.Errorf("%w: error validating the xml schema", err)
	}

	return nil
}
