package http

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

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
	BasePath string    `json:"-" yaml:"-"`
	Path     string    `json:"-" yaml:"-"`
	Request  Request   `json:"request"`
	Response Responses `json:"response"`
	resIdx   int
}

// NextResponse returns the imposter's response.
// If there are multiple responses, it will return them sequentially.
func (i *Imposter) NextResponse() Response {
	r := i.Response[i.resIdx]
	i.resIdx = (i.resIdx + 1) % len(i.Response)
	return r
}

// CalculateFilePath calculate file path based on basePath of imposter's directory
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

// Responses is a wrapper for Response, to allow the use of either a single
// response or an array of responses, while keeping backwards compatibility.
type Responses []Response

func (rr *Responses) MarshalJSON() ([]byte, error) {
	if len(*rr) == 1 {
		return json.Marshal((*rr)[0])
	}
	return json.Marshal(*rr)
}

func (rr *Responses) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*rr = nil
		return nil
	}

	if data[0] == '[' {
		return json.Unmarshal(data, (*[]Response)(rr))
	}

	var r Response
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}

	*rr = Responses{r}
	return nil
}

func (rr *Responses) MarshalYAML() (interface{}, error) {
	if len(*rr) == 1 {
		return (*rr)[0], nil
	}
	return *rr, nil
}

func (rr *Responses) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var r Response
	if err := unmarshal(&r); err == nil {
		*rr = Responses{r}
		return nil
	}

	var tmp []Response
	if err := unmarshal(&tmp); err != nil {
		return err
	}
	*rr = tmp
	return nil
}

type ImposterFs struct {
	path string
	fs   fs.FS
}

func NewImposterFS(path string) (ImposterFs, error) {
	_, err := os.Stat(path)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			return ImposterFs{}, fmt.Errorf("the directory '%s' does not exist", path)
		case os.IsPermission(err):
			return ImposterFs{}, fmt.Errorf("could not read the directory '%s': permission denied", path)
		default:
			return ImposterFs{}, fmt.Errorf("could not read the directory '%s': %w", path, err)
		}
	}

	return ImposterFs{
		path: path,
		fs:   os.DirFS(path),
	}, nil
}

func (ifs ImposterFs) FindImposters(impostersCh chan []Imposter) error {
	err := fs.WalkDir(ifs.fs, ".", func(path string, info fs.DirEntry, err error) error {
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
			imposters, err := ifs.unmarshalImposters(cfg)
			if err != nil {
				return err
			}
			impostersCh <- imposters
		}
		return nil
	})
	close(impostersCh)
	return err
}

func (ifs ImposterFs) unmarshalImposters(imposterConfig ImposterConfig) ([]Imposter, error) {
	imposterFile, _ := ifs.fs.Open(imposterConfig.FilePath)
	defer imposterFile.Close()

	bytes, _ := io.ReadAll(imposterFile)

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
		imposters[i].BasePath = filepath.Dir(filepath.Join(ifs.path, imposterConfig.FilePath))
		imposters[i].Path = imposterConfig.FilePath
	}

	return imposters, nil
}
