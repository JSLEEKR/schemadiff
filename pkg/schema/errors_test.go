package schema

import (
	"errors"
	"fmt"
	"testing"
)

func TestSchemaErrorWithFile(t *testing.T) {
	err := &SchemaError{
		File:    "test.json",
		Err:     ErrParseFailure,
		Message: "parse failed",
	}
	got := err.Error()
	if got != "parse failed (file: test.json): schema parse failure" {
		t.Errorf("unexpected error: %s", got)
	}
}

func TestSchemaErrorWithFileAndPath(t *testing.T) {
	err := &SchemaError{
		File:    "test.json",
		Path:    "properties.name",
		Err:     ErrInvalidSchema,
		Message: "invalid",
	}
	got := err.Error()
	if got != "invalid (file: test.json, path: properties.name): invalid schema" {
		t.Errorf("unexpected error: %s", got)
	}
}

func TestSchemaErrorNoFile(t *testing.T) {
	err := &SchemaError{
		Err:     ErrParseFailure,
		Message: "parse failed",
	}
	got := err.Error()
	if got != "parse failed: schema parse failure" {
		t.Errorf("unexpected error: %s", got)
	}
}

func TestSchemaErrorUnwrap(t *testing.T) {
	inner := fmt.Errorf("inner error")
	err := &SchemaError{Err: inner, Message: "outer"}
	if !errors.Is(err, inner) {
		t.Error("expected Unwrap to return inner error")
	}
}

func TestNewParseError(t *testing.T) {
	err := NewParseError("test.json", fmt.Errorf("bad json"))
	if !IsParseError(err) {
		t.Error("expected parse error")
	}
	if err.File != "test.json" {
		t.Error("expected file")
	}
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("test.json", "type", fmt.Errorf("invalid"))
	if !IsValidationError(err) {
		t.Error("expected validation error")
	}
	if err.Path != "type" {
		t.Error("expected path")
	}
}

func TestIsParseError(t *testing.T) {
	if IsParseError(fmt.Errorf("random error")) {
		t.Error("expected false for non-parse error")
	}
}

func TestIsValidationError(t *testing.T) {
	if IsValidationError(fmt.Errorf("random error")) {
		t.Error("expected false for non-validation error")
	}
}

func TestExitCodes(t *testing.T) {
	if ExitCompatible != 0 {
		t.Error("expected ExitCompatible=0")
	}
	if ExitBreaking != 1 {
		t.Error("expected ExitBreaking=1")
	}
	if ExitError != 2 {
		t.Error("expected ExitError=2")
	}
}

func TestErrorSentinels(t *testing.T) {
	sentinels := []error{
		ErrParseFailure,
		ErrInvalidSchema,
		ErrFileNotFound,
		ErrFileTooLarge,
		ErrUnsupportedFormat,
		ErrDepthExceeded,
		ErrPropertyLimit,
	}
	for _, s := range sentinels {
		if s == nil {
			t.Error("expected non-nil sentinel error")
		}
		if s.Error() == "" {
			t.Error("expected non-empty error message")
		}
	}
}
