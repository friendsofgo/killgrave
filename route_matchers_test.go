package killgrave

import (
	"bytes"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestMatcherBySchema(t *testing.T) {
	bodyA := ioutil.NopCloser(bytes.NewReader([]byte("{\"type\": \"A\"}")))
	bodyB := ioutil.NopCloser(bytes.NewReader([]byte("{\"type\": \"B\"}")))

	schemaAFile := "test/testdata/schemas/type_a.json"
	schemaBFile := "test/testdata/schemas/type_b.json"


	requestWithoutSchema := Request{
		Method:     "POST",
		Endpoint:   "/login",
		SchemaFile: nil,
	}

	requestWithSchema := Request{
		Method:     "POST",
		Endpoint:   "/login",
		SchemaFile: &schemaAFile,
	}

	requestWithNonExistingSchema := Request{
		Method:     "POST",
		Endpoint:   "/login",
		SchemaFile: &schemaBFile,
	}

	okResponse := Response{Status: http.StatusOK}

	var matcherData = []struct {
		name string
		fn   mux.MatcherFunc
		req  *http.Request
		res  bool
	}{
		{"imposter without request schema", MatcherBySchema(Imposter{Request: requestWithoutSchema, Response: okResponse}), &http.Request{Body: bodyA}, true},
		{"correct request schema", MatcherBySchema(Imposter{Request: requestWithSchema, Response: okResponse}), &http.Request{Body: bodyA}, true},
		{"incorrect request schema", MatcherBySchema(Imposter{Request: requestWithSchema, Response: okResponse}), &http.Request{Body: bodyB}, false},
		{"non-existing schema file", MatcherBySchema(Imposter{Request: requestWithNonExistingSchema, Response: okResponse}), &http.Request{Body: bodyB}, false},
	}

	for _, tt := range matcherData {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.fn(tt.req, nil)
			if res != tt.res {
				t.Fatalf("error while matching by request schema - expected: %t, given: %t", tt.res, res)
			}
		})

	}
}
