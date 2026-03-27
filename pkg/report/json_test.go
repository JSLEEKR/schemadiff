package report

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/JSLEEKR/schemadiff/pkg/schema"
)

func TestJSONReporterSingle(t *testing.T) {
	var buf bytes.Buffer
	r := NewJSONReporter(&buf, true)
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

	var report JSONReport
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if report.Version != "1.0.0" {
		t.Errorf("expected version=1.0.0, got %s", report.Version)
	}
	if !report.Summary.HasBreaking {
		t.Error("expected has_breaking=true")
	}
	if report.Summary.Breaking != 1 {
		t.Errorf("expected breaking=1, got %d", report.Summary.Breaking)
	}
}

func TestJSONReporterMultiple(t *testing.T) {
	var buf bytes.Buffer
	r := NewJSONReporter(&buf, false)
	results := map[string]*schema.DiffResult{
		"User": schema.NewDiffResult([]schema.Change{
			{Severity: schema.SeverityBreaking, Message: "b1"},
		}),
		"Post": schema.NewDiffResult([]schema.Change{
			{Severity: schema.SeverityInfo, Message: "i1"},
		}),
	}
	if err := r.ReportMultiple(results); err != nil {
		t.Fatal(err)
	}

	var report JSONReport
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if report.Summary.TotalSchemas != 2 {
		t.Errorf("expected 2 schemas, got %d", report.Summary.TotalSchemas)
	}
	if report.Summary.TotalChanges != 2 {
		t.Errorf("expected 2 changes, got %d", report.Summary.TotalChanges)
	}
}

func TestJSONReporterPretty(t *testing.T) {
	var buf bytes.Buffer
	r := NewJSONReporter(&buf, true)
	result := schema.NewDiffResult([]schema.Change{
		{Severity: schema.SeverityInfo, Message: "test"},
	})
	if err := r.Report(result); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if out[0:2] != "{\n" {
		t.Error("expected pretty-printed JSON with newlines")
	}
}

func TestJSONReporterCompact(t *testing.T) {
	var buf bytes.Buffer
	r := NewJSONReporter(&buf, false)
	result := schema.NewDiffResult([]schema.Change{
		{Severity: schema.SeverityInfo, Message: "test"},
	})
	if err := r.Report(result); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	// Compact JSON should be a single line (before the trailing newline)
	lines := 0
	for _, c := range out {
		if c == '\n' {
			lines++
		}
	}
	if lines != 1 {
		t.Errorf("expected compact single-line JSON, got %d lines", lines)
	}
}

func TestJSONReporterEmpty(t *testing.T) {
	var buf bytes.Buffer
	r := NewJSONReporter(&buf, false)
	result := schema.NewDiffResult(nil)
	if err := r.Report(result); err != nil {
		t.Fatal(err)
	}
	var report JSONReport
	if err := json.Unmarshal(buf.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if report.Summary.TotalChanges != 0 {
		t.Error("expected 0 changes")
	}
}
