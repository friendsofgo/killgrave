package killgrave

import (
	"path"
)

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
	SchemaFile *string            `json:"schema_file"`
	Headers    *map[string]string `json:"headers"`
}

// Response represent the structure of real response
type Response struct {
	Status   int                `json:"status"`
	Body     string             `json:"body"`
	BodyFile *string            `json:"bodyFile"`
	Headers  *map[string]string `json:"headers"`
}
