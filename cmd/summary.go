package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/JSLEEKR/schemadiff/pkg/openapi"
	"github.com/JSLEEKR/schemadiff/pkg/schema"
	"github.com/spf13/cobra"
)

var summaryCmd = &cobra.Command{
	Use:   "summary <old> <new>",
	Short: "Print a concise one-line compatibility summary",
	Long: `Print a concise summary showing whether two schemas are compatible.
Useful for scripts and quick checks.

Output format: COMPATIBLE|BREAKING|WARNING <counts>`,
	Args: cobra.ExactArgs(2),
	RunE: runSummary,
}

var summaryJSON bool

func init() {
	summaryCmd.Flags().StringVarP(&inputFormat, "format", "f", "auto", "Input format: auto, jsonschema, openapi")
	summaryCmd.Flags().BoolVar(&summaryJSON, "json", false, "Output as JSON")
}

type summaryOutput struct {
	Status   string `json:"status"`
	Breaking int    `json:"breaking"`
	Warning  int    `json:"warning"`
	Info     int    `json:"info"`
	Total    int    `json:"total"`
}

func runSummary(cmd *cobra.Command, args []string) error {
	oldPath := args[0]
	newPath := args[1]

	format := inputFormat
	if format == "auto" {
		format = detectFormat(oldPath)
	}

	var result *schema.DiffResult
	switch format {
	case "openapi":
		oldSpec, err := openapi.ParseFile(oldPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(2)
		}
		newSpec, err := openapi.ParseFile(newPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(2)
		}
		results := openapi.DiffSpecs(oldSpec, newSpec)
		var allChanges []schema.Change
		for _, r := range results {
			allChanges = append(allChanges, r.Changes...)
		}
		result = schema.NewDiffResult(allChanges)
	default:
		oldSchema, err := schema.ParseFile(oldPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(2)
		}
		newSchema, err := schema.ParseFile(newPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(2)
		}
		result = schema.Diff(oldSchema, newSchema)
	}

	out := summaryOutput{
		Breaking: result.Summary.Breaking,
		Warning:  result.Summary.Warning,
		Info:     result.Summary.Info,
		Total:    result.Summary.Total,
	}

	if result.HasBreaking {
		out.Status = "BREAKING"
	} else if result.Summary.Warning > 0 {
		out.Status = "WARNING"
	} else {
		out.Status = "COMPATIBLE"
	}

	if summaryJSON {
		data, _ := json.Marshal(out)
		fmt.Fprintln(os.Stdout, string(data))
	} else {
		fmt.Fprintf(os.Stdout, "%s: %d breaking, %d warning, %d info (%d total)\n",
			out.Status, out.Breaking, out.Warning, out.Info, out.Total)
	}

	if result.HasBreaking {
		os.Exit(1)
	}
	return nil
}
