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
		switch filepath.Ext(*imposter.Request.SchemaFile) {
		case ".json":
			err = validateJSONSchema(imposter, req)
		case ".xml", ".xsd":
			err = validateXMLSchema(imposter, req)
		default:
			err = errors.New("unknown schema file extension")
		}

		// TODO: inject the logger
		if err != nil {
			log.Println(err)
			return false
		}
		return true
	}
}

func validateJSONSchema(imposter Imposter, req *http.Request) error {
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

func validateXMLSchema(imposter Imposter, req *http.Request) error {
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

	// Following example from https://godoc.org/github.com/krolaw/xsd
	dir, _ := os.Getwd()
	schemaFilePath := dir + "/" + schemaFile
	schemaBytes, err := ioutil.ReadFile(schemaFilePath)
	if err != nil {
		println(schemaFilePath)
		return fmt.Errorf("%w: error reading from xml schema file location", err)
	}

	schema, err := xsd.ParseSchema(schemaBytes)
	if err != nil {
		return fmt.Errorf("%w: error parsing xml schema", err)
	}

	document := golibxml.ParseDoc(string(b))
	if document == nil {
		return fmt.Errorf("Error parsing the xml request")
	}
	defer document.Free()

	if err := schema.Validate(xsd.DocPtr(unsafe.Pointer(document.Ptr))); err != nil {
		return fmt.Errorf("%w: error validating the xml request", err)
	}

	return nil
}
