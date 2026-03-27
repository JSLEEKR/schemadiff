// Package cmd provides the CLI commands for schemadiff.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is set at build time.
var Version = "1.0.0"

var rootCmd = &cobra.Command{
	Use:   "schemadiff",
	Short: "Detect breaking schema changes before they ship",
	Long: `schemadiff detects backward-incompatible changes in JSON Schema,
OpenAPI 3.x, and Protobuf definitions. Built for CI pipelines.

Use it to catch breaking changes before they reach production.`,
	Version: Version,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(rulesCmd)
	rootCmd.AddCommand(summaryCmd)
	rootCmd.AddCommand(initCmd)
}
