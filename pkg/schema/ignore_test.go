package schema

import "testing"

func TestIgnoreByRule(t *testing.T) {
	ic := NewIgnoreConfig([]string{"format-changed", "pattern-changed"})

	if !ic.ShouldIgnore(Change{Type: ChangeFormatChanged}) {
		t.Error("expected format-changed to be ignored")
	}
	if !ic.ShouldIgnore(Change{Type: ChangePatternChanged}) {
		t.Error("expected pattern-changed to be ignored")
	}
	if ic.ShouldIgnore(Change{Type: ChangeTypeChanged}) {
		t.Error("expected type-changed not to be ignored")
	}
}

func TestIgnoreByPath(t *testing.T) {
	ic := &IgnoreConfig{
		Paths: []string{"properties.internal", "properties.deprecated"},
	}

	if !ic.ShouldIgnore(Change{Path: "properties.internal.name"}) {
		t.Error("expected internal path to be ignored")
	}
	if ic.ShouldIgnore(Change{Path: "properties.name"}) {
		t.Error("expected name path not to be ignored")
	}
}

func TestIgnoreBySeverity(t *testing.T) {
	ic := &IgnoreConfig{
		Severities: []Severity{SeverityInfo},
	}

	if !ic.ShouldIgnore(Change{Severity: SeverityInfo}) {
		t.Error("expected info to be ignored")
	}
	if ic.ShouldIgnore(Change{Severity: SeverityBreaking}) {
		t.Error("expected breaking not to be ignored")
	}
}

func TestIgnoreNilConfig(t *testing.T) {
	var ic *IgnoreConfig
	if ic.ShouldIgnore(Change{Type: ChangeTypeChanged}) {
		t.Error("expected nil config to not ignore anything")
	}
}

func TestApplyIgnoreRules(t *testing.T) {
	changes := []Change{
		{Type: ChangeTypeChanged, Path: "type"},
		{Type: ChangeFormatChanged, Path: "format"},
		{Type: ChangePropertyAdded, Path: "props"},
	}

	ic := NewIgnoreConfig([]string{"format-changed"})
	filtered := ApplyIgnoreRules(changes, ic)
	if len(filtered) != 2 {
		t.Errorf("expected 2 changes after filtering, got %d", len(filtered))
	}
}

func TestApplyIgnoreRulesNil(t *testing.T) {
	changes := []Change{{Type: ChangeTypeChanged}}
	filtered := ApplyIgnoreRules(changes, nil)
	if len(filtered) != 1 {
		t.Error("expected no filtering with nil config")
	}
}

func TestApplyIgnoreRulesEmpty(t *testing.T) {
	filtered := ApplyIgnoreRules(nil, NewIgnoreConfig([]string{"any"}))
	if len(filtered) != 0 {
		t.Error("expected empty for nil changes")
	}
}

func TestIgnoreMultipleCriteria(t *testing.T) {
	ic := &IgnoreConfig{
		Rules: []string{"format-changed"},
		Paths: []string{"properties.internal"},
	}

	// Matches rule
	if !ic.ShouldIgnore(Change{Type: ChangeFormatChanged, Path: "x"}) {
		t.Error("expected rule match to ignore")
	}
	// Matches path
	if !ic.ShouldIgnore(Change{Type: ChangeTypeChanged, Path: "properties.internal.x"}) {
		t.Error("expected path match to ignore")
	}
	// Matches neither
	if ic.ShouldIgnore(Change{Type: ChangeTypeChanged, Path: "properties.name"}) {
		t.Error("expected no match to not ignore")
	}
}
