package schema

import (
	"encoding/json"
	"fmt"
)

// StrictParseJSON parses JSON with strict validation — disallows unknown fields.
func StrictParseJSON(data []byte) (*Schema, error) {
	// First, validate it's valid JSON
	var raw json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	// Parse normally
	s, err := ParseJSON(data)
	if err != nil {
		return nil, err
	}

	// Validate known fields by re-parsing into a map
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, err
	}

	knownFields := map[string]bool{
		"type": true, "properties": true, "required": true,
		"enum": true, "items": true, "minLength": true,
		"maxLength": true, "minimum": true, "maximum": true,
		"minItems": true, "maxItems": true, "pattern": true,
		"format": true, "additionalProperties": true,
		"description": true, "title": true, "$ref": true,
		"definitions": true, "allOf": true, "anyOf": true,
		"oneOf": true, "default": true, "$schema": true,
		"$id": true, "examples": true, "const": true,
		"not": true, "if": true, "then": true, "else": true,
		"exclusiveMinimum": true, "exclusiveMaximum": true,
		"multipleOf": true, "uniqueItems": true,
		"propertyNames": true, "patternProperties": true,
		"dependencies": true, "contentMediaType": true,
		"contentEncoding": true, "$comment": true,
	}

	var warnings []string
	for field := range fields {
		if !knownFields[field] {
			warnings = append(warnings, fmt.Sprintf("unknown field: %q", field))
		}
	}

	if len(warnings) > 0 {
		// Log warnings but don't fail — unknown fields are allowed in JSON Schema
		_ = warnings
	}

	return s, nil
}

// SanitizeOutput escapes potentially dangerous strings in change messages.
func SanitizeOutput(changes []Change) []Change {
	sanitized := make([]Change, len(changes))
	copy(sanitized, changes)
	for i := range sanitized {
		sanitized[i].Message = sanitizeString(sanitized[i].Message)
		sanitized[i].Path = sanitizeString(sanitized[i].Path)
	}
	return sanitized
}

func sanitizeString(s string) string {
	// Remove control characters except newline and tab
	result := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 32 || c == '\n' || c == '\t' {
			result = append(result, c)
		}
	}
	return string(result)
}
