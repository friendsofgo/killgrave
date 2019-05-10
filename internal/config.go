package killgrave

// Config representation of config file yaml
type Config struct {
	ImpostersPath string     `yaml:"imposters_path"`
	Port          int        `yaml:"port"`
	Host          string     `yaml:"host"`
	Cors          ConfigCors `yaml:"cors"`
}

// ConfigCors representation of section CORS of the yaml
type ConfigCors struct {
	Methods       []string `yaml:"methods"`
	Headers       []string `yaml:"headers"`
	Origins       []string `yaml:"origins"`
	ExposeHeaders []string `yaml:"expose_headers"`
}
