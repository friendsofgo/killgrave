package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gorilla/handlers"
)

var binaryContentTypes = []string{
	"application/octet-stream",
	"image/",
	"audio/",
	"video/",
	"application/pdf",
}

// RequestData struct to hold request data
type RequestData struct {
	Method string      `json:"method"`
	Host   string      `json:"host"`
	URL    string      `json:"url"`
	Header http.Header `json:"header"`
	Status int         `json:"statusCode,omitempty"`
	Body   string      `json:"body,omitempty"`
}

func getRequestData(r *http.Request, status int, body string) *RequestData {
	return &RequestData{
		Method: r.Method,
		Host:   r.Host,
		URL:    r.URL.String(),
		Header: r.Header,
		Status: status,
		Body:   body,
	}
}

// ToJSON converts the RequestData struct to JSON.
func (rd *RequestData) toJSON() ([]byte, error) {
	jsonData, err := json.Marshal(rd)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

// isBinaryContent checks to see if the body is a common binary content type
func isBinaryContent(r *http.Request) bool {
	contentType := r.Header.Get("Content-Type")
	for _, binaryType := range binaryContentTypes {
		if strings.HasPrefix(contentType, binaryType) {
			return true
		}
	}
	return false
}

func shouldRecordRequest(s *Server) bool {
	return s.serverCfg.LogWriter != nil && s.dumpCh != nil
}

func getBody(r *http.Request, s *Server) string {
	if s.serverCfg.LogLevel == 0 && !shouldRecordRequest(s) {
		return ""
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v\n", err)
		return ""
	}
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes)) // Reset the body

	// trim if larger than the limit allowed
	if len(bodyBytes) > s.serverCfg.LogBodyMax {
		bodyBytes = bodyBytes[:s.serverCfg.LogBodyMax]
	}

	body := base64.StdEncoding.EncodeToString(bodyBytes)
	// if content is not binary, get it as a string
	if !isBinaryContent(r) {
		body = string(bodyBytes)
	}
	return body
}

// CustomLoggingHandler provides a way to supply a custom log formatter
// while taking advantage of the mechanisms in this package
func CustomLoggingHandler(out io.Writer, h http.Handler, s *Server) http.Handler {
	return handlers.CustomLoggingHandler(out, h, func(writer io.Writer, params handlers.LogFormatterParams) {
		body := getBody(params.Request, s)
		requestData := getRequestData(params.Request, params.StatusCode, body)

		// log the request
		if s.serverCfg.LogLevel > 0 {
			// if we add other formats to log in we can switch here and return the bytes
			data, err := requestData.toJSON()
			if err != nil {
				log.Printf("Error encoding request data: %+v\n", err)
				return
			}

			data = append(data, '\n')
			writer.Write(data)
			if shouldRecordRequest(s) {
				recordRequest(&data, s)
			}
		}
	})
}

func recordRequest(request *[]byte, s *Server) {
	select {
	case s.dumpCh <- request:
		// Successfully sent the request data to the channel
	default:
		// Handle the case where the channel is full
		log.Println("request dump channel is full, could not write request")
	}
}

// Goroutine function to write requests to a JSON file
func RequestWriter(ctx context.Context, wg *sync.WaitGroup, writer io.Writer, requestChan <-chan *[]byte) {
	defer wg.Done()

	for {
		select {
		case requestData := <-requestChan:
			if requestData == nil {
				return // channel closed
			}

			writer.Write(*requestData)
			// call Sync if writer is *os.File
			if f, ok := writer.(*os.File); ok {
				f.Sync()
			}
		case <-ctx.Done():
			return // context cancelled
		}
	}
}
