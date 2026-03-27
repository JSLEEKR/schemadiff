package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitCmdCreatesFile(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	initForce = false
	if err := runInit(initCmd, nil); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(dir, ".schemadiff.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected config file to be created")
	}
}

func TestInitCmdExistingFile(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	// Create existing file
	if err := os.WriteFile(".schemadiff.json", []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}

	initForce = false
	err := runInit(initCmd, nil)
	if err == nil {
		t.Error("expected error for existing file")
	}
}

func TestInitCmdForce(t *testing.T) {
	dir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	// Create existing file
	if err := os.WriteFile(".schemadiff.json", []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}

	initForce = true
	err := runInit(initCmd, nil)
	if err != nil {
		t.Errorf("expected no error with --force: %v", err)
	}
}

func TestRootHasInit(t *testing.T) {
	cmds := rootCmd.Commands()
	found := false
	for _, c := range cmds {
		if c.Name() == "init" {
			found = true
		}
	}
	if !found {
		t.Error("expected init subcommand")
	}
}
