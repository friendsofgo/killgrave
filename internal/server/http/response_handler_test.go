package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRepeatingResponse(t *testing.T) {
	var serverData = []struct {
		name           string
		imposter       Imposter
		expectedBodies []string
	}{
		{
			"repeating response with burst",
			Imposter{
				Request:   Request{Method: "GET", Endpoint: "/burst", ResponseMode: "BURST"},
				Responses: []Response{{Status: 201, Body: "Response 1", Burst: 1}, {Status: 201, Body: "Response 2", Burst: 2}},
			},
			[]string{"Response 1", "Response 2", "Response 2", "Response 1", "Response 2", "Response 2", "Response 1"},
		},
		{
			"repeating response without burst", // Default value checking
			Imposter{
				Request:   Request{Method: "GET", Endpoint: "/repeat", ResponseMode: "BURST"},
				Responses: []Response{{Status: 201, Body: "Response 1"}, {Status: 201, Body: "Response 2"}},
			},
			[]string{"Response 1", "Response 2", "Response 1", "Response 2", "Response 1"},
		},
	}

	for _, tt := range serverData {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.imposter.Request.Method, tt.imposter.Request.Endpoint, bytes.NewBuffer(nil))
			if err != nil {
				t.Fatalf("could not created request: %v", err)
			}
			rec := httptest.NewRecorder()
			handler := http.HandlerFunc(ImposterHandler(tt.imposter))

			for i := 0; i < len(tt.expectedBodies); i++ {
				handler.ServeHTTP(rec, req)
				expectedBody := tt.expectedBodies[i]
				actualBody := rec.Body.String()
				if expectedBody != actualBody {
					t.Fatalf("test-%s expected body is '%s' and got '%s'", tt.name, expectedBody, actualBody)
				}
				rec.Body.Reset()
			}
		})
	}
}

// Checking if valid responses are being generated.
func TestRandomResponse(t *testing.T) {
	var serverData = []struct {
		name     string
		imposter Imposter
	}{
		{
			"random responses",
			Imposter{
				Request:   Request{Method: "GET", Endpoint: "/burst", ResponseMode: "RANDOM"},
				Responses: []Response{{Status: 201, Body: "Response 1"}, {Status: 201, Body: "Response 2"}},
			},
		},
		{
			"random responses with burst",
			Imposter{
				Request:   Request{Method: "GET", Endpoint: "/repeat", ResponseMode: "RANDOM"},
				Responses: []Response{{Status: 201, Body: "Response 1", Burst: 1}, {Status: 201, Body: "Response 2", Burst: 2}},
			},
		},
		{
			"random responses without response mode", // Default value checking
			Imposter{
				Request:   Request{Method: "GET", Endpoint: "/repeat"},
				Responses: []Response{{Status: 201, Body: "Response 1"}, {Status: 201, Body: "Response 2"}},
			},
		},
		{
			"random responses with more than 2 responses",
			Imposter{
				Request:   Request{Method: "GET", Endpoint: "/repeat", ResponseMode: "BURST"},
				Responses: []Response{{Status: 201, Body: "Response 1"}, {Status: 201, Body: "Response 2"}, {Status: 201, Body: "Response 3"}, {Status: 201, Body: "Response 4"}},
			},
		},
	}

	for _, tt := range serverData {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.imposter.Request.Method, tt.imposter.Request.Endpoint, bytes.NewBuffer(nil))
			if err != nil {
				t.Fatalf("could not created request: %v", err)
			}
			rec := httptest.NewRecorder()
			handler := http.HandlerFunc(ImposterHandler(tt.imposter))
			expectedRespMap := map[string]bool{}
			for i := 0; i < len(tt.imposter.Responses); i++ {
				expectedRespMap[tt.imposter.Responses[i].Body] = true
			}

			for i := 0; i < 10; i++ {
				handler.ServeHTTP(rec, req)
				actualBody := rec.Body.String()

				if _, ok := expectedRespMap[actualBody]; !ok {
					t.Fatalf("test-%s invalid response body: '%s'", tt.name, actualBody)
				}

				rec.Body.Reset()
			}
		})
	}
}
