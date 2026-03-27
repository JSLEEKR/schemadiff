package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTestFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestDetectFormatJSONSchema(t *testing.T) {
	dir := t.TempDir()
	path := writeTestFile(t, dir, "schema.json", `{"type":"object"}`)
	format := detectFormat(path)
	if format != "jsonschema" {
		t.Errorf("expected jsonschema, got %s", format)
	}
}

func TestDetectFormatOpenAPI(t *testing.T) {
	dir := t.TempDir()
	path := writeTestFile(t, dir, "api.json", `{"openapi":"3.0.0","info":{"title":"Test","version":"1.0.0"},"paths":{}}`)
	format := detectFormat(path)
	if format != "openapi" {
		t.Errorf("expected openapi, got %s", format)
	}
}

func TestCheckCmdRequiresArgs(t *testing.T) {
	// cobra.ExactArgs(2) requires exactly 2 arguments
	cmd := checkCmd
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error for missing args")
	}
	err = cmd.Args(cmd, []string{"one"})
	if err == nil {
		t.Error("expected error for one arg")
	}
	err = cmd.Args(cmd, []string{"one", "two"})
	if err != nil {
		t.Errorf("expected no error for two args, got %v", err)
	}
}

func TestCheckCmdFlags(t *testing.T) {
	flags := checkCmd.Flags()

	outputFlag := flags.Lookup("output")
	if outputFlag == nil {
		t.Error("expected --output flag")
	}
	if outputFlag.DefValue != "text" {
		t.Errorf("expected default=text, got %s", outputFlag.DefValue)
	}

	formatFlag := flags.Lookup("format")
	if formatFlag == nil {
		t.Error("expected --format flag")
	}
	if formatFlag.DefValue != "auto" {
		t.Errorf("expected default=auto, got %s", formatFlag.DefValue)
	}

	failOnFlag := flags.Lookup("fail-on")
	if failOnFlag == nil {
		t.Error("expected --fail-on flag")
	}
	if failOnFlag.DefValue != "breaking" {
		t.Errorf("expected default=breaking, got %s", failOnFlag.DefValue)
	}
}
