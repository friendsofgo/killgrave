package v2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/friendsofgo/killgrave/internal/server/http"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
	"sort"
	"strconv"
	"strings"
)

func GenerateSwaggerYAML(swagger []byte) (*[]http.Imposter, error) {
	var err error

	var doc openapi2.T

	if !bytes.Contains(swagger, []byte("swagger: \"2.0\"")) {
		return nil, fmt.Errorf("swagger file does not appear to be swagger 2.0")
	}

	if err = yaml.Unmarshal(swagger, &doc); err != nil {
		return nil, err
	}

	return generateInner(doc)
}

func GenerateSwaggerJSON(swagger []byte) (*[]http.Imposter, error) {
	var err error

	var doc openapi2.T

	if !bytes.Contains(swagger, []byte("\"swagger\": \"2.0\"")) &&
		!bytes.Contains(swagger, []byte("\"swagger\":\"2.0\"")) {
		return nil, fmt.Errorf("swagger file does not appear to be swagger 2.0")
	}

	if err = json.Unmarshal(swagger, &doc); err != nil {
		return nil, err
	}

	return generateInner(doc)
}

func generateInner(doc openapi2.T) (*[]http.Imposter, error) {
	var result []http.Imposter

	paths := make([]string, 0)
	for k, _ := range doc.Paths {
		paths = append(paths, k)
	}
	sort.Strings(paths)

	for _, pathKey := range paths {
		path := doc.Paths[pathKey]

		for method, operation := range path.Operations() {
			request, err := extractRequest(path, method); if err != nil {
				return nil, err
			}

			keys := make([]string, 0)
			for k, _ := range operation.Responses {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, key := range keys {
				response := operation.Responses[key]

				var numericResponseCode int

				if key == "default" {
					numericResponseCode = 500
				} else {
					numericResponseCode, err = strconv.Atoi(key); if err != nil {
						return nil, fmt.Errorf("unknown response code: %v", key)
					}
				}

				body, err := stringifyBody(response.Schema); if err != nil {
					return nil, fmt.Errorf("unable to stringify body: %w", err)
				}

				result = append(result, http.Imposter{
					BasePath: "",
					Request:  http.Request{
						Method:     method,
						Endpoint:   pathKey,
						SchemaFile: nil,
						Params:     convertRequestParameters(request.Parameters),
					},
					Response: http.Response{
						Status:   numericResponseCode,
						Body:     body,
						BodyFile: nil,
						Headers:  convertResponseHeaders(response.Headers),
						Delay:    http.NewResponseDelay(0, 0),
					},
				})
			}
		}
	}

	return &result, nil
}

func extractRequest(pathItem *openapi2.PathItem, method string) (*openapi2.Operation, error) {
	switch strings.ToUpper(method) {
	case "GET":
		return pathItem.Get, nil
	case "POST":
		return pathItem.Post, nil
	case "PUT":
		return pathItem.Put, nil
	case "PATCH":
		return pathItem.Patch, nil
	case "DELETE":
		return pathItem.Delete, nil
	case "HEAD":
		return pathItem.Head, nil
	case "OPTIONS":
		return pathItem.Options, nil
	default:
		return nil, fmt.Errorf("unknown operation")
	}
}

func convertRequestParameters(parameters openapi2.Parameters) *map[string]string {
	if parameters == nil || len(parameters) == 0 {
		return nil
	}

	result := map[string]string{}

	for _, param := range parameters {
		result[param.Name] = param.Description
	}

	return &result
}

func convertResponseHeaders(value map[string]*openapi2.Header) *map[string]string {
	if value == nil || len(value) == 0 {
		return nil
	}

	result := map[string]string{}

	for key, header := range value {
		result[key] = header.Description
	}

	return &result
}

func stringifyBody(content *openapi3.SchemaRef) (string, error) {
	if content == nil || content.Value == nil {
		return "", nil
	}

	props := map[string]interface{}{}
	for key, prop := range content.Value.Properties {
		props[key] = prop.Value.Example
	}

	if marshalledProps, err := json.Marshal(props); err != nil {
		return "", err
	} else {
		return string(marshalledProps), nil
	}
}
