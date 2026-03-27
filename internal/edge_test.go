package internal

import (
	"testing"

	"github.com/JSLEEKR/schemadiff/pkg/schema"
)

// Edge case tests for comprehensive coverage

func TestEdgeCaseEmptyProperties(t *testing.T) {
	old := &schema.Schema{Type: "object", Properties: map[string]*schema.Schema{}}
	new := &schema.Schema{Type: "object", Properties: map[string]*schema.Schema{}}
	result := schema.Diff(old, new)
	if result.Summary.Total != 0 {
		t.Error("expected no changes for empty properties")
	}
}

func TestEdgeCaseNilProperties(t *testing.T) {
	old := &schema.Schema{Type: "object"}
	new := &schema.Schema{Type: "object"}
	result := schema.Diff(old, new)
	if result.Summary.Total != 0 {
		t.Error("expected no changes for nil properties")
	}
}

func TestEdgeCaseEnumEmptyToNil(t *testing.T) {
	old := &schema.Schema{Type: "string", Enum: []interface{}{}}
	new := &schema.Schema{Type: "string"}
	result := schema.Diff(old, new)
	// Empty enum to nil is not a change
	if result.HasBreaking {
		t.Error("expected no breaking for empty enum to nil")
	}
}

func TestEdgeCaseTypeArraySingle(t *testing.T) {
	old := &schema.Schema{Type: []interface{}{"string"}}
	new := &schema.Schema{Type: "string"}
	result := schema.Diff(old, new)
	if result.HasBreaking {
		t.Error("expected no breaking for equivalent type representations")
	}
}

func TestEdgeCaseConstraintsSameValues(t *testing.T) {
	min := 5
	max := 100
	old := &schema.Schema{Type: "string", MinLength: &min, MaxLength: &max}
	new := &schema.Schema{Type: "string", MinLength: &min, MaxLength: &max}
	result := schema.Diff(old, new)
	if result.Summary.Total != 0 {
		t.Error("expected no changes for same constraints")
	}
}

func TestEdgeCaseConstraintsLoosened(t *testing.T) {
	oldMin, newMin := 10, 1
	oldMax, newMax := 50, 200
	old := &schema.Schema{Type: "string", MinLength: &oldMin, MaxLength: &oldMax}
	new := &schema.Schema{Type: "string", MinLength: &newMin, MaxLength: &newMax}
	result := schema.Diff(old, new)
	if result.HasBreaking {
		t.Error("expected no breaking for loosened constraints")
	}
}

func TestEdgeCaseMultipleBreaking(t *testing.T) {
	min1, min2 := 1, 10
	old := &schema.Schema{
		Type: "object",
		Properties: map[string]*schema.Schema{
			"a": {Type: "string"},
			"b": {Type: "integer"},
			"c": {Type: "string", MinLength: &min1},
		},
		Required: []string{"a"},
	}
	new := &schema.Schema{
		Type: "object",
		Properties: map[string]*schema.Schema{
			"a": {Type: "number"},        // type changed
			"c": {Type: "string", MinLength: &min2}, // constraint tightened
			// b removed
		},
		Required: []string{"a", "c"}, // c added to required
	}
	result := schema.Diff(old, new)
	if result.Summary.Breaking < 4 {
		t.Errorf("expected at least 4 breaking changes, got %d", result.Summary.Breaking)
	}
}

func TestEdgeCaseDeepNested(t *testing.T) {
	old := &schema.Schema{
		Type: "object",
		Properties: map[string]*schema.Schema{
			"level1": {
				Type: "object",
				Properties: map[string]*schema.Schema{
					"level2": {
						Type: "object",
						Properties: map[string]*schema.Schema{
							"level3": {
								Type: "object",
								Properties: map[string]*schema.Schema{
									"value": {Type: "string"},
								},
							},
						},
					},
				},
			},
		},
	}
	new := &schema.Schema{
		Type: "object",
		Properties: map[string]*schema.Schema{
			"level1": {
				Type: "object",
				Properties: map[string]*schema.Schema{
					"level2": {
						Type: "object",
						Properties: map[string]*schema.Schema{
							"level3": {
								Type: "object",
								Properties: map[string]*schema.Schema{
									"value": {Type: "integer"}, // deep change
								},
							},
						},
					},
				},
			},
		},
	}
	result := schema.Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking change at deep level")
	}
	// Verify path includes full nesting
	for _, c := range result.Changes {
		if c.Type == schema.ChangeTypeChanged {
			expected := "properties.level1.properties.level2.properties.level3.properties.value.type"
			if c.Path != expected {
				t.Errorf("expected path=%s, got %s", expected, c.Path)
			}
		}
	}
}

func TestEdgeCaseAdditionalPropsTrueToTrue(t *testing.T) {
	tr := true
	old := &schema.Schema{Type: "object", AdditionalProperties: &tr}
	new := &schema.Schema{Type: "object", AdditionalProperties: &tr}
	result := schema.Diff(old, new)
	if result.HasBreaking {
		t.Error("expected no breaking for true->true")
	}
}

func TestEdgeCaseAdditionalPropsFalseToTrue(t *testing.T) {
	f, tr := false, true
	old := &schema.Schema{Type: "object", AdditionalProperties: &f}
	new := &schema.Schema{Type: "object", AdditionalProperties: &tr}
	result := schema.Diff(old, new)
	// Loosening is not breaking
	if result.HasBreaking {
		t.Error("expected no breaking for false->true")
	}
}

func TestEdgeCaseItemsAddedNotBreaking(t *testing.T) {
	old := &schema.Schema{Type: "array"}
	new := &schema.Schema{Type: "array", Items: &schema.Schema{Type: "string"}}
	result := schema.Diff(old, new)
	// Adding items schema is not breaking (was unconstrained before)
	if result.HasBreaking {
		t.Error("expected no breaking for adding items schema")
	}
}

func TestEdgeCaseRefResolutionInDiff(t *testing.T) {
	old := &schema.Schema{
		Type: "object",
		Properties: map[string]*schema.Schema{
			"addr": {Ref: "#/definitions/Address"},
		},
		Definitions: map[string]*schema.Schema{
			"Address": {
				Type: "object",
				Properties: map[string]*schema.Schema{
					"city": {Type: "string"},
				},
			},
		},
	}
	new := &schema.Schema{
		Type: "object",
		Properties: map[string]*schema.Schema{
			"addr": {Ref: "#/definitions/Address"},
		},
		Definitions: map[string]*schema.Schema{
			"Address": {
				Type: "object",
				Properties: map[string]*schema.Schema{
					"city": {Type: "integer"}, // changed through ref
				},
			},
		},
	}
	result := schema.Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking through ref resolution")
	}
}

func TestEdgeCaseMergeAllOfThenDiff(t *testing.T) {
	old := &schema.Schema{
		AllOf: []*schema.Schema{
			{Type: "object", Properties: map[string]*schema.Schema{"name": {Type: "string"}}},
			{Required: []string{"name"}},
		},
	}
	new := &schema.Schema{
		AllOf: []*schema.Schema{
			{Type: "object", Properties: map[string]*schema.Schema{"name": {Type: "integer"}}},
			{Required: []string{"name"}},
		},
	}
	result := schema.Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking through allOf merge and diff")
	}
}

func TestEdgeCaseFilterAndReport(t *testing.T) {
	changes := []schema.Change{
		{Type: schema.ChangeTypeChanged, Severity: schema.SeverityBreaking, Path: "type", Message: "changed"},
		{Type: schema.ChangeFormatChanged, Severity: schema.SeverityWarning, Path: "format", Message: "changed"},
		{Type: schema.ChangePropertyAdded, Severity: schema.SeverityInfo, Path: "props.x", Message: "added"},
	}

	// Filter to breaking only
	breaking := schema.FilterBySeverity(changes, schema.SeverityBreaking)
	if len(breaking) != 1 {
		t.Errorf("expected 1 breaking, got %d", len(breaking))
	}

	// Group by severity
	groups := schema.GroupBySeverity(changes)
	if len(groups) != 3 {
		t.Errorf("expected 3 groups, got %d", len(groups))
	}
}

func TestEdgeCaseIgnoreRulesIntegration(t *testing.T) {
	changes := []schema.Change{
		{Type: schema.ChangeFormatChanged, Severity: schema.SeverityWarning, Path: "format"},
		{Type: schema.ChangeTypeChanged, Severity: schema.SeverityBreaking, Path: "type"},
	}

	ic := schema.NewIgnoreConfig([]string{"format-changed"})
	filtered := schema.ApplyIgnoreRules(changes, ic)
	if len(filtered) != 1 {
		t.Errorf("expected 1 after ignore, got %d", len(filtered))
	}
	if filtered[0].Type != schema.ChangeTypeChanged {
		t.Error("expected type-changed to remain")
	}
}

func TestEdgeCaseValidateAndDiff(t *testing.T) {
	s := &schema.Schema{
		Type: "object",
		Properties: map[string]*schema.Schema{
			"good": {Type: "string"},
			"bad":  {Type: "string", MinLength: intPtr(10), MaxLength: intPtr(5)},
		},
	}

	errs := schema.ValidateSchema(s)
	if len(errs) == 0 {
		t.Error("expected validation errors")
	}

	// Should still be able to diff despite validation issues
	s2 := &schema.Schema{Type: "object", Properties: map[string]*schema.Schema{
		"good": {Type: "integer"},
	}}
	result := schema.Diff(s, s2)
	if !result.HasBreaking {
		t.Error("expected breaking changes despite validation issues")
	}
}

func intPtr(i int) *int { return &i }
