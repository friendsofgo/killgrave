package killgrave

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

// Config representation of config file yaml
type Config struct {
	ImpostersPath string      `yaml:"imposters_path"`
	Port          int         `yaml:"port"`
	Host          string      `yaml:"host"`
	CORS          ConfigCORS  `yaml:"cors"`
	Proxy         ConfigProxy `yaml:"proxy"`
	Secure        bool        `yaml:"secure"`
	Watcher       bool        `yaml:"watcher"`
}

// ConfigCORS representation of section CORS of the yaml
type ConfigCORS struct {
	Methods          []string `yaml:"methods"`
	Headers          []string `yaml:"headers"`
	Origins          []string `yaml:"origins"`
	ExposedHeaders   []string `yaml:"exposed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
}

// ConfigProxy is a representation of section proxy of the yaml
type ConfigProxy struct {
	URL  string    `yaml:"url"`
	Mode ProxyMode `yaml:"mode"`
}

// ProxyMode is enumeration of proxy server modes
type ProxyMode uint8

const (
	// ProxyNone server is off
	ProxyNone ProxyMode = iota
	// ProxyMissing handle only missing requests are proxied
	ProxyMissing
	// ProxyAll all requests are proxied
	ProxyAll
)

func (p ProxyMode) String() string {
	m := map[ProxyMode]string{
		ProxyNone:    "none",
		ProxyMissing: "missing",
		ProxyAll:     "all",
	}

	s, ok := m[p]
	if !ok {
		return "none"
	}
	return s
}

// StringToProxyMode convert string into a ProxyMode if not exists return a none mode and an error
func StringToProxyMode(t string) (ProxyMode, error) {
	m := map[string]ProxyMode{
		"none":    ProxyNone,
		"missing": ProxyMissing,
		"all":     ProxyAll,
	}

	p, ok := m[t]
	if !ok {
		return ProxyNone, fmt.Errorf("unknown proxy mode: %s", t)
	}

	return p, nil
}

// UnmarshalYAML implementation of yaml.Unmarshaler interface
func (p *ProxyMode) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var proxyMode string
	if err := unmarshal(&proxyMode); err != nil {
		return err
	}

	m, err := StringToProxyMode(proxyMode)
	if err != nil {
		return err
	}

	*p = m
	return nil
}

// ConfigOpt function to encapsulate optional parameters
type ConfigOpt func(cfg *Config) error

// NewConfig initialize the config
func NewConfig(impostersPath, host string, port int, secure bool, opts ...ConfigOpt) (Config, error) {
	cfg := Config{
		ImpostersPath: impostersPath,
		Host:          host,
		Port:          port,
		Secure:        secure,
	}

	for _, opt := range opts {
		if err := opt(&cfg); err != nil {
			return Config{}, err
		}
	}

	return cfg, nil
}

// WithConfigFile unmarshal content of config file to Config struct
func WithConfigFile(cfgPath string) ConfigOpt {
	return func(cfg *Config) error {
		if cfgPath == "" {
			return nil
		}

		configFile, err := os.Open(cfgPath)
		if err != nil {
			return fmt.Errorf("%w: error trying to read config file: %s, using default configuration instead", err, cfgPath)
		}
		defer configFile.Close()

		bytes, _ := ioutil.ReadAll(configFile)
		if err := yaml.Unmarshal(bytes, cfg); err != nil {
			return fmt.Errorf("%w: error while unmarshall configFile file %s, using default configuration instead", err, cfgPath)
		}

		cfg.ImpostersPath = path.Join(path.Dir(cfgPath), cfg.ImpostersPath)

		return nil
	}
}

// WithProxyConfiguration preparing the server with the proxy configuration that the user has indicated
func WithProxyConfiguration(proxyMode, proxyURL string) ConfigOpt {
	return func(cfg *Config) error {
		mode, _ := StringToProxyMode(proxyMode)
		cfg.Proxy.Mode = mode
		cfg.Proxy.URL = proxyURL

		return nil
	}
}

// WithWatcherConfiguration preparing server to do auto-reload
func WithWatcherConfiguration(watcher bool) ConfigOpt {
	return func(cfg *Config) error {

		cfg.Watcher = watcher
		return nil
	}
}
