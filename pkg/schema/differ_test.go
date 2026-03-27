package schema

import (
	"testing"
)

func intPtr(i int) *int         { return &i }
func floatPtr(f float64) *float64 { return &f }
func boolPtr(b bool) *bool      { return &b }

// --- Type changes ---

func TestDiffTypeChanged(t *testing.T) {
	old := &Schema{Type: "string"}
	new := &Schema{Type: "number"}
	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking change for type change")
	}
	assertChangeType(t, result, ChangeTypeChanged)
}

func TestDiffTypeSame(t *testing.T) {
	old := &Schema{Type: "string"}
	new := &Schema{Type: "string"}
	result := Diff(old, new)
	if result.HasBreaking {
		t.Error("expected no breaking changes")
	}
}

// --- Required fields ---

func TestDiffRequiredAdded(t *testing.T) {
	old := &Schema{Type: "object", Required: []string{"name"}}
	new := &Schema{Type: "object", Required: []string{"name", "email"}}
	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking change for required field added")
	}
	assertChangeType(t, result, ChangeRequiredAdded)
}

func TestDiffRequiredRemoved(t *testing.T) {
	old := &Schema{Type: "object", Required: []string{"name", "email"}}
	new := &Schema{Type: "object", Required: []string{"name"}}
	result := Diff(old, new)
	if result.HasBreaking {
		t.Error("expected no breaking change for required removed")
	}
	assertChangeType(t, result, ChangeRequiredRemoved)
}

// --- Properties ---

func TestDiffPropertyRemoved(t *testing.T) {
	old := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"name":  {Type: "string"},
			"email": {Type: "string"},
		},
	}
	new := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"name": {Type: "string"},
		},
	}
	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking change for property removed")
	}
	assertChangeType(t, result, ChangePropertyRemoved)
}

func TestDiffPropertyAdded(t *testing.T) {
	old := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"name": {Type: "string"},
		},
	}
	new := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"name":  {Type: "string"},
			"email": {Type: "string"},
		},
	}
	result := Diff(old, new)
	if result.HasBreaking {
		t.Error("expected no breaking change for property added")
	}
	assertChangeType(t, result, ChangePropertyAdded)
}

func TestDiffNestedPropertyTypeChanged(t *testing.T) {
	old := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"address": {
				Type: "object",
				Properties: map[string]*Schema{
					"zip": {Type: "string"},
				},
			},
		},
	}
	new := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"address": {
				Type: "object",
				Properties: map[string]*Schema{
					"zip": {Type: "integer"},
				},
			},
		},
	}
	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking change for nested type change")
	}
}

// --- Enum changes ---

func TestDiffEnumValueRemoved(t *testing.T) {
	old := &Schema{Type: "string", Enum: []interface{}{"a", "b", "c"}}
	new := &Schema{Type: "string", Enum: []interface{}{"a", "b"}}
	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking change for enum value removed")
	}
	assertChangeType(t, result, ChangeEnumValueRemoved)
}

func TestDiffEnumValueAdded(t *testing.T) {
	old := &Schema{Type: "string", Enum: []interface{}{"a", "b"}}
	new := &Schema{Type: "string", Enum: []interface{}{"a", "b", "c"}}
	result := Diff(old, new)
	if result.HasBreaking {
		t.Error("expected no breaking change for enum value added")
	}
	assertChangeType(t, result, ChangeEnumValueAdded)
}

// --- Constraint changes ---

func TestDiffMinLengthIncreased(t *testing.T) {
	old := &Schema{Type: "string", MinLength: intPtr(1)}
	new := &Schema{Type: "string", MinLength: intPtr(5)}
	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking for minLength increase")
	}
	assertChangeType(t, result, ChangeMinLengthIncreased)
}

func TestDiffMinLengthAdded(t *testing.T) {
	old := &Schema{Type: "string"}
	new := &Schema{Type: "string", MinLength: intPtr(5)}
	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking for minLength added")
	}
}

func TestDiffMaxLengthDecreased(t *testing.T) {
	old := &Schema{Type: "string", MaxLength: intPtr(100)}
	new := &Schema{Type: "string", MaxLength: intPtr(50)}
	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking for maxLength decrease")
	}
	assertChangeType(t, result, ChangeMaxLengthDecreased)
}

func TestDiffMaxLengthAdded(t *testing.T) {
	old := &Schema{Type: "string"}
	new := &Schema{Type: "string", MaxLength: intPtr(50)}
	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking for maxLength added")
	}
}

func TestDiffMinimumIncreased(t *testing.T) {
	old := &Schema{Type: "number", Minimum: floatPtr(0)}
	new := &Schema{Type: "number", Minimum: floatPtr(10)}
	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking for minimum increase")
	}
	assertChangeType(t, result, ChangeMinimumIncreased)
}

func TestDiffMinimumAdded(t *testing.T) {
	old := &Schema{Type: "number"}
	new := &Schema{Type: "number", Minimum: floatPtr(10)}
	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking for minimum added")
	}
}

func TestDiffMaximumDecreased(t *testing.T) {
	old := &Schema{Type: "number", Maximum: floatPtr(100)}
	new := &Schema{Type: "number", Maximum: floatPtr(50)}
	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking for maximum decrease")
	}
	assertChangeType(t, result, ChangeMaximumDecreased)
}

func TestDiffMaximumAdded(t *testing.T) {
	old := &Schema{Type: "number"}
	new := &Schema{Type: "number", Maximum: floatPtr(50)}
	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking for maximum added")
	}
}

func TestDiffMinItemsIncreased(t *testing.T) {
	old := &Schema{Type: "array", MinItems: intPtr(1)}
	new := &Schema{Type: "array", MinItems: intPtr(3)}
	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking for minItems increase")
	}
	assertChangeType(t, result, ChangeMinItemsIncreased)
}

func TestDiffMinItemsAdded(t *testing.T) {
	old := &Schema{Type: "array"}
	new := &Schema{Type: "array", MinItems: intPtr(3)}
	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking for minItems added")
	}
}

func TestDiffMaxItemsDecreased(t *testing.T) {
	old := &Schema{Type: "array", MaxItems: intPtr(10)}
	new := &Schema{Type: "array", MaxItems: intPtr(5)}
	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking for maxItems decrease")
	}
	assertChangeType(t, result, ChangeMaxItemsDecreased)
}

func TestDiffMaxItemsAdded(t *testing.T) {
	old := &Schema{Type: "array"}
	new := &Schema{Type: "array", MaxItems: intPtr(5)}
	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking for maxItems added")
	}
}

// --- Pattern ---

func TestDiffPatternChanged(t *testing.T) {
	old := &Schema{Type: "string", Pattern: "^[a-z]+$"}
	new := &Schema{Type: "string", Pattern: "^[A-Z]+$"}
	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking for pattern change")
	}
	assertChangeType(t, result, ChangePatternChanged)
}

func TestDiffPatternAdded(t *testing.T) {
	old := &Schema{Type: "string"}
	new := &Schema{Type: "string", Pattern: "^[a-z]+$"}
	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking for pattern added")
	}
}

// --- Format ---

func TestDiffFormatChanged(t *testing.T) {
	old := &Schema{Type: "string", Format: "email"}
	new := &Schema{Type: "string", Format: "uri"}
	result := Diff(old, new)
	if result.Summary.Warning != 1 {
		t.Error("expected warning for format change")
	}
	assertChangeType(t, result, ChangeFormatChanged)
}

// --- AdditionalProperties ---

func TestDiffAdditionalPropertiesRestricted(t *testing.T) {
	old := &Schema{Type: "object"}
	new := &Schema{Type: "object", AdditionalProperties: boolPtr(false)}
	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking for additionalProperties restriction")
	}
	assertChangeType(t, result, ChangeAdditionalPropsFalse)
}

// --- Items ---

func TestDiffItemsTypeChanged(t *testing.T) {
	old := &Schema{Type: "array", Items: &Schema{Type: "string"}}
	new := &Schema{Type: "array", Items: &Schema{Type: "number"}}
	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking for items type change")
	}
}

func TestDiffItemsRemoved(t *testing.T) {
	old := &Schema{Type: "array", Items: &Schema{Type: "string"}}
	new := &Schema{Type: "array"}
	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking for items removal")
	}
}

// --- Edge cases ---

func TestDiffBothNil(t *testing.T) {
	result := Diff(nil, nil)
	if result.HasBreaking {
		t.Error("expected no breaking for nil schemas")
	}
}

func TestDiffOldNil(t *testing.T) {
	result := Diff(nil, &Schema{Type: "string"})
	if result.HasBreaking {
		t.Error("expected no breaking for new schema")
	}
}

func TestDiffNewNil(t *testing.T) {
	result := Diff(&Schema{Type: "string"}, nil)
	if !result.HasBreaking {
		t.Error("expected breaking for removed schema")
	}
}

func TestDiffNoChanges(t *testing.T) {
	s := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"name": {Type: "string"},
		},
		Required: []string{"name"},
	}
	result := Diff(s, s)
	if result.Summary.Total != 0 {
		t.Errorf("expected no changes, got %d", result.Summary.Total)
	}
}

func TestDiffConstraintsNotStricter(t *testing.T) {
	// minLength decreased (loosened) — not breaking
	old := &Schema{Type: "string", MinLength: intPtr(5)}
	new := &Schema{Type: "string", MinLength: intPtr(1)}
	result := Diff(old, new)
	if result.HasBreaking {
		t.Error("expected no breaking for minLength decrease")
	}

	// maxLength increased (loosened) — not breaking
	old2 := &Schema{Type: "string", MaxLength: intPtr(50)}
	new2 := &Schema{Type: "string", MaxLength: intPtr(100)}
	result2 := Diff(old2, new2)
	if result2.HasBreaking {
		t.Error("expected no breaking for maxLength increase")
	}
}

func TestDiffComplexSchema(t *testing.T) {
	old := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"name":  {Type: "string", MinLength: intPtr(1), MaxLength: intPtr(100)},
			"age":   {Type: "integer", Minimum: floatPtr(0), Maximum: floatPtr(150)},
			"email": {Type: "string", Format: "email"},
			"role":  {Type: "string", Enum: []interface{}{"admin", "user", "guest"}},
		},
		Required: []string{"name", "email"},
	}

	new := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"name":  {Type: "string", MinLength: intPtr(2), MaxLength: intPtr(50)},
			"age":   {Type: "string"}, // type changed
			"role":  {Type: "string", Enum: []interface{}{"admin", "user"}}, // guest removed
			"phone": {Type: "string"}, // added
		},
		Required: []string{"name", "email", "phone"}, // phone added
	}

	result := Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking changes")
	}
	// Breaking: minLength increased, maxLength decreased, age type changed,
	//           email removed, guest enum removed, phone required added
	if result.Summary.Breaking < 5 {
		t.Errorf("expected at least 5 breaking changes, got %d", result.Summary.Breaking)
	}
}

// --- Helpers ---

func assertChangeType(t *testing.T, result *DiffResult, ct ChangeType) {
	t.Helper()
	for _, c := range result.Changes {
		if c.Type == ct {
			return
		}
	}
	t.Errorf("expected change type %q not found in results", ct)
}
