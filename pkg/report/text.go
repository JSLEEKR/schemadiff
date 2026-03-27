// Package report provides output formatters for schema diff results.
package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/JSLEEKR/schemadiff/pkg/schema"
)

// TextReporter outputs human-readable text reports.
type TextReporter struct {
	Writer io.Writer
	Color  bool
}

// NewTextReporter creates a TextReporter writing to w.
func NewTextReporter(w io.Writer, color bool) *TextReporter {
	return &TextReporter{Writer: w, Color: color}
}

// Report writes a text report for a single schema diff.
func (r *TextReporter) Report(name string, result *schema.DiffResult) error {
	if result.Summary.Total == 0 {
		fmt.Fprintf(r.Writer, "%s: No changes detected\n", name)
		return nil
	}

	fmt.Fprintf(r.Writer, "\n%s\n", r.header(name))
	fmt.Fprintf(r.Writer, "%s\n\n", strings.Repeat("=", len(name)+4))

	for _, change := range result.Changes {
		icon := r.severityIcon(change.Severity)
		label := r.severityLabel(change.Severity)
		fmt.Fprintf(r.Writer, "  %s %s  %s\n", icon, label, change.Path)
		fmt.Fprintf(r.Writer, "     %s\n", change.Message)
		if change.OldValue != nil && change.NewValue != nil {
			fmt.Fprintf(r.Writer, "     Old: %v → New: %v\n", change.OldValue, change.NewValue)
		}
		fmt.Fprintln(r.Writer)
	}

	// Summary
	fmt.Fprintf(r.Writer, "  Summary: %d breaking, %d warning, %d info (%d total)\n\n",
		result.Summary.Breaking, result.Summary.Warning,
		result.Summary.Info, result.Summary.Total)

	return nil
}

// ReportMultiple writes a text report for multiple schema diffs.
func (r *TextReporter) ReportMultiple(results map[string]*schema.DiffResult) error {
	totalBreaking := 0
	totalWarning := 0
	totalInfo := 0

	for name, result := range results {
		if err := r.Report(name, result); err != nil {
			return err
		}
		totalBreaking += result.Summary.Breaking
		totalWarning += result.Summary.Warning
		totalInfo += result.Summary.Info
	}

	total := totalBreaking + totalWarning + totalInfo
	fmt.Fprintf(r.Writer, "Overall: %d breaking, %d warning, %d info (%d total changes across %d schemas)\n",
		totalBreaking, totalWarning, totalInfo, total, len(results))

	return nil
}

func (r *TextReporter) header(name string) string {
	return fmt.Sprintf("## %s", name)
}

func (r *TextReporter) severityIcon(s schema.Severity) string {
	if !r.Color {
		switch s {
		case schema.SeverityBreaking:
			return "[X]"
		case schema.SeverityWarning:
			return "[!]"
		case schema.SeverityInfo:
			return "[i]"
		default:
			return "[?]"
		}
	}
	switch s {
	case schema.SeverityBreaking:
		return "\033[31m[X]\033[0m"
	case schema.SeverityWarning:
		return "\033[33m[!]\033[0m"
	case schema.SeverityInfo:
		return "\033[36m[i]\033[0m"
	default:
		return "[?]"
	}
}

func (r *TextReporter) severityLabel(s schema.Severity) string {
	label := s.String()
	if !r.Color {
		return fmt.Sprintf("%-8s", label)
	}
	switch s {
	case schema.SeverityBreaking:
		return fmt.Sprintf("\033[31m%-8s\033[0m", label)
	case schema.SeverityWarning:
		return fmt.Sprintf("\033[33m%-8s\033[0m", label)
	case schema.SeverityInfo:
		return fmt.Sprintf("\033[36m%-8s\033[0m", label)
	default:
		return fmt.Sprintf("%-8s", label)
	}
}
