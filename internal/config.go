package killgrave

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Config representation of config file yaml
type Config struct {
	ImpostersPath string     `yaml:"imposters_path"`
	Port          int        `yaml:"port"`
	Host          string     `yaml:"host"`
	CORS          ConfigCORS `yaml:"cors"`
}

// ConfigCORS representation of section CORS of the yaml
type ConfigCORS struct {
	Methods        []string `yaml:"methods"`
	Headers        []string `yaml:"headers"`
	Origins        []string `yaml:"origins"`
	ExposedHeaders []string `yaml:"exposed_headers"`
}

// ReadConfigFile unmarshal content of config file to Config struct
func ReadConfigFile(path string, config *Config) error {
	configFile, err := os.Open(path)
	if err != nil {
		return errors.Wrapf(err, "error trying to read config file: %s", path)
	}
	defer configFile.Close()

	bytes, _ := ioutil.ReadAll(configFile)
	if err := yaml.Unmarshal(bytes, config); err != nil {
		return errors.Wrapf(err, "error while unmarshall configFile file %s", path)
	}
	return nil
}
