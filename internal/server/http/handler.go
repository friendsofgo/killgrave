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
	"text/template"
	"time"
)

type TemplatingData struct {
	RequestBody map[string]interface{}
	PathParams  map[string]string
	QueryParams map[string][]string
}

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

	templateBytes, err := applyTemplate(i, bodyBytes, r)
	if err != nil {
		log.Printf("error applying template: %v\n", err)
	}

	w.Write(templateBytes)
}

func fetchBodyFromFile(bodyFile string) (bytes []byte) {
	if _, err := os.Stat(bodyFile); os.IsNotExist(err) {
		log.Printf("the body file %s not found\n", bodyFile)
		return
	}

	f, _ := os.Open(bodyFile)
	defer f.Close()
	bytes, err := io.ReadAll(f)
	if err != nil {
		log.Printf("imposible read the file %s: %v\n", bodyFile, err)
	}
	return
}

func applyTemplate(i Imposter, bodyBytes []byte, r *http.Request) ([]byte, error) {
	bodyStr := string(bodyBytes)

	// check if the body contains a template
	if !strings.Contains(bodyStr, "{{") {
		return bodyBytes, nil
	}

	tmpl, err := template.New("body").
		Funcs(template.FuncMap{"stringsJoin": strings.Join}).
		Parse(bodyStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %w", err)
	}

	extractedBody, err := extractBody(r)
	if err != nil {
		log.Printf("error extracting body: %v\n", err)
	}

	// parse request body in a generic way
	tmplData := TemplatingData{
		RequestBody: extractedBody,
		PathParams:  extractPathParams(i, r),
		QueryParams: extractQueryParams(r),
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, tmplData)
	if err != nil {
		return nil, fmt.Errorf("error applying template: %w", err)
	}

	return tpl.Bytes(), nil
}

func extractBody(r *http.Request) (map[string]interface{}, error) {
	body := make(map[string]interface{})
	if r.Body == nil {
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

func extractPathParams(i Imposter, r *http.Request) map[string]string {
	params := make(map[string]string)

	path := r.URL.Path
	if path == "" {
		return params
	}

	endpoint := i.Request.Endpoint
	// regex to split either path params using /:paramname or /{paramname}

	// split path and endpoint by /
	pathParts := strings.Split(path, "/")
	imposterParts := strings.Split(endpoint, "/")

	if len(pathParts) != len(imposterParts) {
		log.Printf("request path and imposter endpoint parts do not match: %s, %s\n", path, endpoint)
		return params
	}

	// iterate over pathParts and endpointParts
	for i := range imposterParts {
		if strings.HasPrefix(imposterParts[i], ":") {
			params[imposterParts[i][1:]] = pathParts[i]
		}
		if strings.HasPrefix(imposterParts[i], "{") && strings.HasSuffix(imposterParts[i], "}") {
			params[imposterParts[i][1:len(imposterParts[i])-1]] = pathParts[i]
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
