package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/JSLEEKR/schemadiff/pkg/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a .schemadiff.json config file",
	Long: `Create a .schemadiff.json configuration file in the current directory
with sensible defaults. The config file controls default behavior for
the check command.`,
	RunE: runInit,
}

var initForce bool

func init() {
	initCmd.Flags().BoolVar(&initForce, "force", false, "Overwrite existing config file")
}

func runInit(cmd *cobra.Command, args []string) error {
	path := ".schemadiff.json"

	// Check if file already exists
	if _, err := os.Stat(path); err == nil && !initForce {
		return fmt.Errorf("config file %s already exists (use --force to overwrite)", path)
	}

	cfg := config.DefaultConfig()
	cfg.Paths = []config.PathConfig{
		{
			Name:   "Example",
			Old:    "schemas/old.json",
			New:    "schemas/new.json",
			Format: "auto",
		},
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("generating config: %w", err)
	}

	if err := os.WriteFile(path, append(data, '\n'), 0644); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Created %s\n", path)
	fmt.Fprintln(os.Stdout, "Edit the file to configure your schema paths and settings.")
	return nil
}
