package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

// RequestData struct to hold request data
type RequestData struct {
	Method string              `json:"method"`
	Host   string              `json:"host"`
	URL    string              `json:"url"`
	Header map[string][]string `json:"header"`
	Body   string              `json:"body"`
}

func GetRequestData(r *http.Request) *RequestData {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		return nil
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reset the body

	return &RequestData{
		Method: r.Method,
		Host:   r.Host,
		URL:    r.URL.String(),
		Header: r.Header,
		Body:   string(bodyBytes),
	}
}

func LogRequest(r *RequestData, s *Server) {
	if !s.verbose {
		return
	}

	// rebuild the request
	req, err := http.NewRequest(r.Method, r.URL, bytes.NewBufferString(r.Body))
	req.Host = r.Host
	if err != nil {
		log.Printf("failed to create request for logging: %+v", err)
		return
	}
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	dumped, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Printf("failed to dump request: %v", err)
	}
	if dumped != nil {
		log.Println(string(dumped))
	}

}

func RecordRequest(r *RequestData, s *Server) {
	if len(s.dumpRequestsPath) < 1 || s.dumpCh == nil {
		return
	}
	s.dumpCh <- r
}

// Goroutine function to write requests to a JSON file
func RequestWriter(filePath string, requestChan <-chan *RequestData) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %+v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for requestData := range requestChan {
		if err := encoder.Encode(requestData); err != nil {
			log.Printf("Failed to write to file: %+v", err)
		}
	}
}

// GetRecordedRequests reads the requests from the file and returns them as a slice of RequestData
func getRecordedRequests(filePath string) ([]RequestData, error) {
	// Read the file contents
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Split the contents by the newline separator
	requestDumps := strings.Split(string(fileContent), "\n")
	requestsData := []RequestData{}
	for _, requestDump := range requestDumps {
		if requestDump == "" {
			continue
		}
		// Unmarshal the JSON string into the RequestData struct
		rd := RequestData{}
		err := json.Unmarshal([]byte(requestDump), &rd)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
		}
		requestsData = append(requestsData, rd)
	}
	return requestsData, nil
}
