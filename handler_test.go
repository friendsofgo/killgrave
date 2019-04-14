package killgrave

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestImposterHandler(t *testing.T) {
	bodyRequest := []byte(`{
		"data": {
			"type": "gophers",
		  "attributes": {
			"name": "Zebediah",
			"color": "Purple",
			"age": 55
		  }
		}
	  }`)
	var headers = make(http.Header)
	headers.Add("Content-Type", "application/json")

	req, err := http.NewRequest("POST", "/gophers", bytes.NewBuffer(bodyRequest))
	if err != nil {
		t.Fatalf("could not created request: %v", err)
	}
	req.Header = headers

	schemaFile := "test/testdata/schemas/create_gopher_request.json"
	expectedBody := `{"test":true}`
	expectedStatusCode := 200

	imposter := Imposter{
		Request: Request{
			Method:     "POST",
			Endpoint:   "/gophers",
			SchemaFile: &schemaFile,
			Headers:    &headers,
		},
		Response: Response{
			Status:      expectedStatusCode,
			Body:        expectedBody,
			ContentType: "application/json",
		},
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(ImposterHandler(imposter))

	handler.ServeHTTP(rec, req)

	if status := rec.Code; status != expectedStatusCode {
		t.Errorf("handler expected %d code and got: %d code", expectedStatusCode, status)
	}

	if rec.Body.String() != expectedBody {
		t.Errorf("handler expected %s body and got: %s body", expectedBody, rec.Body.String())
	}
}

func TestInvalidRequestWithSchema(t *testing.T) {
	wrongRequest := []byte(`{
		"data": {
			"type": "gophers",
		  "attributes": {
			"name": "Zebediah",
			"color": "Purple"
		  }
		}
	  }`)
	validRequest := []byte(`{
		"data": {
			"type": "gophers",
		  "attributes": {
			"name": "Zebediah",
			"color": "Purple"
		  }
		}
	  }`)
	notExistFile := "failSchema"
	wrongSchema := "test/testdata/schemas/create_gopher_request_fail.json"
	validSchema := "test/testdata/schemas/create_gopher_request.json"

	var dataTest = []struct {
		name       string
		imposter   Imposter
		statusCode int
		request    []byte
	}{
		{"schema file not found", Imposter{Request: Request{Method: "POST", Endpoint: "/gophers", SchemaFile: &notExistFile}}, http.StatusBadRequest, validRequest},
		{"wrong schema", Imposter{Request: Request{Method: "POST", Endpoint: "/gophers", SchemaFile: &wrongSchema}}, http.StatusBadRequest, validRequest},
		{"request invalid", Imposter{Request: Request{Method: "POST", Endpoint: "/gophers", SchemaFile: &validSchema}}, http.StatusBadRequest, wrongRequest},
		{"valid request no schema", Imposter{Request: Request{Method: "POST", Endpoint: "/gophers"}, Response: Response{Status: http.StatusOK, Body: "test ok"}}, http.StatusOK, validRequest},
	}

	for _, tt := range dataTest {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/gophers", bytes.NewBuffer(tt.request))
			if err != nil {
				t.Fatalf("could not created request: %v", err)
			}
			rec := httptest.NewRecorder()
			handler := http.HandlerFunc(ImposterHandler(tt.imposter))

			handler.ServeHTTP(rec, req)
			if status := rec.Code; status != tt.statusCode {
				t.Fatalf("handler expected %d code and got: %d code", tt.statusCode, status)
			}
		})
	}
}
