// Package openapi provides OpenAPI 3.x specification parsing and schema extraction.
package openapi

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/JSLEEKR/schemadiff/pkg/schema"
	"gopkg.in/yaml.v3"
)

// Spec represents a parsed OpenAPI 3.x specification.
type Spec struct {
	OpenAPI    string                 `json:"openapi" yaml:"openapi"`
	Info       Info                   `json:"info" yaml:"info"`
	Paths      map[string]*PathItem   `json:"paths" yaml:"paths"`
	Components *Components            `json:"components,omitempty" yaml:"components,omitempty"`
}

// Info contains API metadata.
type Info struct {
	Title   string `json:"title" yaml:"title"`
	Version string `json:"version" yaml:"version"`
}

// PathItem represents operations on a single path.
type PathItem struct {
	Get    *Operation `json:"get,omitempty" yaml:"get,omitempty"`
	Post   *Operation `json:"post,omitempty" yaml:"post,omitempty"`
	Put    *Operation `json:"put,omitempty" yaml:"put,omitempty"`
	Patch  *Operation `json:"patch,omitempty" yaml:"patch,omitempty"`
	Delete *Operation `json:"delete,omitempty" yaml:"delete,omitempty"`
}

// Operation represents a single API operation.
type Operation struct {
	OperationID string              `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Summary     string              `json:"summary,omitempty" yaml:"summary,omitempty"`
	Parameters  []Parameter         `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody *RequestBody        `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses   map[string]*Response `json:"responses,omitempty" yaml:"responses,omitempty"`
}

// Parameter represents an API parameter.
type Parameter struct {
	Name     string         `json:"name" yaml:"name"`
	In       string         `json:"in" yaml:"in"`
	Required bool           `json:"required,omitempty" yaml:"required,omitempty"`
	Schema   *schema.Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
}

// RequestBody represents a request body.
type RequestBody struct {
	Required bool                `json:"required,omitempty" yaml:"required,omitempty"`
	Content  map[string]*Content `json:"content,omitempty" yaml:"content,omitempty"`
}

// Response represents an API response.
type Response struct {
	Description string              `json:"description" yaml:"description"`
	Content     map[string]*Content `json:"content,omitempty" yaml:"content,omitempty"`
}

// Content represents media type content.
type Content struct {
	Schema *schema.Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
}

// Components holds reusable API components.
type Components struct {
	Schemas map[string]*schema.Schema `json:"schemas,omitempty" yaml:"schemas,omitempty"`
}

// ParseFile reads and parses an OpenAPI spec from a file.
func ParseFile(path string) (*Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading OpenAPI file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return ParseJSON(data)
	case ".yaml", ".yml":
		return ParseYAML(data)
	default:
		spec, err := ParseJSON(data)
		if err != nil {
			return ParseYAML(data)
		}
		return spec, nil
	}
}

// ParseJSON parses an OpenAPI spec from JSON bytes.
func ParseJSON(data []byte) (*Spec, error) {
	var spec Spec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("parsing OpenAPI JSON: %w", err)
	}
	if !strings.HasPrefix(spec.OpenAPI, "3.") {
		return nil, fmt.Errorf("unsupported OpenAPI version: %q (only 3.x supported)", spec.OpenAPI)
	}
	return &spec, nil
}

// ParseYAML parses an OpenAPI spec from YAML bytes.
func ParseYAML(data []byte) (*Spec, error) {
	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}

	jsonData, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("converting YAML to JSON: %w", err)
	}

	return ParseJSON(jsonData)
}
