package killgrave

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigFromFile(t *testing.T) {
	tests := map[string]struct {
		input     string
		expected  Config
		wantError bool
	}{
		"valid config file": {"test/testdata/config.yml", validConfig(), false},
		"file not found":    {"test/testdata/file.yml", Config{}, true},
		"wrong yaml file":   {"test/testdata/wrong_config.yml", Config{}, true},
		"empty config file": {"", Config{}, true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewConfigFromFile(tc.input)
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
		"valid mode record":  {"record", ProxyRecord, false},
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
			"ProxyMissing must be return missing string",
			ProxyMissing,
			"missing",
		},
		{
			"ProxyRecord must be return record string",
			ProxyRecord,
			"record",
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

func TestNewConfig(t *testing.T) {
	type args struct {
		impostersPath string
		host          string
		port          int
	}
	tests := []struct {
		name string
		args args
		want Config
		err  error
	}{
		{
			name: "empty imposters path",
			args: args{
				impostersPath: "",
				host:          "localhost",
				port:          80,
			},
			want: Config{},
			err:  errEmptyImpostersPath,
		},
		{
			name: "empty host path",
			args: args{
				impostersPath: "imposters",
				host:          "",
				port:          80,
			},
			want: Config{},
			err:  errEmptyHost,
		},
		{
			name: "invalid port",
			args: args{
				impostersPath: "imposters",
				host:          "localhost",
				port:          -1000,
			},
			want: Config{},
			err:  errInvalidPort,
		},
		{
			name: "valid config",
			args: args{
				impostersPath: "imposters",
				host:          "localhost",
				port:          80,
			},
			want: Config{
				ImpostersPath: "imposters",
				Port:          80,
				Host:          "localhost",
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.args.impostersPath, tt.args.host, tt.args.port, false)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConfig_ConfigureProxy(t *testing.T) {
	tests := []struct{
		name string
		url string
		mode ProxyMode
		err error
	}{
		{"valid proxy configuration", "https://friendsofgo.tech", ProxyAll, nil},
		{"none proxy configuration", "", ProxyNone, nil},
		{"invalid proxy configuration", "", ProxyAll, errMandatoryURL},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			cfg, err := NewConfig("imposters", "localhost", 80, false)
			assert.NoError(t, err)

			err = cfg.ConfigureProxy(tt.mode, tt.url, "")
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestConfig_ConfigureProxyRecord(t *testing.T) {
	tests := []struct {
		name           string
		recordFilePath string
		err            error
	}{
		{"valid configuration for Proxy Record", "some_file.imp.json", nil},
		{"invalid configuration for Proxy Record", "", errMandatoryRecordFilePath},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := NewConfig("imposters", "localhost", 80, false)
			assert.NoError(t, err)

			err = cfg.ConfigureProxy(ProxyRecord, "https://friendsofgo.tech", tt.recordFilePath)
			assert.Equal(t, tt.err, err)
		})
	}
}
