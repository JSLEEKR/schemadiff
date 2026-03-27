package schema

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseJSON(t *testing.T) {
	data := []byte(`{
		"type": "object",
		"properties": {
			"name": {"type": "string"},
			"age": {"type": "integer"}
		},
		"required": ["name"]
	}`)

	s, err := ParseJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	if s.Type != "object" {
		t.Errorf("expected type=object, got %v", s.Type)
	}
	if len(s.Properties) != 2 {
		t.Errorf("expected 2 properties, got %d", len(s.Properties))
	}
	if len(s.Required) != 1 || s.Required[0] != "name" {
		t.Errorf("expected required=[name], got %v", s.Required)
	}
}

func TestParseJSONInvalid(t *testing.T) {
	_, err := ParseJSON([]byte(`{invalid`))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestParseYAML(t *testing.T) {
	data := []byte(`
type: object
properties:
  name:
    type: string
  age:
    type: integer
required:
  - name
`)

	s, err := ParseYAML(data)
	if err == nil && s != nil {
		if s.Type != "object" {
			t.Errorf("expected type=object, got %v", s.Type)
		}
	}
}

func TestParseYAMLInvalid(t *testing.T) {
	_, err := ParseYAML([]byte(`\tinvalid: [yaml`))
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestParseFile(t *testing.T) {
	dir := t.TempDir()

	// JSON file
	jsonPath := filepath.Join(dir, "schema.json")
	jsonData := []byte(`{"type":"object","properties":{"id":{"type":"string"}}}`)
	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		t.Fatal(err)
	}

	s, err := ParseFile(jsonPath)
	if err != nil {
		t.Fatal(err)
	}
	if s.Type != "object" {
		t.Errorf("expected type=object, got %v", s.Type)
	}

	// YAML file
	yamlPath := filepath.Join(dir, "schema.yaml")
	yamlData := []byte("type: string\nminLength: 1\n")
	if err := os.WriteFile(yamlPath, yamlData, 0644); err != nil {
		t.Fatal(err)
	}

	s2, err := ParseFile(yamlPath)
	if err != nil {
		t.Fatal(err)
	}
	if s2.Type != "string" {
		t.Errorf("expected type=string, got %v", s2.Type)
	}
}

func TestParseFileNotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/file.json")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestParseFileUnknownExt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "schema.txt")
	data := []byte(`{"type":"string"}`)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	s, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if s.Type != "string" {
		t.Errorf("expected type=string, got %v", s.Type)
	}
}

func TestResolveRefs(t *testing.T) {
	s := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"address": {Ref: "#/definitions/Address"},
		},
		Definitions: map[string]*Schema{
			"Address": {
				Type: "object",
				Properties: map[string]*Schema{
					"street": {Type: "string"},
					"city":   {Type: "string"},
				},
			},
		},
	}

	resolved := ResolveRefs(s)
	addr := resolved.Properties["address"]
	if addr == nil {
		t.Fatal("expected address property")
	}
	if addr.Type != "object" {
		t.Errorf("expected resolved type=object, got %v", addr.Type)
	}
	if len(addr.Properties) != 2 {
		t.Errorf("expected 2 resolved properties, got %d", len(addr.Properties))
	}
}

func TestResolveRefsCircular(t *testing.T) {
	s := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"self": {Ref: "#/definitions/Self"},
		},
		Definitions: map[string]*Schema{
			"Self": {
				Type: "object",
				Properties: map[string]*Schema{
					"nested": {Ref: "#/definitions/Self"},
				},
			},
		},
	}

	// Should not panic on circular refs
	resolved := ResolveRefs(s)
	if resolved == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestResolveRefsNil(t *testing.T) {
	if result := ResolveRefs(nil); result != nil {
		t.Error("expected nil for nil input")
	}
}

func TestResolveRefsNoRef(t *testing.T) {
	s := &Schema{Type: "string"}
	resolved := ResolveRefs(s)
	if resolved.Type != "string" {
		t.Errorf("expected type=string, got %v", resolved.Type)
	}
}
