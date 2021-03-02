package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestImposterHandler(t *testing.T) {
	var headers = make(map[string]string)
	headers["Content-Type"] = "application/json"

	bodyXMLFile := "test/testdata/imposters/responses/create_gopher_response.xml"
	f, _ := os.Open(bodyXMLFile)
	defer f.Close()
	expectedXMLBodyFileData, _ := ioutil.ReadAll(f)

	xmlBody := `<?xml version="1.0" encoding="UTF-8" ?><data><type>gophers</type></data>`

	bodyJSONFile := "test/testdata/imposters/responses/create_gopher_response.json"
	f, _ = os.Open(bodyJSONFile)
	defer f.Close()
	expectedJSONBodyFileData, _ := ioutil.ReadAll(f)

	jsonBody := `{"test":true}`

	bodyFileFake := "test/testdata/imposters/responses/create_gopher_response_fail.json"

	var dataTest = []struct {
		name            string
		imposter        Imposter
		expectedHeaders map[string]string
		expectedBody    string
		statusCode      int
	}{
		{"valid XML imposter with body", Imposter{Response: Response{Status: http.StatusOK, Headers: &headers, Body: xmlBody}}, headers, xmlBody, http.StatusOK},
		{"valid XML imposter with bodyXMLFile", Imposter{Response: Response{Status: http.StatusOK, Headers: &headers, BodyFile: &bodyXMLFile}}, headers, string(expectedXMLBodyFileData), http.StatusOK},

		{"valid JSON imposter with body", Imposter{Response: Response{Status: http.StatusOK, Headers: &headers, Body: jsonBody}}, headers, jsonBody, http.StatusOK},
		{"valid JSON imposter with bodyJSONFile", Imposter{Response: Response{Status: http.StatusOK, Headers: &headers, BodyFile: &bodyJSONFile}}, headers, string(expectedJSONBodyFileData), http.StatusOK},

		{"valid imposter with non-existing body file", Imposter{Response: Response{Status: http.StatusOK, Headers: &headers, BodyFile: &bodyFileFake}}, headers, "", http.StatusOK},
	}

	for _, tt := range dataTest {
		t.Run(tt.name, func(t *testing.T) {
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

			req, err := http.NewRequest("POST", "/gophers", bytes.NewBuffer(bodyRequest))
			if err != nil {
				t.Fatalf("could not created request: %v", err)
			}

			rec := httptest.NewRecorder()
			handler := ImposterHandler(tt.imposter)

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
