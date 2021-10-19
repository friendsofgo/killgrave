package generate

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/friendsofgo/killgrave/internal/generator"
	"github.com/friendsofgo/killgrave/internal/server/http"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	_defaultImpostersPath = "impostersFile"
	_defaultSwaggerFile   = "swagger.yaml"
	_defaultOutputYaml    = false
)

var (
	errGetDataFromImpostersFileFlag = errors.New("error trying to get data from imposters-file flag")
 	errGetDataFromSwaggerFileFlag   = errors.New("error trying to get data from swagger-file flag")
	errGetDataFromOutputYamlFlag    = errors.New("error trying to get data from output-yaml flag")
)

type config struct {
	impostersFile string
	swaggerFile   string
	outputYAML    bool
	useV3         bool
}

// NewGenerate2Cmd returns cobra.Command to run generate sub command, this command will be used to generate impostersFile from swagger files
func NewGenerate2Cmd() *cobra.Command {
	return newGenerateCmd(false)
}

// NewGenerate3Cmd returns cobra.Command to run generate3 sub command, this command will be used to generate impostersFile from swagger files
func NewGenerate3Cmd() *cobra.Command {
	return newGenerateCmd(true)
}

// NewGenerateCmd returns cobra.Command to run generate sub command, this command will be used to generate impostersFile from swagger files
func newGenerateCmd(useV3 bool) *cobra.Command {
	var cfg *config
	var err error

	preRunE := func(cmd *cobra.Command, args []string) error {
		cfg, err = prepareConfig(cmd, useV3)
		if err != nil {
			return err
		}
		return nil
	}

	runE := func(cmd *cobra.Command, args []string) error {
		return runGenerate(cfg)
	}

	var cmd *cobra.Command

	if useV3 {
		cmd = &cobra.Command{
			Use:   "swagger3-gen",
			Short: "Generate impostersFile based on a swagger 3 (OpenAPI 3 Specification) file",
			PreRunE: preRunE,
			RunE: runE,
			Args: cobra.NoArgs,
		}
	} else {
		cmd = &cobra.Command{
			Use:   "swagger-gen",
			Short: "Generate impostersFile based on a swagger (OpenAPI 2 Specification) file",
			PreRunE: preRunE,
			RunE: runE,
			Args: cobra.NoArgs,
		}
	}

	cmd.PersistentFlags().String("impostersFile", _defaultImpostersPath, "directory where your impostersFile are saved. Existing files will be silently overwritten if they have the same name")
	cmd.PersistentFlags().String("swagger-file", _defaultSwaggerFile, "swagger file for use by swagger-gen")
	cmd.PersistentFlags().Bool("output-yaml", _defaultOutputYaml, "true if YAML impostersFile should be generated; false if JSON impostersFile should be generated")

	return cmd
}

func runGenerate(cfg *config) error {
	g := generator.NewGenerator(cfg.swaggerFile)

	swagger, err := ioutil.ReadFile(cfg.swaggerFile)
	if err != nil {
		return fmt.Errorf("%w: error trying to read swagger file: '%s'", err, cfg.swaggerFile)
	}

	var imposters *[]http.Imposter

	imposters, err = g.GenerateSwagger(swagger, cfg.useV3)
	if err != nil {
		return fmt.Errorf("unable to generate imposters: %w", err)
	}

	var output []byte

	if cfg.outputYAML {
		output, err = yaml.Marshal(imposters)
	} else {

		output, err = json.Marshal(imposters)
	}

	if err != nil {
		return fmt.Errorf("unable to marshal imposters: %w", err)
	}

	err = ioutil.WriteFile(cfg.impostersFile, output, 0644)
	if err != nil {
		return fmt.Errorf("%w: error trying to write to imposters file: '%s'", err, cfg.impostersFile)
	}

	return err
}

func prepareConfig(cmd *cobra.Command, useV3 bool) (cfg *config, err error) {
	impostersFilePath, err := cmd.Flags().GetString("imposters-file")
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, errGetDataFromImpostersFileFlag)
	}

	swaggerFilePath, err := cmd.Flags().GetString("swagger-file")
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, errGetDataFromSwaggerFileFlag)
	}

	outputYAML, err := cmd.Flags().GetBool("output-yaml")
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, errGetDataFromOutputYamlFlag)
	}

	return &config{
		impostersFile: impostersFilePath,
		swaggerFile:   swaggerFilePath,
		outputYAML:    outputYAML,
		useV3:         useV3,
	}, nil
}
