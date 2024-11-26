package templating

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyTemplate(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name      string
		bodyStr   string
		templData TemplatingData
		expected  string
	}{
		{
			name:    "simple template",
			bodyStr: `{"message": "Hello, {{ .PathParams.name }}!"}`,
			templData: TemplatingData{
				PathParams: map[string]string{
					"name": "World",
				},
			},
			expected: `{"message": "Hello, World!"}`,
		},
		{
			name:    "template with JSON marshaling",
			bodyStr: `{"data": {{ jsonMarshal .RequestBody.data }}}`,
			templData: TemplatingData{
				RequestBody: map[string]interface{}{
					"data": map[string]string{
						"key": "value",
					},
				},
			},
			expected: `{"data": {"key":"value"}}`,
		},
		{
			name: "template with time functions",
			bodyStr: `{
				"timestamp": "{{ timeFormat (timeUTC (timeNow)) "2006-01-02 15:04" }}",
				"future": "{{ timeFormat (timeAdd (timeUTC (timeNow)) "24h") "2006-01-02" }}"
			}`,
			templData: TemplatingData{
				RequestBody: map[string]interface{}{},
				PathParams:  map[string]string{},
				QueryParams: map[string][]string{},
			},
			expected: `{
				"timestamp": "` + now.UTC().Format("2006-01-02 15:04") + `",
				"future": "` + now.Add(24*time.Hour).UTC().Format("2006-01-02") + `"
			}`,
		},
		{
			name:    "template with string join",
			bodyStr: `{"colors": "{{ stringsJoin .QueryParams.colors "," }}"}`,
			templData: TemplatingData{
				QueryParams: map[string][]string{
					"colors": {"red", "green", "blue"},
				},
			},
			expected: `{"colors": "red,green,blue"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ApplyTemplate(tt.bodyStr, tt.templData)
			require.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(got))
		})
	}
}

func TestApplyTemplate_Error(t *testing.T) {
	tests := []struct {
		name      string
		bodyStr   string
		templData TemplatingData
		errMsg    string
	}{
		{
			name:    "invalid template directive",
			bodyStr: `{"message": "Hello, {{ .PathParams.name | invalidFunc }}`,
			templData: TemplatingData{
				PathParams: map[string]string{
					"name": "World",
				},
			},
			errMsg: "function \"invalidFunc\" not defined",
		},
		{
			name:    "invalid template",
			bodyStr: `{"message": "Hello, {{ .InvalidField }}`,
			templData: TemplatingData{
				PathParams: map[string]string{
					"name": "World",
				},
			},
			errMsg: "error applying template",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ApplyTemplate(tt.bodyStr, tt.templData)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}
