package schema

import (
	"fmt"
	"strings"
)

// DiagnosticInfo provides detailed diagnostic information about a schema comparison.
type DiagnosticInfo struct {
	// OldSchemaStats contains stats about the old schema.
	OldSchemaStats SchemaStats `json:"old_schema_stats"`
	// NewSchemaStats contains stats about the new schema.
	NewSchemaStats SchemaStats `json:"new_schema_stats"`
	// OldDraft is the detected draft version of the old schema.
	OldDraft Draft `json:"old_draft,omitempty"`
	// NewDraft is the detected draft version of the new schema.
	NewDraft Draft `json:"new_draft,omitempty"`
	// CompatibilityWarning is a message if drafts differ.
	CompatibilityWarning string `json:"compatibility_warning,omitempty"`
	// ValidationIssues lists issues found during validation.
	ValidationIssues []string `json:"validation_issues,omitempty"`
}

// SchemaStats provides aggregate statistics about a schema.
type SchemaStats struct {
	// PropertyCount is the total number of properties.
	PropertyCount int `json:"property_count"`
	// RequiredCount is the number of required fields.
	RequiredCount int `json:"required_count"`
	// Depth is the maximum nesting depth.
	Depth int `json:"depth"`
	// HasEnum indicates if any enum is defined.
	HasEnum bool `json:"has_enum"`
	// HasRefs indicates if any $ref is used.
	HasRefs bool `json:"has_refs"`
	// HasAllOf indicates if allOf is used.
	HasAllOf bool `json:"has_all_of"`
	// Type is the root schema type.
	Type string `json:"type"`
}

// CollectStats gathers statistics about a schema.
func CollectStats(s *Schema) SchemaStats {
	if s == nil {
		return SchemaStats{}
	}

	stats := SchemaStats{
		PropertyCount: propertyCount(s),
		RequiredCount: len(s.Required),
		Depth:         schemaDepth(s, 0),
		HasEnum:       hasEnum(s),
		HasRefs:       hasRefs(s),
		HasAllOf:      len(s.AllOf) > 0,
		Type:          normalizeType(s.Type),
	}

	return stats
}

// GenerateDiagnostics creates diagnostic info for a comparison.
func GenerateDiagnostics(old, new *Schema, oldData, newData []byte) *DiagnosticInfo {
	diag := &DiagnosticInfo{
		OldSchemaStats: CollectStats(old),
		NewSchemaStats: CollectStats(new),
	}

	if oldData != nil {
		diag.OldDraft = DetectDraft(oldData)
	}
	if newData != nil {
		diag.NewDraft = DetectDraft(newData)
	}

	if diag.OldDraft != "" && diag.NewDraft != "" {
		_, warning := IsDraftCompatible(diag.OldDraft, diag.NewDraft)
		diag.CompatibilityWarning = warning
	}

	// Collect validation issues
	if old != nil {
		for _, err := range ValidateSchema(old) {
			diag.ValidationIssues = append(diag.ValidationIssues, fmt.Sprintf("old: %s", err.Error()))
		}
	}
	if new != nil {
		for _, err := range ValidateSchema(new) {
			diag.ValidationIssues = append(diag.ValidationIssues, fmt.Sprintf("new: %s", err.Error()))
		}
	}

	return diag
}

// FormatDiagnostics returns a human-readable diagnostic report.
func FormatDiagnostics(diag *DiagnosticInfo) string {
	var sb strings.Builder

	sb.WriteString("=== Diagnostics ===\n\n")
	sb.WriteString("Old Schema:\n")
	sb.WriteString(formatStats(diag.OldSchemaStats))
	if diag.OldDraft != "" && diag.OldDraft != DraftUnknown {
		sb.WriteString(fmt.Sprintf("  Draft: %s\n", DraftInfo(diag.OldDraft)))
	}
	sb.WriteString("\nNew Schema:\n")
	sb.WriteString(formatStats(diag.NewSchemaStats))
	if diag.NewDraft != "" && diag.NewDraft != DraftUnknown {
		sb.WriteString(fmt.Sprintf("  Draft: %s\n", DraftInfo(diag.NewDraft)))
	}

	if diag.CompatibilityWarning != "" {
		sb.WriteString(fmt.Sprintf("\nWarning: %s\n", diag.CompatibilityWarning))
	}

	if len(diag.ValidationIssues) > 0 {
		sb.WriteString("\nValidation Issues:\n")
		for _, issue := range diag.ValidationIssues {
			sb.WriteString(fmt.Sprintf("  - %s\n", issue))
		}
	}

	return sb.String()
}

func formatStats(stats SchemaStats) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("  Type: %s\n", stats.Type))
	sb.WriteString(fmt.Sprintf("  Properties: %d\n", stats.PropertyCount))
	sb.WriteString(fmt.Sprintf("  Required: %d\n", stats.RequiredCount))
	sb.WriteString(fmt.Sprintf("  Depth: %d\n", stats.Depth))
	if stats.HasEnum {
		sb.WriteString("  Has enums: yes\n")
	}
	if stats.HasRefs {
		sb.WriteString("  Has $ref: yes\n")
	}
	if stats.HasAllOf {
		sb.WriteString("  Has allOf: yes\n")
	}
	return sb.String()
}

func hasEnum(s *Schema) bool {
	if s == nil {
		return false
	}
	if len(s.Enum) > 0 {
		return true
	}
	for _, p := range s.Properties {
		if hasEnum(p) {
			return true
		}
	}
	if s.Items != nil {
		return hasEnum(s.Items)
	}
	return false
}

func hasRefs(s *Schema) bool {
	if s == nil {
		return false
	}
	if s.Ref != "" {
		return true
	}
	for _, p := range s.Properties {
		if hasRefs(p) {
			return true
		}
	}
	if s.Items != nil {
		return hasRefs(s.Items)
	}
	return false
}
