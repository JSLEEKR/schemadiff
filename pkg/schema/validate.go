package schema

import (
	"fmt"
	"strings"
)

// ValidationError represents a schema validation issue.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidateSchema checks a schema for basic structural validity.
func ValidateSchema(s *Schema) []ValidationError {
	var errs []ValidationError

	if s == nil {
		return []ValidationError{{Field: "(root)", Message: "schema is nil"}}
	}

	// Check type is valid
	if s.Type != nil {
		validTypes := map[string]bool{
			"string": true, "number": true, "integer": true,
			"boolean": true, "object": true, "array": true, "null": true,
		}
		t := normalizeType(s.Type)
		for _, part := range strings.Split(t, ",") {
			part = strings.TrimSpace(part)
			if part != "" && !validTypes[part] {
				errs = append(errs, ValidationError{
					Field:   "type",
					Message: fmt.Sprintf("invalid type %q", part),
				})
			}
		}
	}

	// Check minLength/maxLength consistency
	if s.MinLength != nil && s.MaxLength != nil && *s.MinLength > *s.MaxLength {
		errs = append(errs, ValidationError{
			Field:   "minLength/maxLength",
			Message: fmt.Sprintf("minLength (%d) > maxLength (%d)", *s.MinLength, *s.MaxLength),
		})
	}

	// Check minimum/maximum consistency
	if s.Minimum != nil && s.Maximum != nil && *s.Minimum > *s.Maximum {
		errs = append(errs, ValidationError{
			Field:   "minimum/maximum",
			Message: fmt.Sprintf("minimum (%v) > maximum (%v)", *s.Minimum, *s.Maximum),
		})
	}

	// Check minItems/maxItems consistency
	if s.MinItems != nil && s.MaxItems != nil && *s.MinItems > *s.MaxItems {
		errs = append(errs, ValidationError{
			Field:   "minItems/maxItems",
			Message: fmt.Sprintf("minItems (%d) > maxItems (%d)", *s.MinItems, *s.MaxItems),
		})
	}

	// Validate nested properties
	for name, prop := range s.Properties {
		subErrs := ValidateSchema(prop)
		for _, e := range subErrs {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("properties.%s.%s", name, e.Field),
				Message: e.Message,
			})
		}
	}

	// Validate items schema
	if s.Items != nil {
		subErrs := ValidateSchema(s.Items)
		for _, e := range subErrs {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("items.%s", e.Field),
				Message: e.Message,
			})
		}
	}

	return errs
}
