package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ParseFile reads and parses a JSON Schema from a file.
// Supports both JSON (.json) and YAML (.yaml, .yml) formats.
func ParseFile(path string) (*Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading schema file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return ParseJSON(data)
	case ".yaml", ".yml":
		return ParseYAML(data)
	default:
		// Try JSON first, then YAML
		s, err := ParseJSON(data)
		if err != nil {
			return ParseYAML(data)
		}
		return s, nil
	}
}

// ParseJSON parses a JSON Schema from JSON bytes.
func ParseJSON(data []byte) (*Schema, error) {
	var s Schema
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing JSON schema: %w", err)
	}
	return &s, nil
}

// ParseYAML parses a JSON Schema from YAML bytes.
func ParseYAML(data []byte) (*Schema, error) {
	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}

	// Convert to JSON then parse as Schema for consistency
	jsonData, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("converting YAML to JSON: %w", err)
	}

	return ParseJSON(jsonData)
}

// ResolveRefs resolves $ref pointers within a schema using its definitions.
func ResolveRefs(s *Schema) *Schema {
	if s == nil {
		return nil
	}
	return resolveRefsInternal(s, s.Definitions, make(map[string]bool))
}

func resolveRefsInternal(s *Schema, definitions map[string]*Schema, seen map[string]bool) *Schema {
	if s == nil {
		return nil
	}

	// Handle $ref
	if s.Ref != "" {
		refName := extractRefName(s.Ref)
		if seen[refName] {
			// Circular reference, return as-is
			return s
		}
		if def, ok := definitions[refName]; ok {
			seen[refName] = true
			resolved := resolveRefsInternal(def, definitions, seen)
			delete(seen, refName)
			return resolved
		}
		return s
	}

	// Deep copy and resolve properties
	resolved := *s
	if s.Properties != nil {
		resolved.Properties = make(map[string]*Schema, len(s.Properties))
		for k, v := range s.Properties {
			resolved.Properties[k] = resolveRefsInternal(v, definitions, seen)
		}
	}

	if s.Items != nil {
		resolved.Items = resolveRefsInternal(s.Items, definitions, seen)
	}

	if s.AllOf != nil {
		resolved.AllOf = make([]*Schema, len(s.AllOf))
		for i, v := range s.AllOf {
			resolved.AllOf[i] = resolveRefsInternal(v, definitions, seen)
		}
	}

	if s.AnyOf != nil {
		resolved.AnyOf = make([]*Schema, len(s.AnyOf))
		for i, v := range s.AnyOf {
			resolved.AnyOf[i] = resolveRefsInternal(v, definitions, seen)
		}
	}

	if s.OneOf != nil {
		resolved.OneOf = make([]*Schema, len(s.OneOf))
		for i, v := range s.OneOf {
			resolved.OneOf[i] = resolveRefsInternal(v, definitions, seen)
		}
	}

	return &resolved
}

func extractRefName(ref string) string {
	// Handle #/definitions/Name format
	parts := strings.Split(ref, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ref
}
