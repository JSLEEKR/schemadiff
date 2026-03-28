package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/JSLEEKR/schemadiff/pkg/openapi"
	"github.com/JSLEEKR/schemadiff/pkg/report"
	"github.com/JSLEEKR/schemadiff/pkg/schema"
	"github.com/spf13/cobra"
)

// rejectTraversal checks that a cleaned path does not contain ".." segments.
func rejectTraversal(path string) error {
	clean := filepath.Clean(path)
	if strings.Contains(clean, "..") {
		return fmt.Errorf("path contains traversal sequence: %s", path)
	}
	return nil
}

var (
	outputFormat string
	inputFormat  string
	failOn       string
	noColor      bool
)

var checkCmd = &cobra.Command{
	Use:   "check <old> <new>",
	Short: "Compare two schema files and report changes",
	Long: `Compare two schema files (JSON Schema or OpenAPI 3.x) and report
breaking changes, warnings, and informational changes.

Exit codes:
  0 = no breaking changes (or --fail-on not triggered)
  1 = breaking changes detected (when --fail-on breaking)
  2 = invalid input or error`,
	Args: cobra.ExactArgs(2),
	RunE: runCheck,
}

func init() {
	checkCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format: text, json, sarif")
	checkCmd.Flags().StringVarP(&inputFormat, "format", "f", "auto", "Input format: auto, jsonschema, openapi")
	checkCmd.Flags().StringVar(&failOn, "fail-on", "breaking", "Fail condition: breaking, warning, any, none")
	checkCmd.Flags().BoolVar(&noColor, "no-color", false, "Disable colored output")
}

func runCheck(cmd *cobra.Command, args []string) error {
	oldPath := args[0]
	newPath := args[1]

	// Reject path traversal attempts
	if err := rejectTraversal(oldPath); err != nil {
		return err
	}
	if err := rejectTraversal(newPath); err != nil {
		return err
	}

	// Detect format
	format := inputFormat
	if format == "auto" {
		format = detectFormat(oldPath)
	}

	switch format {
	case "openapi":
		return runOpenAPICheck(oldPath, newPath)
	default:
		return runJSONSchemaCheck(oldPath, newPath)
	}
}

func runJSONSchemaCheck(oldPath, newPath string) error {
	oldSchema, err := schema.SafeParseFile(oldPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing old schema: %v\n", err)
		os.Exit(2)
	}

	newSchema, err := schema.SafeParseFile(newPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing new schema: %v\n", err)
		os.Exit(2)
	}

	result := schema.Diff(oldSchema, newSchema)
	return outputResult(result, oldPath)
}

func runOpenAPICheck(oldPath, newPath string) error {
	oldSpec, err := openapi.ParseFile(oldPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing old OpenAPI spec: %v\n", err)
		os.Exit(2)
	}

	newSpec, err := openapi.ParseFile(newPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing new OpenAPI spec: %v\n", err)
		os.Exit(2)
	}

	results := openapi.DiffSpecs(oldSpec, newSpec)
	return outputMultipleResults(results, oldPath)
}

func outputResult(result *schema.DiffResult, filePath string) error {
	switch outputFormat {
	case "json":
		r := report.NewJSONReporter(os.Stdout, true)
		if err := r.Report(result); err != nil {
			return err
		}
	case "sarif":
		r := report.NewSARIFReporter(os.Stdout, filePath)
		if err := r.Report(result); err != nil {
			return err
		}
	case "markdown", "md":
		r := report.NewMarkdownReporter(os.Stdout)
		if err := r.Report("Schema", result); err != nil {
			return err
		}
	default:
		r := report.NewTextReporter(os.Stdout, !noColor)
		if err := r.Report("Schema", result); err != nil {
			return err
		}
	}

	exitOnFailCondition(result)
	return nil
}

func outputMultipleResults(results map[string]*schema.DiffResult, filePath string) error {
	// Merge all results for exit code determination
	var allChanges []schema.Change
	for _, r := range results {
		allChanges = append(allChanges, r.Changes...)
	}
	merged := schema.NewDiffResult(allChanges)

	switch outputFormat {
	case "json":
		r := report.NewJSONReporter(os.Stdout, true)
		if err := r.ReportMultiple(results); err != nil {
			return err
		}
	case "sarif":
		r := report.NewSARIFReporter(os.Stdout, filePath)
		if err := r.ReportMultiple(results); err != nil {
			return err
		}
	case "markdown", "md":
		r := report.NewMarkdownReporter(os.Stdout)
		if err := r.ReportMultiple(results); err != nil {
			return err
		}
	default:
		r := report.NewTextReporter(os.Stdout, !noColor)
		if err := r.ReportMultiple(results); err != nil {
			return err
		}
	}

	exitOnFailCondition(merged)
	return nil
}

func exitOnFailCondition(result *schema.DiffResult) {
	switch failOn {
	case "breaking":
		if result.HasBreaking {
			os.Exit(1)
		}
	case "warning":
		if result.HasBreaking || result.Summary.Warning > 0 {
			os.Exit(1)
		}
	case "any":
		if result.Summary.Total > 0 {
			os.Exit(1)
		}
	case "none":
		// Never fail
	}
}

func detectFormat(path string) string {
	// Try to detect OpenAPI by parsing
	spec, err := openapi.ParseFile(path)
	if err == nil && spec.OpenAPI != "" {
		return "openapi"
	}
	return "jsonschema"
}
