package report

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/JSLEEKR/schemadiff/pkg/schema"
)

// JSONReport represents the complete JSON output.
type JSONReport struct {
	Version  string                       `json:"version"`
	Schemas  map[string]*schema.DiffResult `json:"schemas,omitempty"`
	Result   *schema.DiffResult           `json:"result,omitempty"`
	Summary  JSONSummary                  `json:"summary"`
}

// JSONSummary provides aggregate stats.
type JSONSummary struct {
	TotalSchemas  int  `json:"total_schemas"`
	TotalChanges  int  `json:"total_changes"`
	Breaking      int  `json:"breaking"`
	Warning       int  `json:"warning"`
	Info          int  `json:"info"`
	HasBreaking   bool `json:"has_breaking"`
}

// JSONReporter outputs JSON-formatted reports.
type JSONReporter struct {
	Writer io.Writer
	Pretty bool
}

// NewJSONReporter creates a JSONReporter.
func NewJSONReporter(w io.Writer, pretty bool) *JSONReporter {
	return &JSONReporter{Writer: w, Pretty: pretty}
}

// Report writes a JSON report for a single schema diff.
func (r *JSONReporter) Report(result *schema.DiffResult) error {
	report := JSONReport{
		Version: "1.0.0",
		Result:  result,
		Summary: JSONSummary{
			TotalSchemas: 1,
			TotalChanges: result.Summary.Total,
			Breaking:     result.Summary.Breaking,
			Warning:      result.Summary.Warning,
			Info:         result.Summary.Info,
			HasBreaking:  result.HasBreaking,
		},
	}
	return r.write(report)
}

// ReportMultiple writes a JSON report for multiple schema diffs.
func (r *JSONReporter) ReportMultiple(results map[string]*schema.DiffResult) error {
	summary := JSONSummary{
		TotalSchemas: len(results),
	}
	for _, result := range results {
		summary.TotalChanges += result.Summary.Total
		summary.Breaking += result.Summary.Breaking
		summary.Warning += result.Summary.Warning
		summary.Info += result.Summary.Info
		if result.HasBreaking {
			summary.HasBreaking = true
		}
	}

	report := JSONReport{
		Version: "1.0.0",
		Schemas: results,
		Summary: summary,
	}
	return r.write(report)
}

func (r *JSONReporter) write(report JSONReport) error {
	var data []byte
	var err error
	if r.Pretty {
		data, err = json.MarshalIndent(report, "", "  ")
	} else {
		data, err = json.Marshal(report)
	}
	if err != nil {
		return fmt.Errorf("marshaling JSON report: %w", err)
	}

	_, err = r.Writer.Write(data)
	if err != nil {
		return fmt.Errorf("writing JSON report: %w", err)
	}
	_, err = r.Writer.Write([]byte("\n"))
	return err
}
