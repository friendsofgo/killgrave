package http

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
)

func TestMatcherBySchema(t *testing.T) {
	bodyA := ioutil.NopCloser(bytes.NewReader([]byte("{\"type\": \"gopher\"}")))
	bodyB := ioutil.NopCloser(bytes.NewReader([]byte("{\"type\": \"cat\"}")))
	emptyBody := ioutil.NopCloser(bytes.NewReader([]byte("")))
	wrongBody := ioutil.NopCloser(errReader(0))

	schemaGopherFile := "test/testdata/imposters/schemas/type_gopher.json"
	schemaCatFile := "test/testdata/imposters/schemas/type_cat.json"
	schemeFailFile := "test/testdata/imposters/schemas/type_gopher_fail.json"

	requestWithoutSchema := Request{
		Method:     "POST",
		Endpoint:   "/login",
		SchemaFile: nil,
	}

	requestWithSchema := Request{
		Method:     "POST",
		Endpoint:   "/login",
		SchemaFile: &schemaGopherFile,
	}

	requestWithNonExistingSchema := Request{
		Method:     "POST",
		Endpoint:   "/login",
		SchemaFile: &schemaCatFile,
	}

	requestWithWrongSchema := Request{
		Method:     "POST",
		Endpoint:   "/login",
		SchemaFile: &schemeFailFile,
	}

	httpRequestA := &http.Request{Body: bodyA}
	httpRequestB := &http.Request{Body: bodyB}
	okResponse := Response{Status: http.StatusOK}

	var matcherData = map[string]struct {
		fn  mux.MatcherFunc
		req *http.Request
		res bool
	}{
		"correct request schema":               {MatcherBySchema(Imposter{Request: requestWithSchema, Response: okResponse}), httpRequestA, true},
		"imposter without request schema":      {MatcherBySchema(Imposter{Request: requestWithoutSchema, Response: okResponse}), httpRequestA, true},
		"malformatted schema file":             {MatcherBySchema(Imposter{Request: requestWithWrongSchema, Response: okResponse}), httpRequestA, false},
		"incorrect request schema":             {MatcherBySchema(Imposter{Request: requestWithSchema, Response: okResponse}), httpRequestB, false},
		"non-existing schema file":             {MatcherBySchema(Imposter{Request: requestWithNonExistingSchema, Response: okResponse}), httpRequestB, false},
		"empty body with required schema file": {MatcherBySchema(Imposter{Request: requestWithSchema, Response: okResponse}), &http.Request{Body: emptyBody}, false},
		"invalid request body":                 {MatcherBySchema(Imposter{Request: requestWithSchema, Response: okResponse}), &http.Request{Body: wrongBody}, false},
	}

	for name, tt := range matcherData {
		t.Run(name, func(t *testing.T) {
			res := tt.fn(tt.req, nil)
			if res != tt.res {
				t.Fatalf("error while matching by request schema - expected: %t, given: %t", tt.res, res)
			}
		})

	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}
