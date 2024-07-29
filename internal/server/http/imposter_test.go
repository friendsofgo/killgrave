package http

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestNewImposterFS(t *testing.T) {
	t.Run("imposters directory not found", func(t *testing.T) {
		_, err := NewImposterFS("failImposterPath")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "the directory 'failImposterPath' does not exist")
	})

	t.Run("existing imposters directory", func(t *testing.T) {
		_, err := NewImposterFS("test/testdata/imposters")
		assert.NoError(t, err)
	})
}

func TestImposterFS_FindImposters(t *testing.T) {
	// Set up
	const expected = 7
	ifs, err := NewImposterFS("test/testdata/imposters")
	require.NoError(t, err)

	// We trigger the imposters search.
	// We expect exactly [expected] imposters.
	ch := make(chan []Imposter, expected)
	err = ifs.FindImposters(ch)
	require.NoError(t, err)

	// We collect all the imposters.
	received := make([]Imposter, 0, expected)
	for ii := range ch {
		received = append(received, ii...)
	}
	require.Len(t, received, expected)

	// Imposter 1
	schemaFile := "schemas/create_gopher_request.json"
	bodyFile := "responses/create_gopher_response.json"
	assert.EqualValues(t, Imposter{
		BasePath: "test/testdata/imposters",
		Path:     "create_gopher.imp.json",
		Request: Request{
			Method:     "POST",
			Endpoint:   "/gophers",
			SchemaFile: &schemaFile,
			Params: &map[string]string{
				"gopherColor": "{v:[a-z]+}",
			},
			Headers: &map[string]string{
				"Content-Type": "application/json",
			},
		},
		Response: Responses{{
			Status: 200,
			Headers: &map[string]string{
				"Content-Type": "application/json",
			},
			BodyFile: &bodyFile,
		}},
	}, received[0])

	// Imposter 2
	assert.EqualValues(t, Imposter{
		BasePath: "test/testdata/imposters",
		Path:     "create_gopher.imp.json",
		Request:  Request{},
	}, received[1])

	// Imposter 3
	assert.EqualValues(t, Imposter{
		BasePath: "test/testdata/imposters",
		Path:     "test_request.imp.json",
		Request: Request{
			Method:   "GET",
			Endpoint: "/testRequest",
		},
		Response: Responses{{
			Status: 200,
			Body:   "Handled",
		}},
	}, received[2])

	// Imposter 4
	assert.EqualValues(t, Imposter{
		BasePath: "test/testdata/imposters",
		Path:     "test_request.imp.yaml",
		Request: Request{
			Method:   "GET",
			Endpoint: "/yamlTestRequest",
		},
		Response: Responses{{
			Status: 200,
			Body:   "Yaml Handled",
		}},
	}, received[3])

	// Imposter 5
	assert.EqualValues(t, Imposter{
		BasePath: "test/testdata/imposters",
		Path:     "test_request.imp.yml",
		Request: Request{
			Method:   "GET",
			Endpoint: "/ymlTestRequest",
		},
		Response: Responses{{
			Status: 200,
			Body:   "Yml Handled",
			Delay: ResponseDelay{
				delay:  1000000000,
				offset: 4000000000,
			},
		}},
	}, received[4])

	// Imposter 6
	assert.EqualValues(t, Imposter{
		BasePath: "test/testdata/imposters",
		Path:     "test_request.imp.yml",
		Request: Request{
			Method:   "POST",
			Endpoint: "/yamlGophers",
			Headers: &map[string]string{
				"Content-Type": "application/json",
			},
		},
		Response: Responses{{
			Status: 201,
			Headers: &map[string]string{
				"Content-Type": "application/json",
				"X-Source":     "YAML",
			},
			BodyFile: &bodyFile,
		}},
	}, received[5])

	// Imposter 7
	assert.EqualValues(t, Imposter{
		BasePath: "test/testdata/imposters",
		Path:     "test_request.imp.yml",
		Request:  Request{},
	}, received[6])

	// Finally, once the search is done,
	// the channel must be closed.
	_, open := <-ch
	require.False(t, open)
}

func TestResponses_MarshalJSON(t *testing.T) {
	tcs := map[string]struct {
		rr  *Responses
		exp string
	}{
		"single response": {
			rr:  &Responses{{Status: 200, Body: "OK"}},
			exp: `{"status":200,"body":"OK","bodyFile":null,"headers":null,"delay":{}}`,
		},
		"multiple response": {
			rr:  &Responses{{Status: 200, Body: "OK"}, {Status: 404, Body: "Not Found"}},
			exp: `[{"status":200,"body":"OK","bodyFile":null,"headers":null,"delay":{}},{"status":404,"body":"Not Found","bodyFile":null,"headers":null,"delay":{}}]`,
		},
		"empty array": {
			rr:  &Responses{},
			exp: `[]`,
		},
		"null array": {
			rr:  nil,
			exp: `null`,
		},
	}
	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			got, err := json.Marshal(tc.rr)
			assert.NoError(t, err)
			assert.Equal(t, tc.exp, string(got))
		})
	}
}

func TestResponses_UnmarshalJSON(t *testing.T) {
	tcs := map[string]struct {
		data string
		exp  Imposter
	}{
		"single response": {
			data: `{"response": {"status":200,"body":"OK"}}`,
			exp:  Imposter{Response: Responses{{Status: 200, Body: "OK"}}},
		},
		"single array response": {
			data: `{"response": [{"status":200,"body":"OK"}]}`,
			exp:  Imposter{Response: Responses{{Status: 200, Body: "OK"}}},
		},
		"multiple array response": {
			data: `{"response": [{"status":200,"body":"OK"}, {"status":404,"body":"Not Found"}]}`,
			exp:  Imposter{Response: Responses{{Status: 200, Body: "OK"}, {Status: 404, Body: "Not Found"}}},
		},
		"empty array": {
			data: `{"response": []}`,
			exp:  Imposter{Response: Responses{}},
		},
		"null array": {
			data: `{"response": null}`,
			exp:  Imposter{Response: nil},
		}}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			var got Imposter
			err := json.Unmarshal([]byte(tc.data), &got)
			require.NoError(t, err)
			assert.Equal(t, tc.exp, got)
		})
	}
}

func TestResponses_MarshalYAML(t *testing.T) {
	tcs := map[string]struct {
		rr  *Responses
		exp string
	}{
		"single response": {
			rr:  &Responses{{Status: 200, Body: "OK"}},
			exp: "status: 200\nbody: OK\nbodyFile: null\nheaders: null\ndelay: {}\n",
		},
		"multiple response": {
			rr:  &Responses{{Status: 200, Body: "OK"}, {Status: 404, Body: "Not Found"}},
			exp: "- status: 200\n  body: OK\n  bodyFile: null\n  headers: null\n  delay: {}\n- status: 404\n  body: Not Found\n  bodyFile: null\n  headers: null\n  delay: {}\n",
		},
		"empty array": {
			rr:  &Responses{},
			exp: "[]\n",
		},
		"null array": {
			rr:  nil,
			exp: "null\n",
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			got, err := yaml.Marshal(tc.rr)
			assert.NoError(t, err)
			assert.Equal(t, tc.exp, string(got))
		})
	}
}

func TestResponses_UnmarshalYAML(t *testing.T) {
	tcs := map[string]struct {
		data string
		exp  Imposter
	}{
		"single response": {
			data: "response:\n  status: 200\n  body: OK\n",
			exp:  Imposter{Response: Responses{{Status: 200, Body: "OK"}}},
		},
		"single array response": {
			data: "response:\n- status: 200\n  body: OK\n",
			exp:  Imposter{Response: Responses{{Status: 200, Body: "OK"}}},
		},
		"multiple array response": {
			data: "response:\n- status: 200\n  body: OK\n- status: 404\n  body: Not Found\n",
			exp:  Imposter{Response: Responses{{Status: 200, Body: "OK"}, {Status: 404, Body: "Not Found"}}},
		},
		"empty array": {
			data: "response: []\n",
			exp:  Imposter{Response: Responses{}},
		},
		"null array": {
			data: "response: \n",
			exp:  Imposter{Response: nil},
		}}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			var got Imposter
			err := yaml.Unmarshal([]byte(tc.data), &got)
			require.NoError(t, err)
			assert.Equal(t, tc.exp, got)
		})
	}
}
