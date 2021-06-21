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

	// Fields for handling burst response
	count        int
	currentIndex int
	scheduleMap  map[int]int
}

// CalculateFilePath calculate file path based on basePath of imposter directory
func (i *Imposter) CalculateFilePath(filePath string) string {
	return path.Join(i.BasePath, filePath)
}

// GetResponse method is used to get response. (Implemented/Changed in burst mode)
func (i *Imposter) GetResponse() Response {
	// Checking if imposter has at least one response available
	if len(i.Responses) == 0 {
		return Response{}
	}

	// Filling default values in the struct if not done already.
	if i.scheduleMap == nil {
		i.fillDefaults()
	}

	var ind = i.currentIndex
	i.updateCounter() // Updating counters for burst mode.
	return i.Responses[ind]
}

// updateCounter method updates neccesary counters/indexes to be used in burst mode.
func (i *Imposter) updateCounter() {

	// Calculating next index
	i.count += 1
	if i.scheduleMap[i.currentIndex] < i.count {
		i.currentIndex += 1
	}

	// Wrapping logic for counter and index
	if i.currentIndex > len(i.Responses)-1 {
		i.currentIndex = 0
		i.count = 1
	}
}

// fillDefaults method is used to populate default values for imposters fields.
func (i *Imposter) fillDefaults() {
	var scheduleMap = make(map[int]int, len(i.Responses))
	for ind, resp := range i.Responses {
		resp.fillDefaults()

		if ind != 0 {
			scheduleMap[ind] = scheduleMap[ind-1] + resp.Burst
		} else {
			scheduleMap[ind] = resp.Burst
		}
	}
	i.scheduleMap = scheduleMap
	i.count = 1
	i.currentIndex = 0
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
	RDelay   ResponseDelay      `json:"delay" yaml:"delay"`
	Burst    int                `json:"burst" yaml:"burst"`
}

// Delay returns delay for response that user can specify in imposter config
func (r *Response) Delay() time.Duration {
	return r.RDelay.Delay()
}

// fillDefaults method is used to populate default values for response fields.
func (r *Response) fillDefaults() {
	if r.Burst <= 0 {
		r.Burst = 1
	}
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
