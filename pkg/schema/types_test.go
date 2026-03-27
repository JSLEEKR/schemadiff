package schema

import (
	"testing"
)

func TestSeverityString(t *testing.T) {
	tests := []struct {
		sev  Severity
		want string
	}{
		{SeverityInfo, "INFO"},
		{SeverityWarning, "WARNING"},
		{SeverityBreaking, "BREAKING"},
		{Severity(99), "UNKNOWN"},
	}
	for _, tt := range tests {
		if got := tt.sev.String(); got != tt.want {
			t.Errorf("Severity(%d).String() = %q, want %q", tt.sev, got, tt.want)
		}
	}
}

func TestChangeString(t *testing.T) {
	c := Change{
		Path:     "properties.name.type",
		Type:     ChangeTypeChanged,
		Severity: SeverityBreaking,
		Message:  "Type changed from \"string\" to \"number\"",
	}
	got := c.String()
	if got != `[BREAKING] properties.name.type: Type changed from "string" to "number"` {
		t.Errorf("Change.String() = %q", got)
	}
}

func TestNewDiffResult(t *testing.T) {
	changes := []Change{
		{Severity: SeverityBreaking},
		{Severity: SeverityBreaking},
		{Severity: SeverityWarning},
		{Severity: SeverityInfo},
		{Severity: SeverityInfo},
		{Severity: SeverityInfo},
	}
	result := NewDiffResult(changes)
	if !result.HasBreaking {
		t.Error("expected HasBreaking=true")
	}
	if result.Summary.Breaking != 2 {
		t.Errorf("expected Breaking=2, got %d", result.Summary.Breaking)
	}
	if result.Summary.Warning != 1 {
		t.Errorf("expected Warning=1, got %d", result.Summary.Warning)
	}
	if result.Summary.Info != 3 {
		t.Errorf("expected Info=3, got %d", result.Summary.Info)
	}
	if result.Summary.Total != 6 {
		t.Errorf("expected Total=6, got %d", result.Summary.Total)
	}
}

func TestNewDiffResultNoBreaking(t *testing.T) {
	changes := []Change{
		{Severity: SeverityInfo},
	}
	result := NewDiffResult(changes)
	if result.HasBreaking {
		t.Error("expected HasBreaking=false")
	}
}

func TestNewDiffResultEmpty(t *testing.T) {
	result := NewDiffResult(nil)
	if result.HasBreaking {
		t.Error("expected HasBreaking=false for empty")
	}
	if result.Summary.Total != 0 {
		t.Error("expected Total=0 for empty")
	}
}
