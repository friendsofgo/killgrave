package v3

import (
	"encoding/json"
	"fmt"
	"github.com/friendsofgo/killgrave/internal/server/http"
	"github.com/getkin/kin-openapi/openapi3"
	"sort"
	"strconv"
	"strings"
)

func GenerateSwagger(swagger []byte) (*[]http.Imposter, error) {
	var result []http.Imposter

	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(swagger)
	if err != nil {
		return nil, err
	}

	err = doc.Validate(loader.Context)
	if err != nil {
		return nil, err
	}

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

				body, err := stringifyBody(response.Value); if err != nil {
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
						Headers:  convertResponseHeaders(response.Value),
						Delay:    http.NewResponseDelay(0, 0),
					},
				})
			}
		}
	}

	return &result, nil
}

func convertRequestParameters(parameters openapi3.Parameters) *map[string]string {
	if parameters == nil {
		return nil
	}

	result := map[string]string{}

	for _, param := range parameters {
		result[param.Value.Name] = param.Value.Description
	}

	return &result
}

func convertResponseHeaders(value *openapi3.Response) *map[string]string {
	if value == nil || len(value.Headers) == 0 {
		return nil
	}

	result := map[string]string{}

	for _, header := range value.Headers {
		result[header.Value.Name] = header.Value.Description
	}

	return &result
}

func stringifyBody(response *openapi3.Response) (string, error) {
	if response == nil || response.Content == nil {
		return "", nil
	}

	if val := response.Content.Get("application/json"); val != nil {
		if val.Schema == nil || val.Schema.Value == nil {
			return "", nil
		}

		props := map[string]interface{}{}
		for key, prop := range val.Schema.Value.Properties {
			if prop.Value != nil {
				props[key] = prop.Value.Example
			}
		}

		if marshalledProps, err := json.Marshal(props); err != nil {
			return "", err
		} else {
			return string(marshalledProps), nil
		}
	}

	return "", nil
}

func extractRequest(pathItem *openapi3.PathItem, method string) (*openapi3.Operation, error) {
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
	case "TRACE":
		return pathItem.Trace, nil
	default:
		return nil, fmt.Errorf("unknown operation")
	}
}
