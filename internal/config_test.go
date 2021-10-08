package killgrave

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	tests := map[string]struct {
		input     string
		expected  Config
		wantError bool
	}{
		"valid config file": {"test/testdata/config.yml", validConfig(), false},
		"file not found":    {"test/testdata/file.yml", Config{}, true},
		"wrong yaml file":   {"test/testdata/wrong_config.yml", Config{}, true},
		"empty config file": {"", Config{}, false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewConfig(
				"",
				"",
				0,
				false,
				WithProxyConfiguration(ProxyNone.String(), ""),
				WithConfigFile(tc.input),
				WithWatcherConfiguration(tc.expected.Watcher))

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestProxyModeParseString(t *testing.T) {
	testCases := map[string]struct {
		input     string
		expected  ProxyMode
		wantError bool
	}{
		"valid mode":   {"all", ProxyAll, false},
		"unknown mode": {"UnKnOwn1", ProxyNone, true},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mode, err := StringToProxyMode(tc.input)

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expected, mode)

		})
	}
}

func TestProxyModeUnmarshal(t *testing.T) {
	testCases := map[string]struct {
		input     interface{}
		expected  ProxyMode
		wantError bool
	}{
		"valid mode all":     {"all", ProxyAll, false},
		"valid mode missing": {"missing", ProxyMissing, false},
		"valid mode none":    {"none", ProxyNone, false},
		"empty mode":         {"", ProxyNone, true},
		"invalid mode":       {"nonsens23e", ProxyNone, true},
		"error input":        {123, ProxyNone, true},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var mode ProxyMode
			err := mode.UnmarshalYAML(func(i interface{}) error {
				s := i.(*string)
				input, ok := tc.input.(string)
				if !ok {
					return errors.New("error")
				}
				*s = input
				return nil
			})

			if tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expected, mode)
		})
	}
}

func validConfig() Config {
	return Config{
		ImpostersPath: "test/testdata/imposters",
		Port:          3000,
		Host:          "localhost",
		CORS: ConfigCORS{
			Methods:          []string{"GET"},
			Origins:          []string{"*"},
			Headers:          []string{"Content-Type"},
			ExposedHeaders:   []string{"Cache-Control"},
			AllowCredentials: true,
		},
		Watcher: true,
		Secure:  true,
	}
}

func TestProxyMode_String(t *testing.T) {
	tests := []struct {
		name string
		p    ProxyMode
		want string
	}{
		{
			"ProxyNone must be return none string",
			ProxyNone,
			"none",
		},
		{
			"ProxyNone must be return missing string",
			ProxyMissing,
			"missing",
		},
		{
			"ProxyNone must be return all string",
			ProxyAll,
			"all",
		},
		{
			"An invalid mode must return none string",
			ProxyMode(33),
			"none",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.p.String())
		})
	}
}
