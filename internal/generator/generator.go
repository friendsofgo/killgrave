package generator

import (
	"fmt"
	v2 "github.com/friendsofgo/killgrave/internal/generator/v2"
	v3 "github.com/friendsofgo/killgrave/internal/generator/v3"
	"github.com/friendsofgo/killgrave/internal/server/http"
	"strings"
)

type Generator struct {
	fileName string
}

type GenerationMode int

const (
	YAML GenerationMode = iota
	JSON
)

func NewGenerator(fileName string) *Generator {
	return &Generator{
		fileName: strings.ToLower(fileName),
	}
}

func (g *Generator) generationMode() (*GenerationMode, error) {
	var result GenerationMode

	if strings.HasSuffix(g.fileName, ".yml") ||
		strings.HasSuffix(g.fileName, ".yaml") {
		result = YAML
		return &result, nil
	}

	if strings.HasSuffix(g.fileName, ".json") {
		result = JSON
		return &result, nil
	}

	return nil, fmt.Errorf("unknown file extension for swagger file name: %v", g.fileName)
}

func (g *Generator) GenerateSwagger(swagger []byte, useV3 bool) (*[]http.Imposter, error) {
	mode, err := g.generationMode()
	if err != nil {
		return nil, err
	}

	if useV3 {
		return v3.GenerateSwagger(swagger)
	} else {
		if *mode == JSON {
			return v2.GenerateSwaggerJSON(swagger)
		} else {
			return v2.GenerateSwaggerYAML(swagger)
		}
	}
}
