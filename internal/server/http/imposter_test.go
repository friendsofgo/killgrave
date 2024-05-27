package http

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"

	"github.com/stretchr/testify/require"
)

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
