package report

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/JSLEEKR/schemadiff/pkg/schema"
)

func TestSARIFReporterSingle(t *testing.T) {
	var buf bytes.Buffer
	r := NewSARIFReporter(&buf, "old.json")
	result := schema.NewDiffResult([]schema.Change{
		{
			Path:     "type",
			Type:     schema.ChangeTypeChanged,
			Severity: schema.SeverityBreaking,
			Message:  "Type changed",
		},
	})
	if err := r.Report(result); err != nil {
		t.Fatal(err)
	}

	var report SARIFReport
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatalf("invalid SARIF JSON: %v", err)
	}
	if report.Version != "2.1.0" {
		t.Errorf("expected SARIF version=2.1.0, got %s", report.Version)
	}
	if len(report.Runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(report.Runs))
	}
	if report.Runs[0].Tool.Driver.Name != "schemadiff" {
		t.Errorf("expected tool name=schemadiff, got %s", report.Runs[0].Tool.Driver.Name)
	}
	if len(report.Runs[0].Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(report.Runs[0].Results))
	}
	if report.Runs[0].Results[0].Level != "error" {
		t.Errorf("expected level=error, got %s", report.Runs[0].Results[0].Level)
	}
}

func TestSARIFReporterMultiple(t *testing.T) {
	var buf bytes.Buffer
	r := NewSARIFReporter(&buf, "schema.json")
	results := map[string]*schema.DiffResult{
		"User": schema.NewDiffResult([]schema.Change{
			{Type: schema.ChangeTypeChanged, Severity: schema.SeverityBreaking, Message: "b1", Path: "type"},
			{Type: schema.ChangeRequiredAdded, Severity: schema.SeverityBreaking, Message: "b2", Path: "required"},
		}),
		"Post": schema.NewDiffResult([]schema.Change{
			{Type: schema.ChangePropertyAdded, Severity: schema.SeverityInfo, Message: "i1", Path: "props"},
		}),
	}
	if err := r.ReportMultiple(results); err != nil {
		t.Fatal(err)
	}

	var report SARIFReport
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatalf("invalid SARIF JSON: %v", err)
	}
	if len(report.Runs[0].Results) != 3 {
		t.Errorf("expected 3 results, got %d", len(report.Runs[0].Results))
	}
}

func TestSARIFReporterNoFile(t *testing.T) {
	var buf bytes.Buffer
	r := NewSARIFReporter(&buf, "")
	result := schema.NewDiffResult([]schema.Change{
		{Type: schema.ChangeTypeChanged, Severity: schema.SeverityBreaking, Message: "test"},
	})
	if err := r.Report(result); err != nil {
		t.Fatal(err)
	}

	var report SARIFReport
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if len(report.Runs[0].Results[0].Locations) != 0 {
		t.Error("expected no locations when no file path")
	}
}

func TestSARIFReporterWithLocation(t *testing.T) {
	var buf bytes.Buffer
	r := NewSARIFReporter(&buf, "api.yaml")
	result := schema.NewDiffResult([]schema.Change{
		{
			Path:     "properties.name",
			Type:     schema.ChangePropertyRemoved,
			Severity: schema.SeverityBreaking,
			Message:  "Property removed",
		},
	})
	if err := r.Report(result); err != nil {
		t.Fatal(err)
	}

	var report SARIFReport
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	loc := report.Runs[0].Results[0].Locations
	if len(loc) != 1 {
		t.Fatal("expected 1 location")
	}
	if loc[0].PhysicalLocation.ArtifactLocation.URI != "api.yaml" {
		t.Error("expected artifact URI=api.yaml")
	}
}

func TestSARIFReporterEmpty(t *testing.T) {
	var buf bytes.Buffer
	r := NewSARIFReporter(&buf, "")
	result := schema.NewDiffResult(nil)
	if err := r.Report(result); err != nil {
		t.Fatal(err)
	}
	var report SARIFReport
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if len(report.Runs[0].Results) != 0 {
		t.Error("expected 0 results")
	}
}

func TestSeverityToLevel(t *testing.T) {
	tests := []struct {
		sev  schema.Severity
		want string
	}{
		{schema.SeverityBreaking, "error"},
		{schema.SeverityWarning, "warning"},
		{schema.SeverityInfo, "note"},
		{schema.Severity(99), "none"},
	}
	for _, tt := range tests {
		if got := severityToLevel(tt.sev); got != tt.want {
			t.Errorf("severityToLevel(%d) = %q, want %q", tt.sev, got, tt.want)
		}
	}
}

func TestSARIFSchemaURL(t *testing.T) {
	var buf bytes.Buffer
	r := NewSARIFReporter(&buf, "")
	if err := r.Report(schema.NewDiffResult(nil)); err != nil {
		t.Fatal(err)
	}
	var report SARIFReport
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.Schema == "" {
		t.Error("expected schema URL")
	}
}

func TestSARIFRulesOnlyIncludeUsed(t *testing.T) {
	var buf bytes.Buffer
	r := NewSARIFReporter(&buf, "")
	result := schema.NewDiffResult([]schema.Change{
		{Type: schema.ChangeTypeChanged, Severity: schema.SeverityBreaking, Message: "test"},
	})
	if err := r.Report(result); err != nil {
		t.Fatal(err)
	}
	var report SARIFReport
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	rules := report.Runs[0].Tool.Driver.Rules
	if len(rules) != 1 {
		t.Errorf("expected 1 used rule, got %d", len(rules))
	}
	if rules[0].ID != string(schema.ChangeTypeChanged) {
		t.Errorf("expected rule ID=%s, got %s", schema.ChangeTypeChanged, rules[0].ID)
	}
}
