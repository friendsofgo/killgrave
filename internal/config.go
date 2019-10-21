package killgrave

import (
	"errors"
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
	Url  string    `yaml:"url"`
	Mode ProxyMode `yaml:"mode"`
}

// ProxyMode is enumeration of proxy server modes
type ProxyMode int

const (
	// ProxyNone server is off
	ProxyNone ProxyMode = iota
	// ProxyMissing handle only missing requests are proxied
	ProxyMissing
	// ProxyAll all requests are proxied
	ProxyAll
)

// UnmarshalYAML implementation of yaml.Unmarshaler interface
func (mode *ProxyMode) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var textMode string
	if err := unmarshal(&textMode); err != nil {
		return err
	}
	switch textMode {
	case "all":
		*mode = ProxyAll
	case "missing":
		*mode = ProxyMissing
	case "none":
		*mode = ProxyNone
	default:
		return errors.New("unknown proxy mode: " + textMode)
	}
	return nil
}

// ConfigOpt function to encapsulate optional parameters
type ConfigOpt func(cfg *Config) error

// NewConfig initialize the config
func NewConfig(impostersPath, host string, port int, opts ...ConfigOpt) (Config, error) {
	cfg := Config{
		ImpostersPath: impostersPath,
		Host:          host,
		Port:          port,
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
