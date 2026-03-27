package schema

import "testing"

func TestFilterBySeverity(t *testing.T) {
	changes := []Change{
		{Severity: SeverityBreaking},
		{Severity: SeverityWarning},
		{Severity: SeverityInfo},
		{Severity: SeverityBreaking},
	}

	breaking := FilterBySeverity(changes, SeverityBreaking)
	if len(breaking) != 2 {
		t.Errorf("expected 2 breaking, got %d", len(breaking))
	}

	warningPlus := FilterBySeverity(changes, SeverityWarning)
	if len(warningPlus) != 3 {
		t.Errorf("expected 3 warning+, got %d", len(warningPlus))
	}

	all := FilterBySeverity(changes, SeverityInfo)
	if len(all) != 4 {
		t.Errorf("expected 4 info+, got %d", len(all))
	}
}

func TestFilterBySeverityEmpty(t *testing.T) {
	result := FilterBySeverity(nil, SeverityBreaking)
	if len(result) != 0 {
		t.Error("expected empty")
	}
}

func TestFilterByType(t *testing.T) {
	changes := []Change{
		{Type: ChangeTypeChanged},
		{Type: ChangePropertyRemoved},
		{Type: ChangePropertyAdded},
		{Type: ChangeTypeChanged},
	}

	typeChanged := FilterByType(changes, ChangeTypeChanged)
	if len(typeChanged) != 2 {
		t.Errorf("expected 2 type-changed, got %d", len(typeChanged))
	}

	multi := FilterByType(changes, ChangeTypeChanged, ChangePropertyRemoved)
	if len(multi) != 3 {
		t.Errorf("expected 3 multi-type, got %d", len(multi))
	}
}

func TestFilterByTypeEmpty(t *testing.T) {
	result := FilterByType(nil, ChangeTypeChanged)
	if len(result) != 0 {
		t.Error("expected empty")
	}
}

func TestFilterByPath(t *testing.T) {
	changes := []Change{
		{Path: "properties.name.type"},
		{Path: "properties.age.type"},
		{Path: "required"},
		{Path: "properties.name.minLength"},
	}

	nameChanges := FilterByPath(changes, "properties.name")
	if len(nameChanges) != 2 {
		t.Errorf("expected 2 name changes, got %d", len(nameChanges))
	}

	propChanges := FilterByPath(changes, "properties")
	if len(propChanges) != 3 {
		t.Errorf("expected 3 property changes, got %d", len(propChanges))
	}
}

func TestFilterByPathEmpty(t *testing.T) {
	result := FilterByPath(nil, "test")
	if len(result) != 0 {
		t.Error("expected empty")
	}
}

func TestGroupBySeverity(t *testing.T) {
	changes := []Change{
		{Severity: SeverityBreaking},
		{Severity: SeverityInfo},
		{Severity: SeverityBreaking},
	}
	groups := GroupBySeverity(changes)
	if len(groups[SeverityBreaking]) != 2 {
		t.Error("expected 2 breaking")
	}
	if len(groups[SeverityInfo]) != 1 {
		t.Error("expected 1 info")
	}
}

func TestGroupByPath(t *testing.T) {
	changes := []Change{
		{Path: "properties.name.type"},
		{Path: "properties.age.type"},
		{Path: "required"},
	}
	groups := GroupByPath(changes)
	if len(groups["properties"]) != 2 {
		t.Errorf("expected 2 properties, got %d", len(groups["properties"]))
	}
	if len(groups["required"]) != 1 {
		t.Error("expected 1 required")
	}
}

func TestGroupByPathNoSeparator(t *testing.T) {
	changes := []Change{{Path: "type"}}
	groups := GroupByPath(changes)
	if len(groups["type"]) != 1 {
		t.Error("expected 1 type change")
	}
}
