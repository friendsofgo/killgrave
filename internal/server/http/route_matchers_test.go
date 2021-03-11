package http

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
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
		im  Imposter
		req *http.Request
		res int
	}{
		"correct request schema":               {Imposter{Request: requestWithSchema, Response: okResponse}, httpRequestA, http.StatusOK},
		"imposter without request schema":      {Imposter{Request: requestWithoutSchema, Response: okResponse}, httpRequestA, http.StatusOK},
		"malformatted schema file":             {Imposter{Request: requestWithWrongSchema, Response: okResponse}, httpRequestA, http.StatusBadRequest},
		"incorrect request schema":             {Imposter{Request: requestWithSchema, Response: okResponse}, httpRequestB, http.StatusBadRequest},
		"non-existing schema file":             {Imposter{Request: requestWithNonExistingSchema, Response: okResponse}, httpRequestB, http.StatusBadRequest},
		"empty body with required schema file": {Imposter{Request: requestWithSchema, Response: okResponse}, &http.Request{Body: emptyBody}, http.StatusBadRequest},
		"invalid request body":                 {Imposter{Request: requestWithSchema, Response: okResponse}, &http.Request{Body: wrongBody}, http.StatusBadRequest},
	}

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	for name, tt := range matcherData {
		rr := httptest.NewRecorder()
		h := SchemaValidationMiddleware(tt.im, dummyHandler)
		t.Run(name, func(t *testing.T) {
			h.ServeHTTP(rr, tt.req)
			code := rr.Code
			if code != tt.res {
				t.Fatalf(
					"error while matching by request schema - expected: %d, given: %d, error: %v",
					tt.res, code, rr.Body.String(),
				)
			}
		})

	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}
