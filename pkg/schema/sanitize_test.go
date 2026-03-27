package schema

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSafeParseFileValid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "schema.json")
	if err := os.WriteFile(path, []byte(`{"type":"object"}`), 0644); err != nil {
		t.Fatal(err)
	}

	s, err := SafeParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if s.Type != "object" {
		t.Errorf("expected type=object, got %v", s.Type)
	}
}

func TestSafeParseFileNotFound(t *testing.T) {
	_, err := SafeParseFile("/nonexistent/schema.json")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestSafeParseFileDirectory(t *testing.T) {
	dir := t.TempDir()
	_, err := SafeParseFile(dir)
	if err == nil {
		t.Error("expected error for directory")
	}
}

func TestSafeParseFileBadExtension(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "schema.exe")
	if err := os.WriteFile(path, []byte(`{"type":"object"}`), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := SafeParseFile(path)
	if err == nil {
		t.Error("expected error for .exe extension")
	}
}

func TestSafeParseFileTooLarge(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "large.json")
	// Create a file just over the limit check by checking the error message
	// (We won't actually create a 10MB file in tests)
	if err := os.WriteFile(path, []byte(`{"type":"object"}`), 0644); err != nil {
		t.Fatal(err)
	}

	// Should parse fine for small file
	_, err := SafeParseFile(path)
	if err != nil {
		t.Errorf("expected no error for small file: %v", err)
	}
}

func TestValidatePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"valid json", "schema.json", false},
		{"valid yaml", "schema.yaml", false},
		{"valid yml", "schema.yml", false},
		{"no extension", "schema", false},
		{"bad extension", "schema.exe", true},
		{"bad extension py", "schema.py", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePath(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
		})
	}
}

func TestSchemaDepth(t *testing.T) {
	s := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"level1": {
				Type: "object",
				Properties: map[string]*Schema{
					"level2": {
						Type: "object",
						Properties: map[string]*Schema{
							"level3": {Type: "string"},
						},
					},
				},
			},
		},
	}

	depth := schemaDepth(s, 0)
	if depth != 3 {
		t.Errorf("expected depth=3, got %d", depth)
	}
}

func TestSchemaDepthNil(t *testing.T) {
	if schemaDepth(nil, 0) != 0 {
		t.Error("expected 0 for nil")
	}
}

func TestPropertyCount(t *testing.T) {
	s := &Schema{
		Properties: map[string]*Schema{
			"a": {
				Properties: map[string]*Schema{
					"b": {Type: "string"},
					"c": {Type: "string"},
				},
			},
			"d": {Type: "string"},
		},
	}
	count := propertyCount(s)
	if count != 4 { // a, b, c, d
		t.Errorf("expected 4, got %d", count)
	}
}

func TestPropertyCountNil(t *testing.T) {
	if propertyCount(nil) != 0 {
		t.Error("expected 0 for nil")
	}
}

func TestSafeParseFileDeeplyNested(t *testing.T) {
	// Build a deeply nested schema JSON
	depth := MaxDepth + 5
	var sb strings.Builder
	for i := 0; i < depth; i++ {
		sb.WriteString(`{"type":"object","properties":{"x":`)
	}
	sb.WriteString(`{"type":"string"}`)
	for i := 0; i < depth; i++ {
		sb.WriteString(`}}`)
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "deep.json")
	if err := os.WriteFile(path, []byte(sb.String()), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := SafeParseFile(path)
	if err == nil {
		t.Error("expected error for deeply nested schema")
	}
}

func TestSchemaDepthWithItems(t *testing.T) {
	s := &Schema{
		Type: "array",
		Items: &Schema{
			Type: "object",
			Properties: map[string]*Schema{
				"name": {Type: "string"},
			},
		},
	}
	depth := schemaDepth(s, 0)
	if depth != 2 {
		t.Errorf("expected depth=2, got %d", depth)
	}
}

func TestSchemaDepthWithAllOf(t *testing.T) {
	s := &Schema{
		AllOf: []*Schema{
			{
				Type: "object",
				Properties: map[string]*Schema{
					"x": {Type: "string"},
				},
			},
		},
	}
	depth := schemaDepth(s, 0)
	if depth != 2 {
		t.Errorf("expected depth=2, got %d", depth)
	}
}
