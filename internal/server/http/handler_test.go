package http

import (
	"bytes"
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

func TestImposterHandlerTemplating(t *testing.T) {
	bodyRequest := []byte(`{
		"data": {
		  "type": "gophers",
		  "attributes": {
			"name": "Natalissa"
		  }
		}
	  }`)
	var headers = make(map[string]string)
	headers["Content-Type"] = "application/json"

	schemaFile := "test/testdata/imposters_templating/schemas/create_gopher_request.json"
	bodyFile := "test/testdata/imposters_templating/responses/create_gopher_response.json.tmpl"
	bodyFileFake := "test/testdata/imposters_templating/responses/create_gopher_response_fail.json"
	body := `{"test":true}`

	validRequest := Request{
		Method:     "POST",
		Endpoint:   "/gophers/{GopherID}",
		SchemaFile: &schemaFile,
		Headers:    &headers,
	}

	f, _ := os.Open(bodyFile)
	defer f.Close()

	expectedBody := `{
    "data": {
        "type": "gophers",
        "id": "bca49e8a-82dd-4c5d-b886-13a6ceb3744b",
        "attributes": {
            "name": "Natalissa",
            "color": "Blue,Purple",
            "age": 42
        }
    }
}
`

	var dataTest = []struct {
		name         string
		imposter     Imposter
		expectedBody string
		statusCode   int
	}{
		{"valid imposter with body", Imposter{Request: validRequest, Response: Responses{{Status: http.StatusOK, Headers: &headers, Body: body}}}, body, http.StatusOK},
		{"valid imposter with bodyFile", Imposter{Request: validRequest, Response: Responses{{Status: http.StatusOK, Headers: &headers, BodyFile: &bodyFile}}}, expectedBody, http.StatusOK},
		{"valid imposter with not exists bodyFile", Imposter{Request: validRequest, Response: Responses{{Status: http.StatusOK, Headers: &headers, BodyFile: &bodyFileFake}}}, "", http.StatusOK},
	}

	for _, tt := range dataTest {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/gophers/bca49e8a-82dd-4c5d-b886-13a6ceb3744b?gopherColor=Blue&gopherColor=Purple&gopherAge=42", bytes.NewBuffer(bodyRequest))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			handler := ImposterHandler(tt.imposter)

			handler.ServeHTTP(rec, req)
			assert.Equal(t, rec.Code, tt.statusCode)
			assert.Equal(t, tt.expectedBody, rec.Body.String())

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
