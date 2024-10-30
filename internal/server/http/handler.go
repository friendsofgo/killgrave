package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/friendsofgo/killgrave/internal/templating"
)

// ImposterHandler create specific handler for the received imposter
func ImposterHandler(i Imposter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := i.NextResponse()
		if res.Delay.Delay() > 0 {
			time.Sleep(res.Delay.Delay())
		}
		writeHeaders(res, w)
		w.WriteHeader(res.Status)
		writeBody(i, res, w, r)
	}
}

func writeHeaders(r Response, w http.ResponseWriter) {
	if r.Headers == nil {
		return
	}

	for key, val := range *r.Headers {
		w.Header().Set(key, val)
	}
}

func writeBody(i Imposter, res Response, w http.ResponseWriter, r *http.Request) {
	bodyBytes := []byte(res.Body)

	if res.BodyFile != nil {
		bodyFile := i.CalculateFilePath(*res.BodyFile)
		bodyBytes = fetchBodyFromFile(bodyFile)
	}

	bodyStr := string(bodyBytes)

	// early return if body does not contain templating
	if !strings.Contains(bodyStr, "{{") {
		w.Write([]byte(bodyStr))
		return
	}

	structuredBody, err := extractBody(r)
	if err != nil {
		log.Printf("error extracting body: %v\n", err)
	}

	templData := templating.TemplatingData{
		RequestBody: structuredBody,
		PathParams:  extractPathParams(r, i.Request.Endpoint),
		QueryParams: extractQueryParams(r),
	}

	templateBytes, err := templating.ApplyTemplate(bodyStr, templData)
	if err != nil {
		log.Printf("error applying template: %v\n", err)
	}

	w.Write(templateBytes)
}

func fetchBodyFromFile(bodyFile string) []byte {
	if _, err := os.Stat(bodyFile); os.IsNotExist(err) {
		log.Printf("the body file %s not found\n", bodyFile)
		return nil
	}

	f, _ := os.Open(bodyFile)
	defer f.Close()
	bytes, err := io.ReadAll(f)
	if err != nil {
		log.Printf("imposible read the file %s: %v\n", bodyFile, err)
		return nil
	}
	return bytes
}

func extractBody(r *http.Request) (map[string]interface{}, error) {
	body := make(map[string]interface{})
	if r.Body == http.NoBody {
		return body, nil
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return body, fmt.Errorf("error reading request body: %w", err)
	}

	// Restore the body for further use
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	contentType := r.Header.Get("Content-Type")

	switch {
	case strings.Contains(contentType, "application/json"):
		err = json.Unmarshal(bodyBytes, &body)
	default:
		return body, fmt.Errorf("unsupported content type: %s", contentType)
	}

	if err != nil {
		return body, fmt.Errorf("error unmarshaling request body: %w", err)
	}

	return body, nil
}

func extractPathParams(r *http.Request, endpoint string) map[string]string {
	params := make(map[string]string)

	path := r.URL.Path
	if path == "" {
		return params
	}

	// split path and endpoint by /
	pathParts := strings.Split(path, "/")
	endpointParts := strings.Split(endpoint, "/")

	if len(pathParts) != len(endpointParts) {
		log.Printf("request path and endpoint parts do not match: %s, %s\n", path, endpoint)
		return params
	}

	// iterate over pathParts and endpointParts
	for i := range endpointParts {
		if strings.HasPrefix(endpointParts[i], ":") {
			params[endpointParts[i][1:]] = pathParts[i]
		}
		if strings.HasPrefix(endpointParts[i], "{") && strings.HasSuffix(endpointParts[i], "}") {
			params[endpointParts[i][1:len(endpointParts[i])-1]] = pathParts[i]
		}
	}

	return params
}

func extractQueryParams(r *http.Request) map[string][]string {
	params := make(map[string][]string)
	query := r.URL.Query()
	for k, v := range query {
		params[k] = v
	}
	return params
}
