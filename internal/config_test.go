package killgrave

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestReadConfigFile(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected Config
		err      error
	}{
		"valid config file": {"test/testdata/config.yml", validConfig(), nil},
		"file not found":    {"test/testdata/file.yml", Config{}, errors.New("error")},
		"wrong yaml file":   {"test/testdata/wrong_config.yml", Config{}, errors.New("error")},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var got Config
			err := ReadConfigFile(tc.input, &got)

			if err != nil && tc.err == nil {
				t.Fatalf("not expected any erros and got %v", err)
			}

			if err == nil && tc.err != nil {
				t.Fatalf("expected an error and got nil")
			}

			if !reflect.DeepEqual(tc.expected, got) {
				t.Fatalf("expected: %v, got: %v", tc.expected, got)
			}

		})
	}
}

func validConfig() Config {
	return Config{
		ImpostersPath: "imposters",
		Port:          3000,
		Host:          "localhost",
		CORS: ConfigCORS{
			Methods:          []string{"GET"},
			Origins:          []string{"*"},
			Headers:          []string{"Content-Type"},
			ExposedHeaders:   []string{"Cache-Control"},
			AllowCredentials: true,
		},
	}
}
