package schema

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadFromStdinOrFileFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.json")
	if err := os.WriteFile(path, []byte(`{"type":"string"}`), 0644); err != nil {
		t.Fatal(err)
	}

	data, err := ReadFromStdinOrFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"type":"string"}` {
		t.Errorf("unexpected data: %s", data)
	}
}

func TestReadFromStdinOrFileNotFound(t *testing.T) {
	_, err := ReadFromStdinOrFile("/nonexistent/file.json")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestParseFromReader(t *testing.T) {
	r := strings.NewReader(`{"type":"object","properties":{"name":{"type":"string"}}}`)
	s, err := ParseFromReader(r)
	if err != nil {
		t.Fatal(err)
	}
	if s.Type != "object" {
		t.Errorf("expected type=object, got %v", s.Type)
	}
}

func TestParseFromReaderYAML(t *testing.T) {
	r := strings.NewReader("type: string\nminLength: 1\n")
	s, err := ParseFromReader(r)
	if err == nil && s != nil {
		if s.Type != "string" {
			t.Errorf("expected type=string, got %v", s.Type)
		}
	}
}

func TestParseFromReaderEmpty(t *testing.T) {
	r := strings.NewReader("")
	_, err := ParseFromReader(r)
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestParseFromReaderBytes(t *testing.T) {
	r := bytes.NewReader([]byte(`{"type":"integer"}`))
	s, err := ParseFromReader(r)
	if err != nil {
		t.Fatal(err)
	}
	if s.Type != "integer" {
		t.Errorf("expected type=integer, got %v", s.Type)
	}
}
