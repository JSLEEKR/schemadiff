package report

import (
	"fmt"
	"io"

	"github.com/JSLEEKR/schemadiff/pkg/schema"
)

// MarkdownReporter outputs GitHub-flavored Markdown reports.
type MarkdownReporter struct {
	Writer io.Writer
}

// NewMarkdownReporter creates a MarkdownReporter.
func NewMarkdownReporter(w io.Writer) *MarkdownReporter {
	return &MarkdownReporter{Writer: w}
}

// Report writes a Markdown report for a single schema diff.
func (r *MarkdownReporter) Report(name string, result *schema.DiffResult) error {
	if result.Summary.Total == 0 {
		fmt.Fprintf(r.Writer, "### %s\n\nNo changes detected.\n\n", name)
		return nil
	}

	fmt.Fprintf(r.Writer, "### %s\n\n", name)

	// Breaking changes
	if result.Summary.Breaking > 0 {
		fmt.Fprintf(r.Writer, "**Breaking Changes (%d)**\n\n", result.Summary.Breaking)
		fmt.Fprintln(r.Writer, "| Path | Description |")
		fmt.Fprintln(r.Writer, "|------|-------------|")
		for _, c := range result.Changes {
			if c.Severity == schema.SeverityBreaking {
				fmt.Fprintf(r.Writer, "| `%s` | %s |\n", c.Path, c.Message)
			}
		}
		fmt.Fprintln(r.Writer)
	}

	// Warnings
	if result.Summary.Warning > 0 {
		fmt.Fprintf(r.Writer, "**Warnings (%d)**\n\n", result.Summary.Warning)
		fmt.Fprintln(r.Writer, "| Path | Description |")
		fmt.Fprintln(r.Writer, "|------|-------------|")
		for _, c := range result.Changes {
			if c.Severity == schema.SeverityWarning {
				fmt.Fprintf(r.Writer, "| `%s` | %s |\n", c.Path, c.Message)
			}
		}
		fmt.Fprintln(r.Writer)
	}

	// Info
	if result.Summary.Info > 0 {
		fmt.Fprintf(r.Writer, "**Non-breaking (%d)**\n\n", result.Summary.Info)
		fmt.Fprintln(r.Writer, "| Path | Description |")
		fmt.Fprintln(r.Writer, "|------|-------------|")
		for _, c := range result.Changes {
			if c.Severity == schema.SeverityInfo {
				fmt.Fprintf(r.Writer, "| `%s` | %s |\n", c.Path, c.Message)
			}
		}
		fmt.Fprintln(r.Writer)
	}

	return nil
}

// ReportMultiple writes a Markdown report for multiple schema diffs.
func (r *MarkdownReporter) ReportMultiple(results map[string]*schema.DiffResult) error {
	fmt.Fprint(r.Writer, "## Schema Compatibility Report\n\n")

	totalBreaking := 0
	totalWarning := 0
	totalInfo := 0

	for _, result := range results {
		totalBreaking += result.Summary.Breaking
		totalWarning += result.Summary.Warning
		totalInfo += result.Summary.Info
	}

	// Summary badge
	if totalBreaking > 0 {
		fmt.Fprintf(r.Writer, "> **BREAKING** - %d breaking, %d warning, %d info\n\n", totalBreaking, totalWarning, totalInfo)
	} else if totalWarning > 0 {
		fmt.Fprintf(r.Writer, "> **WARNING** - %d warning, %d info\n\n", totalWarning, totalInfo)
	} else {
		fmt.Fprintf(r.Writer, "> **COMPATIBLE** - %d informational changes\n\n", totalInfo)
	}

	for name, result := range results {
		if err := r.Report(name, result); err != nil {
			return err
		}
	}

	return nil
}
