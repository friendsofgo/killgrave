package killgrave

import "net/http"

// Imposter define an imposter structure
type Imposter struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

// Request represent the structure of real request
type Request struct {
	Method     string       `json:"method"`
	Endpoint   string       `json:"endpoint"`
	SchemaFile *string      `json:"schema_file"`
	Headers    *http.Header `json:"headers"`
}

// Response represent the structure of real response
type Response struct {
	Status      int     `json:"status"`
	Body        string  `json:"body"`
	BodyFile    *string `json:"bodyFile"`
	ContentType string  `json:"content_type"`
}
