package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/JSLEEKR/schemadiff/pkg/schema"
)

func TestTextReporterNoChanges(t *testing.T) {
	var buf bytes.Buffer
	r := NewTextReporter(&buf, false)
	result := schema.NewDiffResult(nil)
	if err := r.Report("Test", result); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No changes") {
		t.Error("expected 'No changes' message")
	}
}

func TestTextReporterBreaking(t *testing.T) {
	var buf bytes.Buffer
	r := NewTextReporter(&buf, false)
	result := schema.NewDiffResult([]schema.Change{
		{
			Path:     "properties.name.type",
			Type:     schema.ChangeTypeChanged,
			Severity: schema.SeverityBreaking,
			Message:  "Type changed from string to number",
			OldValue: "string",
			NewValue: "number",
		},
	})
	if err := r.Report("UserSchema", result); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "BREAKING") {
		t.Error("expected BREAKING label")
	}
	if !strings.Contains(out, "properties.name.type") {
		t.Error("expected path")
	}
	if !strings.Contains(out, "Old: string") {
		t.Error("expected old value")
	}
}

func TestTextReporterColor(t *testing.T) {
	var buf bytes.Buffer
	r := NewTextReporter(&buf, true)
	result := schema.NewDiffResult([]schema.Change{
		{
			Path:     "type",
			Type:     schema.ChangeTypeChanged,
			Severity: schema.SeverityBreaking,
			Message:  "Type changed",
		},
	})
	if err := r.Report("Test", result); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "\033[31m") {
		t.Error("expected red ANSI code for breaking")
	}
}

func TestTextReporterWarning(t *testing.T) {
	var buf bytes.Buffer
	r := NewTextReporter(&buf, false)
	result := schema.NewDiffResult([]schema.Change{
		{
			Path:     "format",
			Type:     schema.ChangeFormatChanged,
			Severity: schema.SeverityWarning,
			Message:  "Format changed",
		},
	})
	if err := r.Report("Test", result); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "WARNING") {
		t.Error("expected WARNING label")
	}
}

func TestTextReporterInfo(t *testing.T) {
	var buf bytes.Buffer
	r := NewTextReporter(&buf, false)
	result := schema.NewDiffResult([]schema.Change{
		{
			Path:     "properties.phone",
			Type:     schema.ChangePropertyAdded,
			Severity: schema.SeverityInfo,
			Message:  "Property added",
		},
	})
	if err := r.Report("Test", result); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "INFO") {
		t.Error("expected INFO label")
	}
}

func TestTextReporterMultiple(t *testing.T) {
	var buf bytes.Buffer
	r := NewTextReporter(&buf, false)
	results := map[string]*schema.DiffResult{
		"User": schema.NewDiffResult([]schema.Change{
			{Severity: schema.SeverityBreaking, Message: "break1"},
		}),
		"Post": schema.NewDiffResult([]schema.Change{
			{Severity: schema.SeverityInfo, Message: "info1"},
		}),
	}
	if err := r.ReportMultiple(results); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "Overall:") {
		t.Error("expected Overall summary")
	}
	if !strings.Contains(out, "2 schemas") {
		t.Error("expected schema count")
	}
}

func TestTextReporterSummary(t *testing.T) {
	var buf bytes.Buffer
	r := NewTextReporter(&buf, false)
	result := schema.NewDiffResult([]schema.Change{
		{Severity: schema.SeverityBreaking, Message: "b1"},
		{Severity: schema.SeverityBreaking, Message: "b2"},
		{Severity: schema.SeverityWarning, Message: "w1"},
		{Severity: schema.SeverityInfo, Message: "i1"},
	})
	if err := r.Report("Test", result); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "2 breaking") {
		t.Error("expected '2 breaking' in summary")
	}
	if !strings.Contains(out, "4 total") {
		t.Error("expected '4 total' in summary")
	}
}
