package report

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/JSLEEKR/schemadiff/pkg/schema"
)

// SARIF 2.1.0 types for CI integration

// SARIFReport is the top-level SARIF 2.1.0 structure.
type SARIFReport struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []SARIFRun `json:"runs"`
}

// SARIFRun represents a single analysis run.
type SARIFRun struct {
	Tool    SARIFTool    `json:"tool"`
	Results []SARIFResult `json:"results"`
}

// SARIFTool describes the analysis tool.
type SARIFTool struct {
	Driver SARIFDriver `json:"driver"`
}

// SARIFDriver describes the tool driver.
type SARIFDriver struct {
	Name           string      `json:"name"`
	Version        string      `json:"version"`
	InformationURI string      `json:"informationUri"`
	Rules          []SARIFRule `json:"rules"`
}

// SARIFRule describes a rule that detected an issue.
type SARIFRule struct {
	ID               string          `json:"id"`
	ShortDescription SARIFMessage    `json:"shortDescription"`
	FullDescription  *SARIFMessage   `json:"fullDescription,omitempty"`
	DefaultConfig    *SARIFRuleConfig `json:"defaultConfiguration,omitempty"`
}

// SARIFRuleConfig is rule-level configuration.
type SARIFRuleConfig struct {
	Level string `json:"level"`
}

// SARIFResult is a single finding.
type SARIFResult struct {
	RuleID    string        `json:"ruleId"`
	Level     string        `json:"level"`
	Message   SARIFMessage  `json:"message"`
	Locations []SARIFLocation `json:"locations,omitempty"`
}

// SARIFMessage is a text message.
type SARIFMessage struct {
	Text string `json:"text"`
}

// SARIFLocation is a code location.
type SARIFLocation struct {
	PhysicalLocation *SARIFPhysicalLocation `json:"physicalLocation,omitempty"`
	LogicalLocations []SARIFLogicalLocation  `json:"logicalLocations,omitempty"`
}

// SARIFPhysicalLocation is a file location.
type SARIFPhysicalLocation struct {
	ArtifactLocation SARIFArtifactLocation `json:"artifactLocation"`
}

// SARIFArtifactLocation identifies a file.
type SARIFArtifactLocation struct {
	URI string `json:"uri"`
}

// SARIFLogicalLocation is a named location within a schema.
type SARIFLogicalLocation struct {
	Name string `json:"name"`
	Kind string `json:"kind,omitempty"`
}

// SARIFReporter outputs SARIF 2.1.0 formatted reports.
type SARIFReporter struct {
	Writer   io.Writer
	FilePath string
}

// NewSARIFReporter creates a SARIFReporter.
func NewSARIFReporter(w io.Writer, filePath string) *SARIFReporter {
	return &SARIFReporter{Writer: w, FilePath: filePath}
}

// Report writes a SARIF report for a single schema diff.
func (r *SARIFReporter) Report(result *schema.DiffResult) error {
	changes := make(map[string]*schema.DiffResult)
	changes["schema"] = result
	return r.ReportMultiple(changes)
}

// ReportMultiple writes a SARIF report for multiple schema diffs.
func (r *SARIFReporter) ReportMultiple(results map[string]*schema.DiffResult) error {
	rules := buildRules()
	ruleIndex := buildRuleIndex(rules)

	var sarifResults []SARIFResult
	for schemaName, result := range results {
		for _, change := range result.Changes {
			ruleID := string(change.Type)
			level := severityToLevel(change.Severity)

			sarifResult := SARIFResult{
				RuleID:  ruleID,
				Level:   level,
				Message: SARIFMessage{Text: change.Message},
			}

			if r.FilePath != "" {
				sarifResult.Locations = []SARIFLocation{
					{
						PhysicalLocation: &SARIFPhysicalLocation{
							ArtifactLocation: SARIFArtifactLocation{URI: r.FilePath},
						},
						LogicalLocations: []SARIFLogicalLocation{
							{Name: fmt.Sprintf("%s.%s", schemaName, change.Path), Kind: "jsonPointer"},
						},
					},
				}
			}

			sarifResults = append(sarifResults, sarifResult)
		}
	}

	// Use only rules that have results
	var usedRules []SARIFRule
	usedRuleIDs := make(map[string]bool)
	for _, r := range sarifResults {
		if !usedRuleIDs[r.RuleID] {
			usedRuleIDs[r.RuleID] = true
			if idx, ok := ruleIndex[r.RuleID]; ok {
				usedRules = append(usedRules, rules[idx])
			}
		}
	}

	report := SARIFReport{
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Version: "2.1.0",
		Runs: []SARIFRun{
			{
				Tool: SARIFTool{
					Driver: SARIFDriver{
						Name:           "schemadiff",
						Version:        "1.0.0",
						InformationURI: "https://github.com/JSLEEKR/schemadiff",
						Rules:          usedRules,
					},
				},
				Results: sarifResults,
			},
		},
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling SARIF report: %w", err)
	}

	_, err = r.Writer.Write(data)
	if err != nil {
		return fmt.Errorf("writing SARIF report: %w", err)
	}
	_, err = r.Writer.Write([]byte("\n"))
	return err
}

func buildRules() []SARIFRule {
	allTypes := []schema.ChangeType{
		schema.ChangeRequiredAdded, schema.ChangeRequiredRemoved,
		schema.ChangePropertyRemoved, schema.ChangePropertyAdded,
		schema.ChangeTypeChanged, schema.ChangeEnumValueRemoved,
		schema.ChangeEnumValueAdded, schema.ChangeMinLengthIncreased,
		schema.ChangeMaxLengthDecreased, schema.ChangeMinimumIncreased,
		schema.ChangeMaximumDecreased, schema.ChangeMinItemsIncreased,
		schema.ChangeMaxItemsDecreased, schema.ChangePatternChanged,
		schema.ChangeFormatChanged, schema.ChangeAdditionalPropsFalse,
		schema.ChangeItemsChanged,
	}

	var rules []SARIFRule
	for _, ct := range allTypes {
		desc := schema.RuleDescription[ct]
		level := severityToLevel(schema.DefaultSeverity(ct))
		rules = append(rules, SARIFRule{
			ID:               string(ct),
			ShortDescription: SARIFMessage{Text: desc},
			DefaultConfig:    &SARIFRuleConfig{Level: level},
		})
	}
	return rules
}

func buildRuleIndex(rules []SARIFRule) map[string]int {
	index := make(map[string]int, len(rules))
	for i, rule := range rules {
		index[rule.ID] = i
	}
	return index
}

func severityToLevel(s schema.Severity) string {
	switch s {
	case schema.SeverityBreaking:
		return "error"
	case schema.SeverityWarning:
		return "warning"
	case schema.SeverityInfo:
		return "note"
	default:
		return "none"
	}
}
