package killgrave

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewConfigFromFile(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected Config
		err      error
	}{
		"valid config file": {"test/testdata/config.yml", validConfig(), nil},
		"file not found":    {"test/testdata/file.yml", Config{}, errors.New("error")},
		"wrong yaml file":   {"test/testdata/wrong_config.yml", Config{}, errors.New("error")},
		"empty config file": {"", Config{}, errors.New("error")},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewConfigFromFile(tc.input)

			if err != nil && tc.err == nil {
				t.Fatalf("not expected any errors and got %v", err)
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

func TestProxyModeParseString(t *testing.T) {
	testCases := map[string]struct {
		input    string
		expected ProxyMode
		err      error
	}{
		"valid mode":   {"all", ProxyAll, nil},
		"unknown mode": {"UnKnOwn1", ProxyNone, errors.New("error")},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mode, err := StringToProxyMode(tc.input)

			if err != nil && tc.err == nil {
				t.Fatalf("not expected any erros and got %v", err)
			}
			if err == nil && tc.err != nil {
				t.Fatalf("expected an error and got nil")
			}
			if tc.expected != mode {
				t.Fatalf("expected: %v, got: %v", tc.expected, mode)
			}
		})
	}
}

func TestProxyModeUnmarshal(t *testing.T) {
	testCases := map[string]struct {
		input    interface{}
		expected ProxyMode
		err      error
	}{
		"valid mode all":     {"all", ProxyAll, nil},
		"valid mode missing": {"missing", ProxyMissing, nil},
		"valid mode none":    {"none", ProxyNone, nil},
		"empty mode":         {"", ProxyNone, errors.New("error")},
		"invalid mode":       {"nonsens23e", ProxyNone, errors.New("error")},
		"error input":        {123, ProxyNone, errors.New("error")},
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
			if err != nil && tc.err == nil {
				t.Fatalf("not expected any erros and got %v", err)
			}

			if err == nil && tc.err != nil {
				t.Fatalf("expected an error and got nil")
			}
			if tc.expected != mode {
				t.Fatalf("expected: %v, got: %v", tc.expected, mode)
			}
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
			if got := tt.p.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
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
			got, err := NewConfig(tt.args.impostersPath, tt.args.host, tt.args.port, false, "", "")
			if tt.err != nil && !errors.Is(err, tt.err) {
				t.Errorf("NewConfig() error got = %v, err want %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_ConfigureProxy(t *testing.T) {
	expected := Config{
		ImpostersPath: "imposters",
		Port:          80,
		Host:          "localhost",
		Proxy: ConfigProxy{
			Url:  "https://friendsofgo.tech",
			Mode: ProxyAll,
		},
	}

	got, err := NewConfig("imposters", "localhost", 80, false, "", "")
	if err != nil {
		t.Fatalf("error not expected: %v", err)
	}
	got.ConfigureProxy(ProxyAll, "https://friendsofgo.tech")

	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("got = %v, want %v", expected, got)
	}
}
