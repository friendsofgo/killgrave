package http

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

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
			},
			"friends": [
				{
					"name": "Zebediah",
					"color": "Purple",
					"age": 55
				}
			]
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

	f, err := os.Open(bodyFile)
	require.NoError(t, err)
	defer f.Close()

	expectedBody := `{
    "data": {
        "type": "gophers",
        "id": "bca49e8a-82dd-4c5d-b886-13a6ceb3744b",
        "timestamp": "` + time.Now().UTC().Format("2006-01-02 15:04") + `",
        "birthday": "` + time.Now().UTC().Add(time.Hour*24).Format("2006-01-02") + `",
        "attributes": {
            "name": "Natalissa",
            "color": "Blue,Purple",
            "age": 42
        },
        "friends": [{"age":55,"color":"Purple","name":"Zebediah"}]
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
			require.NoError(t, err)
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

func TestExtractPathParams(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		url      string
		expected map[string]string
	}{
		{
			name:     "simple path params",
			endpoint: "/gophers/{id}",
			url:      "/gophers/123",
			expected: map[string]string{
				"id": "123",
			},
		},
		{
			name:     "no path params",
			endpoint: "/gophers",
			url:      "/gophers",
			expected: map[string]string{},
		},
		{
			name:     "mismatched path parts",
			endpoint: "/gophers/{id}/details",
			url:      "/gophers/123",
			expected: map[string]string{},
		},
		{
			name:     "colon path params",
			endpoint: "/gophers/:id",
			url:      "/gophers/123",
			expected: map[string]string{
				"id": "123",
			},
		},
		{
			name:     "multiple path params",
			endpoint: "/gophers/{id}/friends/{friendID}",
			url:      "/gophers/123/friends/456",
			expected: map[string]string{
				"id":       "123",
				"friendID": "456",
			},
		},
		{
			name:     "URL with extra slashes",
			endpoint: "/gophers/{id}",
			url:      "/gophers//123",
			expected: map[string]string{},
		},
		{
			name:     "URL with special characters",
			endpoint: "/gophers/{id}",
			url:      "/gophers/123@!$",
			expected: map[string]string{
				"id": "123@!$",
			},
		},
		{
			name:     "URL with missing path parts",
			endpoint: "/gophers/{id}/friends/{friendID}",
			url:      "/gophers/123/friends",
			expected: map[string]string{},
		},
		{
			name:     "URL with query parameters",
			endpoint: "/gophers/{id}",
			url:      "/gophers/123?color=blue",
			expected: map[string]string{
				"id": "123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			params := extractPathParams(req, tt.endpoint)
			if len(params) != len(tt.expected) {
				t.Errorf("expected %d params, got %d", len(tt.expected), len(params))
			}
			for k, v := range tt.expected {
				if params[k] != v {
					t.Errorf("expected param %s to be %s, got %s", k, v, params[k])
				}
			}
		})
	}
}

func TestExtractQueryParams(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected map[string][]string
	}{
		{
			name: "simple query params",
			url:  "/gophers?color=blue&age=5",
			expected: map[string][]string{
				"color": {"blue"},
				"age":   {"5"},
			},
		},
		{
			name:     "no query params",
			url:      "/gophers",
			expected: map[string][]string{},
		},
		{
			name: "multiple values for a query param",
			url:  "/gophers?color=blue&color=green",
			expected: map[string][]string{
				"color": {"blue", "green"},
			},
		},
		{
			name: "empty query parameter",
			url:  "/gophers?color=",
			expected: map[string][]string{
				"color": {""},
			},
		},
		{
			name: "query parameter with no value",
			url:  "/gophers?color",
			expected: map[string][]string{
				"color": {""},
			},
		},
		{
			name: "query parameter with special characters",
			url:  "/gophers?color=blue&name=John%20Doe",
			expected: map[string][]string{
				"color": {"blue"},
				"name":  {"John Doe"},
			},
		},
		{
			name: "query parameter with mixed case",
			url:  "/gophers?Color=blue&color=green",
			expected: map[string][]string{
				"Color": {"blue"},
				"color": {"green"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			params := extractQueryParams(req)
			if len(params) != len(tt.expected) {
				t.Errorf("expected %d params, got %d", len(tt.expected), len(params))
			}
			for k, v := range tt.expected {
				if len(params[k]) != len(v) {
					t.Errorf("expected %d values for param %s, got %d", len(v), k, len(params[k]))
				}
				for i := range v {
					if params[k][i] != v[i] {
						t.Errorf("expected param %s to be %s, got %s", k, v[i], params[k][i])
					}
				}
			}
		})
	}
}
