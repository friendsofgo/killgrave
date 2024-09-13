package http

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	expectedBodyFileData, _ := io.ReadAll(f)

	var dataTest = []struct {
		name         string
		imposter     Imposter
		expectedBody string
		statusCode   int
	}{
		{"valid imposter with body", Imposter{Request: validRequest, Response: Responses{{Status: http.StatusOK, Headers: &headers, Body: body}}}, body, http.StatusOK},
		{"valid imposter with bodyFile", Imposter{Request: validRequest, Response: Responses{{Status: http.StatusOK, Headers: &headers, BodyFile: &bodyFile}}}, string(expectedBodyFileData), http.StatusOK},
		{"valid imposter with not exists bodyFile", Imposter{Request: validRequest, Response: Responses{{Status: http.StatusOK, Headers: &headers, BodyFile: &bodyFileFake}}}, "", http.StatusOK},
	}

	for _, tt := range dataTest {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/gophers", bytes.NewBuffer(bodyRequest))
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			handler := ImposterHandler(tt.imposter)

			handler.ServeHTTP(rec, req)
			assert.Equal(t, rec.Code, tt.statusCode)
			assert.Equal(t, tt.expectedBody, rec.Body.String())

		})
	}
}

func TestImposterHandler_Variables(t *testing.T) {
	var headers = make(map[string]string)
	headers["Content-Type"] = "application/json"

	responseId1 := "test/testdata/imposters_variables/responses/gopher_1_response.json"
	responseId2 := "test/testdata/imposters_variables/responses/gopher_2_response.json"
	responseId1Variable1 := "test/testdata/imposters_variables/responses/gopher_1_1_response.json"
	responseId1Variable2 := "test/testdata/imposters_variables/responses/gopher_1_2_response.json"
	responseId1WithoutVariable := "test/testdata/imposters_variables/responses/gopher_1_without_variable_response.json"

	imposterFilePath := "test/testdata/imposters_variables/gopher_variables.imp.json"
	imposterFile, err := os.Open(imposterFilePath)
	require.NoError(t, err)
	defer imposterFile.Close()

	imposterBytes, err := io.ReadAll(imposterFile)
	require.NoError(t, err)

	var imposters []Imposter
	err = json.Unmarshal(imposterBytes, &imposters)
	require.NoError(t, err)

	var dataTest = []struct {
		name             string
		imposter         Imposter
		url              string
		expectedBodyPath string
		statusCode       int
	}{
		{"valid imposter with id 1 in path", imposters[0], "/gophers/1", responseId1, http.StatusOK},
		{"valid imposter with id 2 in path", imposters[0], "/gophers/2", responseId2, http.StatusOK},
		{"valid imposter with id 1 and second variable 1 in path", imposters[1], "/gophers/1/1", responseId1Variable1, http.StatusOK},
		{"valid imposter with id 1 and second variable 2 in path", imposters[1], "/gophers/1/2", responseId1Variable2, http.StatusOK},
		{"valid imposter without variable but body file has variable", imposters[2], "/gophers/1", responseId1WithoutVariable, http.StatusOK},
	}

	for _, tt := range dataTest {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			require.NoError(t, err)
			rec := httptest.NewRecorder()
			handler := ImposterHandler(tt.imposter)

			m := mux.NewRouter()
			m.Handle(tt.imposter.Request.Endpoint, handler)
			m.ServeHTTP(rec, req)

			expectedBodyPathFile, _ := os.Open(tt.expectedBodyPath)
			defer expectedBodyPathFile.Close()
			expectedBody, err := io.ReadAll(expectedBodyPathFile)
			require.NoError(t, err)

			assert.Equal(t, rec.Code, tt.statusCode)
			assert.Equal(t, string(expectedBody), rec.Body.String())

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
		{"valid request no schema", Imposter{Request: Request{Method: "POST", Endpoint: "/gophers"}, Response: Responses{{Status: http.StatusOK, Body: "test ok"}}}, http.StatusOK, validRequest},
	}

	for _, tt := range dataTest {

		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/gophers", bytes.NewBuffer(tt.request))
			assert.Nil(t, err)
			rec := httptest.NewRecorder()
			handler := ImposterHandler(tt.imposter)

			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.statusCode, rec.Code)
		})
	}
}

func TestImposterHandler_MultipleRequests(t *testing.T) {
	req, err := http.NewRequest("POST", "/gophers", bytes.NewBuffer([]byte(`{
		"data": {
			"type": "gophers",
		  "attributes": {
			"name": "Zebediah",
			"color": "Purple"
		  }
		}
	  }`)))
	require.NoError(t, err)

	t.Run("created then conflict", func(t *testing.T) {
		imp := Imposter{
			Request: Request{Method: "POST", Endpoint: "/gophers"},
			Response: Responses{
				{Status: http.StatusCreated, Body: "Created"},
				{Status: http.StatusConflict, Body: "Conflict"},
			},
		}

		handler := ImposterHandler(imp)

		// First request
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "Created", rec.Body.String())

		// Second request
		rec = httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusConflict, rec.Code)
		assert.Equal(t, "Conflict", rec.Body.String())
	})

	t.Run("idempotent", func(t *testing.T) {
		handler := ImposterHandler(Imposter{
			Request: Request{Method: "POST", Endpoint: "/gophers"},
			Response: Responses{
				{Status: http.StatusAccepted, Body: "Accepted"},
			},
		})

		// First request
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusAccepted, rec.Code)
		assert.Equal(t, "Accepted", rec.Body.String())

		// Second request
		rec = httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusAccepted, rec.Code)
		assert.Equal(t, "Accepted", rec.Body.String())
	})
}
