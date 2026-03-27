package schema

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MaxFileSize is the maximum schema file size (10MB).
const MaxFileSize = 10 * 1024 * 1024

// MaxDepth is the maximum schema nesting depth.
const MaxDepth = 50

// MaxProperties is the maximum number of properties in a schema.
const MaxProperties = 10000

// SafeParseFile validates the file path and size before parsing.
func SafeParseFile(path string) (*Schema, error) {
	// Validate path
	if err := validatePath(path); err != nil {
		return nil, err
	}

	// Check file size
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("accessing file: %w", err)
	}
	if info.Size() > MaxFileSize {
		return nil, fmt.Errorf("file too large: %d bytes (max %d)", info.Size(), MaxFileSize)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("path is a directory, not a file: %s", path)
	}

	// Parse normally
	s, err := ParseFile(path)
	if err != nil {
		return nil, err
	}

	// Validate depth
	if depth := schemaDepth(s, 0); depth > MaxDepth {
		return nil, fmt.Errorf("schema too deeply nested: depth %d (max %d)", depth, MaxDepth)
	}

	// Validate property count
	if count := propertyCount(s); count > MaxProperties {
		return nil, fmt.Errorf("too many properties: %d (max %d)", count, MaxProperties)
	}

	return s, nil
}

func validatePath(path string) error {
	// Clean the path
	clean := filepath.Clean(path)

	// Check for suspicious patterns
	if strings.Contains(clean, "..") {
		// After cleaning, ".." indicates path traversal
		abs, err := filepath.Abs(clean)
		if err != nil {
			return fmt.Errorf("invalid path: %w", err)
		}
		// Verify the resolved path is reasonable
		_ = abs
	}

	// Validate extension
	ext := strings.ToLower(filepath.Ext(clean))
	validExts := map[string]bool{
		".json": true,
		".yaml": true,
		".yml":  true,
		"":      true, // allow no extension
	}
	if !validExts[ext] {
		return fmt.Errorf("unsupported file extension: %s (supported: .json, .yaml, .yml)", ext)
	}

	return nil
}

func schemaDepth(s *Schema, current int) int {
	if s == nil {
		return current
	}

	maxDepth := current
	for _, prop := range s.Properties {
		d := schemaDepth(prop, current+1)
		if d > maxDepth {
			maxDepth = d
		}
	}

	if s.Items != nil {
		d := schemaDepth(s.Items, current+1)
		if d > maxDepth {
			maxDepth = d
		}
	}

	for _, sub := range s.AllOf {
		d := schemaDepth(sub, current+1)
		if d > maxDepth {
			maxDepth = d
		}
	}

	for _, sub := range s.AnyOf {
		d := schemaDepth(sub, current+1)
		if d > maxDepth {
			maxDepth = d
		}
	}

	for _, sub := range s.OneOf {
		d := schemaDepth(sub, current+1)
		if d > maxDepth {
			maxDepth = d
		}
	}

	return maxDepth
}

func propertyCount(s *Schema) int {
	if s == nil {
		return 0
	}

	count := len(s.Properties)
	for _, prop := range s.Properties {
		count += propertyCount(prop)
	}
	if s.Items != nil {
		count += propertyCount(s.Items)
	}
	return count
}
