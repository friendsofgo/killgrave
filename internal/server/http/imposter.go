package http

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const imposterExtension = ".imp.json"

// Imposter define an imposter structure
type Imposter struct {
	BasePath string
	Request  Request  `json:"request"`
	Response Response `json:"response"`
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
	BodyFile *string            `json:"bodyFile"`
	Headers  *map[string]string `json:"headers"`
}

func findImposters(impostersDirectory string, done <-chan struct{}) (<-chan string, <-chan error) {
	imposterFiles := make(chan string)
	errc := make(chan error, 1)
	go func() {
		defer close(imposterFiles)
		errc <- filepath.Walk(impostersDirectory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("%w: error accsessing path %q", err, path)
			}
			if info.IsDir() || strings.HasSuffix(info.Name(), imposterExtension) {
				return nil
			}
			select {
			case imposterFiles <- path:
				return nil
			case <-done:
				return fmt.Errorf("find imposters canceled")
			}
		})
	}()
	return imposterFiles, errc
}
