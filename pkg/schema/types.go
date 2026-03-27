// Package schema provides JSON Schema parsing, diffing, and breaking change detection.
package schema

import "fmt"

// Severity indicates the impact level of a schema change.
type Severity int

const (
	// SeverityInfo is a non-breaking, informational change.
	SeverityInfo Severity = iota
	// SeverityWarning is a potentially breaking change.
	SeverityWarning
	// SeverityBreaking is a confirmed backward-incompatible change.
	SeverityBreaking
)

// String returns the human-readable severity label.
func (s Severity) String() string {
	switch s {
	case SeverityInfo:
		return "INFO"
	case SeverityWarning:
		return "WARNING"
	case SeverityBreaking:
		return "BREAKING"
	default:
		return "UNKNOWN"
	}
}

// ChangeType identifies the kind of schema change detected.
type ChangeType string

const (
	ChangeRequiredAdded       ChangeType = "required-field-added"
	ChangeRequiredRemoved     ChangeType = "required-field-removed"
	ChangePropertyRemoved     ChangeType = "property-removed"
	ChangePropertyAdded       ChangeType = "property-added"
	ChangeTypeChanged         ChangeType = "type-changed"
	ChangeEnumValueRemoved    ChangeType = "enum-value-removed"
	ChangeEnumValueAdded      ChangeType = "enum-value-added"
	ChangeMinLengthIncreased  ChangeType = "min-length-increased"
	ChangeMaxLengthDecreased  ChangeType = "max-length-decreased"
	ChangeMinimumIncreased    ChangeType = "minimum-increased"
	ChangeMaximumDecreased    ChangeType = "maximum-decreased"
	ChangeMinItemsIncreased   ChangeType = "min-items-increased"
	ChangeMaxItemsDecreased   ChangeType = "max-items-decreased"
	ChangePatternChanged      ChangeType = "pattern-changed"
	ChangeFormatChanged       ChangeType = "format-changed"
	ChangeAdditionalPropsFalse ChangeType = "additional-properties-restricted"
	ChangeItemsChanged        ChangeType = "items-schema-changed"
)

// Change represents a single detected schema change.
type Change struct {
	// Path is the JSON path to the changed element (e.g., "properties.name.type").
	Path string `json:"path"`
	// Type is the kind of change detected.
	Type ChangeType `json:"type"`
	// Severity indicates whether the change is breaking.
	Severity Severity `json:"severity"`
	// Message is a human-readable description of the change.
	Message string `json:"message"`
	// OldValue is the previous value (if applicable).
	OldValue interface{} `json:"old_value,omitempty"`
	// NewValue is the new value (if applicable).
	NewValue interface{} `json:"new_value,omitempty"`
}

// String returns a human-readable change description.
func (c Change) String() string {
	return fmt.Sprintf("[%s] %s: %s", c.Severity, c.Path, c.Message)
}

// Schema represents a parsed JSON Schema.
type Schema struct {
	// Type is the JSON Schema type (string, number, object, array, boolean, integer, null).
	Type interface{} `json:"type,omitempty"`
	// Properties maps property names to their schemas (for object types).
	Properties map[string]*Schema `json:"properties,omitempty"`
	// Required lists required property names.
	Required []string `json:"required,omitempty"`
	// Enum lists allowed values.
	Enum []interface{} `json:"enum,omitempty"`
	// Items is the schema for array items.
	Items *Schema `json:"items,omitempty"`
	// MinLength is the minimum string length.
	MinLength *int `json:"minLength,omitempty"`
	// MaxLength is the maximum string length.
	MaxLength *int `json:"maxLength,omitempty"`
	// Minimum is the minimum numeric value.
	Minimum *float64 `json:"minimum,omitempty"`
	// Maximum is the maximum numeric value.
	Maximum *float64 `json:"maximum,omitempty"`
	// MinItems is the minimum array length.
	MinItems *int `json:"minItems,omitempty"`
	// MaxItems is the maximum array length.
	MaxItems *int `json:"maxItems,omitempty"`
	// Pattern is the regex pattern for strings.
	Pattern string `json:"pattern,omitempty"`
	// Format is the string format (e.g., "email", "date-time").
	Format string `json:"format,omitempty"`
	// AdditionalProperties controls whether extra properties are allowed.
	AdditionalProperties *bool `json:"additionalProperties,omitempty"`
	// Description is a human-readable description.
	Description string `json:"description,omitempty"`
	// Title is the schema title.
	Title string `json:"title,omitempty"`
	// Ref is a JSON Schema $ref pointer.
	Ref string `json:"$ref,omitempty"`
	// Definitions stores reusable schema definitions.
	Definitions map[string]*Schema `json:"definitions,omitempty"`
	// AllOf combines schemas with AND logic.
	AllOf []*Schema `json:"allOf,omitempty"`
	// AnyOf combines schemas with OR logic.
	AnyOf []*Schema `json:"anyOf,omitempty"`
	// OneOf combines schemas with XOR logic.
	OneOf []*Schema `json:"oneOf,omitempty"`
	// Default is the default value.
	Default interface{} `json:"default,omitempty"`
}

// DiffResult holds the complete result of comparing two schemas.
type DiffResult struct {
	// Changes is the list of all detected changes.
	Changes []Change `json:"changes"`
	// HasBreaking is true if any breaking change was detected.
	HasBreaking bool `json:"has_breaking"`
	// Summary provides counts by severity.
	Summary DiffSummary `json:"summary"`
}

// DiffSummary provides aggregate counts of changes by severity.
type DiffSummary struct {
	Breaking    int `json:"breaking"`
	Warning     int `json:"warning"`
	Info        int `json:"info"`
	Total       int `json:"total"`
}

// NewDiffResult creates a DiffResult from a list of changes.
func NewDiffResult(changes []Change) *DiffResult {
	result := &DiffResult{
		Changes: changes,
	}
	for _, c := range changes {
		switch c.Severity {
		case SeverityBreaking:
			result.Summary.Breaking++
			result.HasBreaking = true
		case SeverityWarning:
			result.Summary.Warning++
		case SeverityInfo:
			result.Summary.Info++
		}
	}
	result.Summary.Total = len(changes)
	return result
}
