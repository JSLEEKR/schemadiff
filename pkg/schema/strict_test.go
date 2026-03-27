package schema

import "testing"

func TestStrictParseJSONValid(t *testing.T) {
	data := []byte(`{"type":"object","properties":{"name":{"type":"string"}}}`)
	s, err := StrictParseJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	if s.Type != "object" {
		t.Errorf("expected type=object, got %v", s.Type)
	}
}

func TestStrictParseJSONInvalid(t *testing.T) {
	_, err := StrictParseJSON([]byte(`{invalid`))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestStrictParseJSONWithUnknownFields(t *testing.T) {
	data := []byte(`{"type":"string","customExtension":"value","x-vendor":"data"}`)
	s, err := StrictParseJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	if s.Type != "string" {
		t.Errorf("expected type=string, got %v", s.Type)
	}
}

func TestStrictParseJSONAllKnownFields(t *testing.T) {
	data := []byte(`{
		"type": "object",
		"properties": {},
		"required": [],
		"enum": [],
		"items": {},
		"minLength": 0,
		"maxLength": 100,
		"minimum": 0,
		"maximum": 100,
		"minItems": 0,
		"maxItems": 100,
		"pattern": ".*",
		"format": "email",
		"additionalProperties": true,
		"description": "test",
		"title": "test",
		"$ref": "",
		"definitions": {},
		"allOf": [],
		"anyOf": [],
		"oneOf": [],
		"default": null
	}`)
	_, err := StrictParseJSON(data)
	if err != nil {
		t.Fatalf("expected no error for known fields: %v", err)
	}
}

func TestSanitizeOutput(t *testing.T) {
	changes := []Change{
		{
			Path:    "normal.path",
			Message: "Normal message",
		},
		{
			Path:    "path\x00with\x01control",
			Message: "message\x02with\x03chars",
		},
	}

	sanitized := SanitizeOutput(changes)
	if sanitized[0].Path != "normal.path" {
		t.Error("expected normal path unchanged")
	}
	if sanitized[0].Message != "Normal message" {
		t.Error("expected normal message unchanged")
	}
	if sanitized[1].Path != "pathwithcontrol" {
		t.Errorf("expected control chars removed from path, got %q", sanitized[1].Path)
	}
	if sanitized[1].Message != "messagewithchars" {
		t.Errorf("expected control chars removed from message, got %q", sanitized[1].Message)
	}
}

func TestSanitizeOutputPreservesNewlines(t *testing.T) {
	changes := []Change{
		{Message: "line1\nline2\ttab"},
	}
	sanitized := SanitizeOutput(changes)
	if sanitized[0].Message != "line1\nline2\ttab" {
		t.Error("expected newlines and tabs preserved")
	}
}

func TestSanitizeOutputEmpty(t *testing.T) {
	sanitized := SanitizeOutput(nil)
	if len(sanitized) != 0 {
		t.Error("expected empty for nil input")
	}
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"normal", "normal"},
		{"with\x00null", "withnull"},
		{"tabs\tok", "tabs\tok"},
		{"newline\nok", "newline\nok"},
		{"\x01\x02\x03", ""},
		{"", ""},
	}
	for _, tt := range tests {
		got := sanitizeString(tt.input)
		if got != tt.want {
			t.Errorf("sanitizeString(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
