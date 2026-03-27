package schema

import (
	"errors"
	"fmt"
)

// Exit codes for CLI usage.
const (
	ExitCompatible = 0 // No breaking changes
	ExitBreaking   = 1 // Breaking changes detected
	ExitError      = 2 // Invalid input or error
)

// Error types for structured error handling.
var (
	ErrParseFailure     = errors.New("schema parse failure")
	ErrInvalidSchema    = errors.New("invalid schema")
	ErrFileNotFound     = errors.New("file not found")
	ErrFileTooLarge     = errors.New("file too large")
	ErrUnsupportedFormat = errors.New("unsupported format")
	ErrDepthExceeded    = errors.New("schema depth exceeded")
	ErrPropertyLimit    = errors.New("property count exceeded")
)

// SchemaError wraps errors with schema context.
type SchemaError struct {
	File    string
	Path    string
	Err     error
	Message string
}

func (e *SchemaError) Error() string {
	if e.File != "" && e.Path != "" {
		return fmt.Sprintf("%s (file: %s, path: %s): %s", e.Message, e.File, e.Path, e.Err)
	}
	if e.File != "" {
		return fmt.Sprintf("%s (file: %s): %s", e.Message, e.File, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Message, e.Err)
}

func (e *SchemaError) Unwrap() error {
	return e.Err
}

// NewParseError creates a parse failure error.
func NewParseError(file string, err error) *SchemaError {
	return &SchemaError{
		File:    file,
		Err:     fmt.Errorf("%w: %v", ErrParseFailure, err),
		Message: "failed to parse schema",
	}
}

// NewValidationError creates a validation error.
func NewValidationError(file, path string, err error) *SchemaError {
	return &SchemaError{
		File:    file,
		Path:    path,
		Err:     fmt.Errorf("%w: %v", ErrInvalidSchema, err),
		Message: "schema validation failed",
	}
}

// IsParseError checks if an error is a parse failure.
func IsParseError(err error) bool {
	return errors.Is(err, ErrParseFailure)
}

// IsValidationError checks if an error is a validation error.
func IsValidationError(err error) bool {
	return errors.Is(err, ErrInvalidSchema)
}
