package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/JSLEEKR/schemadiff/pkg/schema"
)

func TestMarkdownReporterNoChanges(t *testing.T) {
	var buf bytes.Buffer
	r := NewMarkdownReporter(&buf)
	if err := r.Report("Test", schema.NewDiffResult(nil)); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No changes") {
		t.Error("expected 'No changes' message")
	}
}

func TestMarkdownReporterBreaking(t *testing.T) {
	var buf bytes.Buffer
	r := NewMarkdownReporter(&buf)
	result := schema.NewDiffResult([]schema.Change{
		{Path: "type", Severity: schema.SeverityBreaking, Message: "Type changed"},
	})
	if err := r.Report("Schema", result); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "Breaking Changes") {
		t.Error("expected 'Breaking Changes' header")
	}
	if !strings.Contains(out, "`type`") {
		t.Error("expected code-formatted path")
	}
}

func TestMarkdownReporterWarning(t *testing.T) {
	var buf bytes.Buffer
	r := NewMarkdownReporter(&buf)
	result := schema.NewDiffResult([]schema.Change{
		{Path: "format", Severity: schema.SeverityWarning, Message: "Format changed"},
	})
	if err := r.Report("Schema", result); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "Warnings") {
		t.Error("expected 'Warnings' header")
	}
}

func TestMarkdownReporterInfo(t *testing.T) {
	var buf bytes.Buffer
	r := NewMarkdownReporter(&buf)
	result := schema.NewDiffResult([]schema.Change{
		{Path: "props.x", Severity: schema.SeverityInfo, Message: "Added"},
	})
	if err := r.Report("Schema", result); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "Non-breaking") {
		t.Error("expected 'Non-breaking' header")
	}
}

func TestMarkdownReporterMultiple(t *testing.T) {
	var buf bytes.Buffer
	r := NewMarkdownReporter(&buf)
	results := map[string]*schema.DiffResult{
		"User": schema.NewDiffResult([]schema.Change{
			{Severity: schema.SeverityBreaking, Message: "break"},
		}),
	}
	if err := r.ReportMultiple(results); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "Schema Compatibility Report") {
		t.Error("expected report title")
	}
	if !strings.Contains(out, "BREAKING") {
		t.Error("expected BREAKING status")
	}
}

func TestMarkdownReporterMultipleCompatible(t *testing.T) {
	var buf bytes.Buffer
	r := NewMarkdownReporter(&buf)
	results := map[string]*schema.DiffResult{
		"User": schema.NewDiffResult([]schema.Change{
			{Severity: schema.SeverityInfo, Message: "info"},
		}),
	}
	if err := r.ReportMultiple(results); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "COMPATIBLE") {
		t.Error("expected COMPATIBLE status")
	}
}

func TestMarkdownReporterTable(t *testing.T) {
	var buf bytes.Buffer
	r := NewMarkdownReporter(&buf)
	result := schema.NewDiffResult([]schema.Change{
		{Path: "a", Severity: schema.SeverityBreaking, Message: "m1"},
		{Path: "b", Severity: schema.SeverityBreaking, Message: "m2"},
	})
	if err := r.Report("Test", result); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "| `a` |") {
		t.Error("expected table row for path a")
	}
	if !strings.Contains(out, "| `b` |") {
		t.Error("expected table row for path b")
	}
}
