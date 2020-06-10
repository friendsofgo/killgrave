package http

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const imposterExtension = ".imp.json"

// Impostor define an imposter structure
type Impostor struct {
	BasePath string
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

// Delay returns delay for response that user can specify in imposter config
func (i *Impostor) Delay() time.Duration {
	return i.Response.Delay.Delay()
}

// CalculateFilePath calculate file path based on basePath of imposter directory
func (i *Impostor) CalculateFilePath(filePath string) string {
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
	Delay    ResponseDelay      `json:"delay"`
}

func findImposters(impostorsDirectory string, imposterFileCh chan string) error {
	err := filepath.Walk(impostorsDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("%w: error finding impostors", err)
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), imposterExtension) {
			imposterFileCh <- path
		}
		return nil
	})
	return err
}
