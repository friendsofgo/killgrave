package http

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

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

// Copied from the gorilla/handlers package
const lowerhex = "0123456789abcdef"

// Copied from the gorilla/handlers package
func appendQuoted(buf []byte, s string) []byte {
	var runeTmp [utf8.UTFMax]byte
	for width := 0; len(s) > 0; s = s[width:] {
		r := rune(s[0])
		width = 1
		if r >= utf8.RuneSelf {
			r, width = utf8.DecodeRuneInString(s)
		}
		if width == 1 && r == utf8.RuneError {
			buf = append(buf, `\x`...)
			buf = append(buf, lowerhex[s[0]>>4])
			buf = append(buf, lowerhex[s[0]&0xF])
			continue
		}
		if r == rune('"') || r == '\\' { // always backslashed
			buf = append(buf, '\\')
			buf = append(buf, byte(r))
			continue
		}
		if strconv.IsPrint(r) {
			n := utf8.EncodeRune(runeTmp[:], r)
			buf = append(buf, runeTmp[:n]...)
			continue
		}
		switch r {
		case '\a':
			buf = append(buf, `\a`...)
		case '\b':
			buf = append(buf, `\b`...)
		case '\f':
			buf = append(buf, `\f`...)
		case '\n':
			buf = append(buf, `\n`...)
		case '\r':
			buf = append(buf, `\r`...)
		case '\t':
			buf = append(buf, `\t`...)
		case '\v':
			buf = append(buf, `\v`...)
		default:
			switch {
			case r < ' ':
				buf = append(buf, `\x`...)
				buf = append(buf, lowerhex[s[0]>>4])
				buf = append(buf, lowerhex[s[0]&0xF])
			case r > utf8.MaxRune:
				r = 0xFFFD
				fallthrough
			case r < 0x10000:
				buf = append(buf, `\u`...)
				for s := 12; s >= 0; s -= 4 {
					buf = append(buf, lowerhex[r>>uint(s)&0xF])
				}
			default:
				buf = append(buf, `\U`...)
				for s := 28; s >= 0; s -= 4 {
					buf = append(buf, lowerhex[r>>uint(s)&0xF])
				}
			}
		}
	}
	return buf
}

// Copied from the gorilla/handlers package
func buildCommonLogLine(req *http.Request, url url.URL, ts time.Time, status int, size int) []byte {
	username := "-"
	if url.User != nil {
		if name := url.User.Username(); name != "" {
			username = name
		}
	}

	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		host = req.RemoteAddr
	}

	uri := req.RequestURI

	// Requests using the CONNECT method over HTTP/2.0 must use
	// the authority field (aka r.Host) to identify the target.
	// Refer: https://httpwg.github.io/specs/rfc7540.html#CONNECT
	if req.ProtoMajor == 2 && req.Method == "CONNECT" {
		uri = req.Host
	}
	if uri == "" {
		uri = url.RequestURI()
	}

	buf := make([]byte, 0, 3*(len(host)+len(username)+len(req.Method)+len(uri)+len(req.Proto)+50)/2)
	buf = append(buf, host...)
	buf = append(buf, " - "...)
	buf = append(buf, username...)
	buf = append(buf, " ["...)
	buf = append(buf, ts.Format("02/Jan/2006:15:04:05 -0700")...)
	buf = append(buf, `] "`...)
	buf = append(buf, req.Method...)
	buf = append(buf, " "...)
	buf = appendQuoted(buf, uri)
	buf = append(buf, " "...)
	buf = append(buf, req.Proto...)
	buf = append(buf, `" `...)
	buf = append(buf, strconv.Itoa(status)...)
	buf = append(buf, " "...)
	buf = append(buf, strconv.Itoa(size)...)
	return buf
}

// copied from gorilla/handlers and modified to add the body
// writeLog writes a log entry for req to w in Apache Common Log Format.
// ts is the timestamp with which the entry should be logged.
// status and size are used to provide the response HTTP status and size.
func writeLog(writer io.Writer, params handlers.LogFormatterParams, body string) {
	buf := buildCommonLogLine(params.Request, params.URL, params.TimeStamp, params.StatusCode, params.Size)
	// Append body if present
	if len(body) > 0 {
		buf = append(buf, " "...)
		buf = append(buf, body...)
	}
	buf = append(buf, '\n')
	writer.Write(buf)
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
	return len(s.dumpRequestsPath) > 0 && s.dumpCh != nil
}

func getBody(r *http.Request, s *Server) string {
	if s.logLevel == 0 && !shouldRecordRequest(s) {
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
		if s.logLevel >= 2 {
			writeLog(writer, params, body)
		} else if err != nil || s.logLevel > 0 {
			writeLog(writer, params, "")
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
		file.Sync()
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
