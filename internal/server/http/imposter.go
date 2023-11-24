package http

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// ImposterType allows to know the imposter type we're dealing with
type ImposterType int

const (
	jsonImposterExtension = ".imp.json"
	ymlImposterExtension  = ".imp.yml"
	yamlImposterExtension = ".imp.yaml"
)

const (
	// JSONImposter allows to know when we're dealing with a JSON imposter
	JSONImposter ImposterType = iota
	// YAMLImposter allows to know when we're dealing with a YAML imposter
	YAMLImposter
)

// ImposterConfig is used to load imposters based on which type they are
type ImposterConfig struct {
	Type     ImposterType
	FilePath string
}

// Imposter define an imposter structure
type Imposter struct {
	BasePath string   `json:"-" yaml:"-"`
	Path     string   `json:"-" yaml:"-"`
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

// Delay returns delay for response that user can specify in imposter config
func (i *Imposter) Delay() time.Duration {
	return i.Response.Delay.Delay()
}

// CalculateFilePath calculate file path based on basePath of imposter directory
func (i *Imposter) CalculateFilePath(filePath string) string {
	return path.Join(i.BasePath, filePath)
}

// Request represent the structure of real request
type Request struct {
	Method     string             `json:"method"`
	Endpoint   string             `json:"endpoint"`
	SchemaFile *string            `json:"schemaFile"`
	Params     *map[string]string `json:"params"`
	Headers    *map[string]string `json:"headers"`
}

// Response represent the structure of real response
type Response struct {
	Status   int                `json:"status"`
	Body     string             `json:"body"`
	BodyFile *string            `json:"bodyFile" yaml:"bodyFile"`
	Headers  *map[string]string `json:"headers"`
	Delay    ResponseDelay      `json:"delay" yaml:"delay"`
}

type ImposterFs struct {
	path string
	fs   fs.FS
}

func NewImposterFS(path string) (ImposterFs, error) {
	// TODO: What if user lacks permissions?
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ImposterFs{}, fmt.Errorf("%w: the directory %s doesn't exists", err, path)
	}

	return ImposterFs{
		path: path,
		fs:   os.DirFS(path),
	}, nil
}

func (i ImposterFs) FindImposters(impostersCh chan []Imposter) error {
	err := fs.WalkDir(i.fs, ".", func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("%w: error finding imposters", err)
		}

		var cfg ImposterConfig
		filename := info.Name()
		if !info.IsDir() {
			switch {
			case strings.HasSuffix(filename, jsonImposterExtension):
				cfg = ImposterConfig{JSONImposter, path}
			case strings.HasSuffix(filename, yamlImposterExtension), strings.HasSuffix(filename, ymlImposterExtension):
				cfg = ImposterConfig{YAMLImposter, path}
			default:
				return nil
			}
			imposters, err := i.unmarshalImposters(cfg)
			if err != nil {
				return err
			}
			impostersCh <- imposters
		}
		return nil
	})
	return err
}

func (i ImposterFs) unmarshalImposters(imposterConfig ImposterConfig) ([]Imposter, error) {
	// TODO: Error handling?
	bytes, _ := fs.ReadFile(i.fs, imposterConfig.FilePath)

	var parseError error
	var imposters []Imposter

	switch imposterConfig.Type {
	case JSONImposter:
		parseError = json.Unmarshal(bytes, &imposters)
	case YAMLImposter:
		parseError = yaml.Unmarshal(bytes, &imposters)
	default:
		parseError = fmt.Errorf("unsupported imposter type %v", imposterConfig.Type)
	}

	if parseError != nil {
		return nil, fmt.Errorf("%w: error while unmarshalling imposter's file %s", parseError, imposterConfig.FilePath)
	}

	for i := range imposters {
		imposters[i].BasePath = filepath.Dir(imposterConfig.FilePath)
		imposters[i].Path = imposterConfig.FilePath
	}

	return imposters, nil
}
