package killgrave

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
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

	schemaFile := "test/testdata/imposters/schemas/create_gopher_request.json"
	bodyFile := "test/testdata/imposters/responses/create_gopher_response.json"
	bodyFileFake := "test/testdata/imposters/responses/create_gopher_response_fail.json"
	body := `{"test":true}`

	validRequest := Request{
		Method:     "POST",
		Endpoint:   "/gophers",
		SchemaFile: &schemaFile,
		Headers:    &headers,
	}

	f, _ := os.Open(bodyFile)
	defer f.Close()
	expectedBodyFileData, _ := ioutil.ReadAll(f)

	var dataTest = []struct {
		name         string
		imposter     Imposter
		expectedBody string
		statusCode   int
	}{
		{"valid imposter with body", Imposter{Request: validRequest, Response: Response{Status: http.StatusOK, Headers: &http.Header{"Content-Type": []string{"application/json"}}, Body: body}}, body, http.StatusOK},
		{"valid imposter with bodyFile", Imposter{Request: validRequest, Response: Response{Status: http.StatusOK, Headers: &http.Header{"Content-Type": []string{"application/json"}}, BodyFile: &bodyFile}}, string(expectedBodyFileData), http.StatusOK},
		{"valid imposter with not exists bodyFile", Imposter{Request: validRequest, Response: Response{Status: http.StatusOK, Headers: &http.Header{"Content-Type": []string{"application/json"}}, BodyFile: &bodyFileFake}}, "", http.StatusOK},
	}

	for _, tt := range dataTest {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/gophers", bytes.NewBuffer(bodyRequest))
			if err != nil {
				t.Fatalf("could not created request: %v", err)
			}
			req.Header = headers

			rec := httptest.NewRecorder()
			handler := http.HandlerFunc(ImposterHandler(tt.imposter))

			handler.ServeHTTP(rec, req)
			if status := rec.Code; status != tt.statusCode {
				t.Errorf("handler expected %d code and got: %d code", tt.statusCode, status)
			}

			if rec.Body.String() != tt.expectedBody {
				t.Errorf("handler expected %s body and got: %s body", tt.expectedBody, rec.Body.String())
			}

		})
	}
}

func TestInvalidRequestWithSchema(t *testing.T) {
	validRequest := []byte(`{
		"data": {
			"type": "gophers",
		  "attributes": {
			"name": "Zebediah",
			"color": "Purple"
		  }
		}
	  }`)

	var dataTest = []struct {
		name       string
		imposter   Imposter
		statusCode int
		request    []byte
	}{
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

func TestInvalidHeaders(t *testing.T) {
	validRequest := []byte(`{
		"data": {
			"type": "gophers",
		  "attributes": {
			"name": "Zebediah",
			"color": "Purple",
			"age": 55
		  }
		}
	  }`)
	schemaFile := "test/testdata/imposters/schemas/create_gopher_request.json"
	var expectedHeaders = make(http.Header)
	expectedHeaders.Add("Content-Type", "application/json")
	expectedHeaders.Add("Authorization", "Bearer gopher")

	var wrongHeaders = make(http.Header)
	wrongHeaders.Add("Content-Type", "application/json")
	wrongHeaders.Add("Authorization", "Bearer bad gopher")

	var dataTest = []struct {
		name       string
		imposter   Imposter
		statusCode int
		headers    http.Header
	}{
		{"missing headers", Imposter{Request: Request{Method: "POST", Endpoint: "/gophers", SchemaFile: &schemaFile, Headers: &expectedHeaders}}, http.StatusBadRequest, make(http.Header)},
		{"wrong headers", Imposter{Request: Request{Method: "POST", Endpoint: "/gophers", SchemaFile: &schemaFile, Headers: &expectedHeaders}}, http.StatusBadRequest, wrongHeaders},
	}

	for _, tt := range dataTest {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/gophers", bytes.NewBuffer(validRequest))
			if err != nil {
				t.Fatalf("could not created request: %v", err)
			}
			req.Header = tt.headers

			rec := httptest.NewRecorder()
			handler := http.HandlerFunc(ImposterHandler(tt.imposter))

			handler.ServeHTTP(rec, req)
			if status := rec.Code; status != tt.statusCode {
				t.Fatalf("handler expected %d code and got: %d code", tt.statusCode, status)
			}
		})
	}
}
