package schema

import (
	"strings"
	"testing"
)

func TestCollectStats(t *testing.T) {
	s := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"name":  {Type: "string"},
			"age":   {Type: "integer"},
			"role":  {Type: "string", Enum: []interface{}{"admin", "user"}},
			"nested": {
				Type: "object",
				Properties: map[string]*Schema{
					"x": {Type: "string"},
				},
			},
		},
		Required: []string{"name", "age"},
	}

	stats := CollectStats(s)
	if stats.PropertyCount != 5 { // name, age, role, nested, x
		t.Errorf("expected 5 properties, got %d", stats.PropertyCount)
	}
	if stats.RequiredCount != 2 {
		t.Errorf("expected 2 required, got %d", stats.RequiredCount)
	}
	if stats.Depth != 2 {
		t.Errorf("expected depth=2, got %d", stats.Depth)
	}
	if !stats.HasEnum {
		t.Error("expected HasEnum=true")
	}
	if stats.Type != "object" {
		t.Errorf("expected type=object, got %s", stats.Type)
	}
}

func TestCollectStatsNil(t *testing.T) {
	stats := CollectStats(nil)
	if stats.PropertyCount != 0 {
		t.Error("expected 0 properties for nil")
	}
}

func TestCollectStatsWithRefs(t *testing.T) {
	s := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"addr": {Ref: "#/definitions/Address"},
		},
	}
	stats := CollectStats(s)
	if !stats.HasRefs {
		t.Error("expected HasRefs=true")
	}
}

func TestCollectStatsWithAllOf(t *testing.T) {
	s := &Schema{
		AllOf: []*Schema{{Type: "object"}},
	}
	stats := CollectStats(s)
	if !stats.HasAllOf {
		t.Error("expected HasAllOf=true")
	}
}

func TestGenerateDiagnostics(t *testing.T) {
	old := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"name": {Type: "string"},
		},
	}
	new := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"name": {Type: "integer"},
		},
	}

	oldData := []byte(`{"$schema":"http://json-schema.org/draft-07/schema#","type":"object"}`)
	newData := []byte(`{"$schema":"http://json-schema.org/draft-07/schema#","type":"object"}`)

	diag := GenerateDiagnostics(old, new, oldData, newData)
	if diag.OldDraft != Draft7 {
		t.Errorf("expected old draft=draft-07, got %v", diag.OldDraft)
	}
	if diag.OldSchemaStats.PropertyCount != 1 {
		t.Error("expected 1 property in old")
	}
}

func TestGenerateDiagnosticsCrossDraft(t *testing.T) {
	old := &Schema{Type: "string"}
	new := &Schema{Type: "string"}

	oldData := []byte(`{"$schema":"http://json-schema.org/draft-04/schema#","type":"string"}`)
	newData := []byte(`{"$schema":"https://json-schema.org/draft/2020-12/schema","type":"string"}`)

	diag := GenerateDiagnostics(old, new, oldData, newData)
	if diag.CompatibilityWarning == "" {
		t.Error("expected compatibility warning for cross-draft")
	}
}

func TestGenerateDiagnosticsWithValidationIssues(t *testing.T) {
	old := &Schema{Type: "string", MinLength: intPtr(10), MaxLength: intPtr(5)}
	new := &Schema{Type: "string"}

	diag := GenerateDiagnostics(old, new, nil, nil)
	if len(diag.ValidationIssues) == 0 {
		t.Error("expected validation issues")
	}
}

func TestFormatDiagnostics(t *testing.T) {
	diag := &DiagnosticInfo{
		OldSchemaStats: SchemaStats{Type: "object", PropertyCount: 5, RequiredCount: 2, Depth: 2},
		NewSchemaStats: SchemaStats{Type: "object", PropertyCount: 6, RequiredCount: 3, Depth: 2},
		OldDraft:       Draft7,
		NewDraft:       Draft7,
	}

	output := FormatDiagnostics(diag)
	if !strings.Contains(output, "Diagnostics") {
		t.Error("expected Diagnostics header")
	}
	if !strings.Contains(output, "Properties: 5") {
		t.Error("expected old property count")
	}
	if !strings.Contains(output, "Draft 7") {
		t.Error("expected draft info")
	}
}

func TestFormatDiagnosticsWithWarning(t *testing.T) {
	diag := &DiagnosticInfo{
		CompatibilityWarning: "different drafts",
	}
	output := FormatDiagnostics(diag)
	if !strings.Contains(output, "Warning:") {
		t.Error("expected warning in output")
	}
}

func TestFormatDiagnosticsWithIssues(t *testing.T) {
	diag := &DiagnosticInfo{
		ValidationIssues: []string{"old: type: invalid", "new: minLength > maxLength"},
	}
	output := FormatDiagnostics(diag)
	if !strings.Contains(output, "Validation Issues") {
		t.Error("expected validation issues section")
	}
}

func TestHasEnum(t *testing.T) {
	if hasEnum(nil) {
		t.Error("expected false for nil")
	}
	if !hasEnum(&Schema{Enum: []interface{}{"a"}}) {
		t.Error("expected true for enum")
	}
	nested := &Schema{
		Properties: map[string]*Schema{
			"x": {Enum: []interface{}{"a"}},
		},
	}
	if !hasEnum(nested) {
		t.Error("expected true for nested enum")
	}
}

func TestHasRefs(t *testing.T) {
	if hasRefs(nil) {
		t.Error("expected false for nil")
	}
	if !hasRefs(&Schema{Ref: "#/definitions/X"}) {
		t.Error("expected true for ref")
	}
}
