package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestImposterJSONHandler(t *testing.T) {
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
	var headers = make(map[string]string)
	headers["Content-Type"] = "application/json"

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
		{"valid json imposter with body", Imposter{Request: validRequest, Response: Response{Status: http.StatusOK, Headers: &headers, Body: body}}, body, http.StatusOK},
		{"valid json imposter with bodyFile", Imposter{Request: validRequest, Response: Response{Status: http.StatusOK, Headers: &headers, BodyFile: &bodyFile}}, string(expectedBodyFileData), http.StatusOK},
		{"valid json imposter with not exists bodyFile", Imposter{Request: validRequest, Response: Response{Status: http.StatusOK, Headers: &headers, BodyFile: &bodyFileFake}}, "", http.StatusOK},
	}

	for _, tt := range dataTest {
		t.Run(tt.name, func(t *testing.T) {
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

	schemaFile := "test/testdata/imposters/schemas/create_gopher_request"

	var dataTest = []struct {
		name       string
		imposter   Imposter
		statusCode int
		request    []byte
	}{
		{"valid request no schema", Imposter{Request: Request{Method: "POST", Endpoint: "/gophers"}, Response: Response{Status: http.StatusOK, Body: "test ok"}}, http.StatusOK, validRequest},
		{"valid request schema file with no extension", Imposter{Request: Request{Method: "POST", Endpoint: "/gophers", SchemaFile: &schemaFile}, Response: Response{Status: http.StatusOK, Body: "test ok"}}, http.StatusOK, validRequest},
	}

	for _, tt := range dataTest {

		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/gophers", bytes.NewBuffer(tt.request))
			if err != nil {
				t.Fatalf("could not created request: %v", err)
			}
			rec := httptest.NewRecorder()
			handler := ImposterHandler(tt.imposter)

			handler.ServeHTTP(rec, req)
			if status := rec.Code; status != tt.statusCode {
				t.Fatalf("handler expected %d code and got: %d code", tt.statusCode, status)
			}
		})
	}
}

func TestImposterXMLHandler(t *testing.T) {
	bodyRequest := []byte(`<?xml version="1.0" encoding="UTF-8" ?>
	<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">
	  <xs:element name="data">
		<xs:complexType>
		  <xs:sequence>
			<xs:element type="xs:string" name="type"/>
			<xs:element type="xs:string" name="id"/>
			<xs:element name="attributes">
			  <xs:complexType>
				<xs:sequence>
				  <xs:element type="xs:string" name="name"/>
				  <xs:element type="xs:string" name="color"/>
				  <xs:element type="xs:integer" name="age"/>
				</xs:sequence>
			  </xs:complexType>
			</xs:element>
		  </xs:sequence>
		</xs:complexType>
	  </xs:element>
	</xs:schema>`)

	var headers = make(map[string]string)
	headers["Content-Type"] = "application/xml"

	schemaFile := "test/testdata/imposters/schemas/create_gopher_request.xsd"
	bodyFile := "test/testdata/imposters/responses/create_gopher_response.xml"
	bodyFileFake := "test/testdata/imposters/responses/create_gopher_response_fail.xml"
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
		{"valid xml imposter with body", Imposter{Request: validRequest, Response: Response{Status: http.StatusOK, Headers: &headers, Body: body}}, body, http.StatusOK},
		{"valid xml imposter with bodyFile", Imposter{Request: validRequest, Response: Response{Status: http.StatusOK, Headers: &headers, BodyFile: &bodyFile}}, string(expectedBodyFileData), http.StatusOK},
		{"valid xml imposter with not exists bodyFile", Imposter{Request: validRequest, Response: Response{Status: http.StatusOK, Headers: &headers, BodyFile: &bodyFileFake}}, "", http.StatusOK},
	}

	for _, tt := range dataTest {
		t.Run(tt.name, func(t *testing.T) {
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
