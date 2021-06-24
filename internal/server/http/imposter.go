package http

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
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
	BasePath  string
	Request   Request    `json:"request"`
	Responses []Response `json:"responses"`

	// Field for handling burst response
	responseHandler ResponseHandler
}

// CalculateFilePath calculate file path based on basePath of imposter directory
func (i *Imposter) CalculateFilePath(filePath string) string {
	return path.Join(i.BasePath, filePath)
}

// GetResponse is used to get dynamic/random response.
func (i *Imposter) GetResponse() Response {
	// If no response provided then returning 404 response
	if len(i.Responses) == 0 {
		return Response{Status: 404}
	}
	if i.responseHandler.scheduleMap == nil {
		i.fillDefaults()
	}
	index := i.responseHandler.GetIndex()
	return i.Responses[index]
}

// fillDefaults is used to populate default values for imposters fields.
func (i *Imposter) fillDefaults() {
	i.responseHandler.fillDefaults(i)
}

// Request represent the structure of real request
type Request struct {
	Method       string             `json:"method"`
	Endpoint     string             `json:"endpoint"`
	SchemaFile   *string            `json:"schemaFile"`
	Params       *map[string]string `json:"params"`
	Headers      *map[string]string `json:"headers"`
	ResponseMode string             `json:"responseMode"`
}

// Response represent the structure of real response
type Response struct {
	Status   int                `json:"status"`
	Body     string             `json:"body"`
	BodyFile *string            `json:"bodyFile" yaml:"bodyFile"`
	Headers  *map[string]string `json:"headers"`
	RDelay   ResponseDelay      `json:"delay" yaml:"delay"`
	Burst    int                `json:"burst" yaml:"burst"`
}

// Delay returns delay for response that user can specify in imposter config
func (r *Response) Delay() time.Duration {
	return r.RDelay.Delay()
}

func findImposters(impostersDirectory string, imposterConfigCh chan ImposterConfig) error {
	err := filepath.Walk(impostersDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("%w: error finding imposters", err)
		}

		filename := info.Name()
		if !info.IsDir() {
			if strings.HasSuffix(filename, jsonImposterExtension) {
				imposterConfigCh <- ImposterConfig{JSONImposter, path}
			} else if strings.HasSuffix(filename, yamlImposterExtension) || strings.HasSuffix(filename, ymlImposterExtension) {
				imposterConfigCh <- ImposterConfig{YAMLImposter, path}
			}
		}
		return nil
	})
	return err
}
