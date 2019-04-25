package killgrave

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"net/http"
	"os"
)

// MatcherBySchema ...
func MatcherBySchema(imposter Imposter) mux.MatcherFunc {
	return func(req *http.Request, rm *mux.RouteMatch) bool {
		err := validateSchema(imposter, req)
		return err == nil
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
