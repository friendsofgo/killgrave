package generator

import (
	"encoding/json"
	"errors"
	"github.com/friendsofgo/killgrave/internal/server/http"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestGenerate(t *testing.T) {
	var testCases = []struct {
		name       string
		inputPath  string
		exampleOutputPath string
		generateV3 bool
		err        error
	}{
		{
			"malformatted swagger (3 -> 2) yml",
			"test/testdata/oas3_sample.yml",
			"test/testdata/test_request.imp.json",
			false,
			errors.New("cannot generate v2 swagger from v3 file"),
		},
		{
			"malformatted swagger (2 -> 3) yml",
			"test/testdata/oas2_sample.yml",
			"test/testdata/test_request.imp.json",
			true,
			errors.New("cannot generate v3 swagger from v2 file"),
		},
		{
			"valid v2 yml",
			"test/testdata/oas2_sample.yml",
			"test/testdata/test_request.imp.json",
			false,
			nil,
		},
		{
			"valid v3 yml",
			"test/testdata/oas3_sample.yml",
			"test/testdata/test_request.imp.json",
			true,
			nil,
		},
		{
			"malformatted swagger (3 -> 2) json",
			"test/testdata/oas3_sample.json",
			"test/testdata/test_request.imp.json",
			false,
			errors.New("cannot generate v2 swagger from v3 file"),
		},
		{
			"malformatted swagger (2 -> 3) json",
			"test/testdata/oas2_sample.json",
			"test/testdata/test_request.imp.json",
			true,
			errors.New("cannot generate v3 swagger from v2 file"),
		},
		{
			"valid v2 json",
			"test/testdata/oas2_sample.json",
			"test/testdata/test_request.imp.json",
			false,
			nil,
		},
		{
			"valid v3 json",
			"test/testdata/oas3_sample.json",
			"test/testdata/test_request.imp.json",
			true,
			nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGenerator(tt.inputPath)

			swagger, err := ioutil.ReadFile(tt.inputPath)
			if err != nil {
				t.Fatalf("%v: error trying to read swagger file: '%s'", err, tt.inputPath)
			}

			var expectedImpostersJson []byte
			expectedImpostersJson, err = ioutil.ReadFile(tt.exampleOutputPath)
			if err != nil {
				t.Fatalf("%v: error trying to read example imposter file: '%s'", err, tt.exampleOutputPath)
			}

			var expectedImposters []http.Imposter
			err = json.Unmarshal(expectedImpostersJson, &expectedImposters)
			if err != nil {
				t.Fatalf("%v: error trying to unmarshal example imposter file: '%s'", err, tt.exampleOutputPath)
			}

			var actualImposters *[]http.Imposter
			actualImposters, err = g.GenerateSwagger(swagger, tt.generateV3)

			if err == nil {
				if tt.err != nil {
					t.Fatalf("expected an error and got nil")
				}

				assert.Equal(t, len(expectedImposters), len(*actualImposters), "expected same number of imposters")
				for i, _ := range expectedImposters {
					assert.Equal(t, expectedImposters[i], (*actualImposters)[i], "expected imposter %d and actual imposter %d should be equal", i, i)
				}
			}

			if err != nil {
				if tt.err == nil {
					t.Fatalf("did not expect any errors and got %+v", err)
				}
			}
		})
	}
}
