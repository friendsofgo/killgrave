package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	internal_handlers "github.com/friendsofgo/killgrave/internal/gorilla/handlers"
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
	Body   string      `json:"body"`
}

func getRequestData(r *http.Request, body string) *RequestData {
	return &RequestData{
		Method: r.Method,
		Host:   r.Host,
		URL:    r.URL.String(),
		Header: r.Header,
		Body:   body,
	}
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
	body := string(bodyBytes)
	// if content is binary, encode it to base64
	if isBinaryContent(r) {
		body = base64.StdEncoding.EncodeToString(bodyBytes)
	}
	return body
}

// CustomLoggingHandler provides a way to supply a custom log formatter
// while taking advantage of the mechanisms in this package
func CustomLoggingHandler(out io.Writer, h http.Handler, s *Server) http.Handler {
	return handlers.CustomLoggingHandler(out, h, func(writer io.Writer, params handlers.LogFormatterParams) {
		body := getBody(params.Request, s)

		var err error
		if shouldRecordRequest(s) {
			err = recordRequest(params.Request, s, body)
		}

		// log the request based on the log level
		// if err is set, log the request, but only add the body if the log level is 2 or higher
		if s.serverCfg.LogLevel >= 2 {
			internal_handlers.WriteLog(writer, params, body)
		} else if err != nil || s.serverCfg.LogLevel > 0 {
			internal_handlers.WriteLog(writer, params, "")
		}
	})
}

func recordRequest(r *http.Request, s *Server, body string) error {
	rd := getRequestData(r, body)
	select {
	case s.dumpCh <- rd:
		// Successfully sent the request data to the channel
	default:
		// Handle the case where the channel is full
		log.Println("Channel is full, dropping request and logging it instead:")
		return fmt.Errorf("request dump channel is full")
	}
	return nil
}

// Goroutine function to write requests to a JSON file
func RequestWriter(ctx context.Context, wg *sync.WaitGroup, writer io.Writer, requestChan <-chan *RequestData) {
	defer wg.Done()
	// defer file.Close()

	encoder := json.NewEncoder(writer)
	for {
		select {
		case requestData := <-requestChan:
			if requestData == nil {
				return // channel closed
			}

			if err := encoder.Encode(requestData); err != nil {
				log.Printf("Failed to write to file: %+v", err)
				fmt.Printf("Failed to write to file: %+v", err)
			}
			// Type assertion to call Sync if writer is *os.File
			if f, ok := writer.(*os.File); ok {
				f.Sync()
			}
		case <-ctx.Done():
			return // context cancelled
		}
	}
}
