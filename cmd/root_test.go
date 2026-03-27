package cmd

import (
	"testing"
)

func TestRootCmdVersion(t *testing.T) {
	if Version == "" {
		t.Error("expected non-empty version")
	}
}

func TestRootCmdHasSubcommands(t *testing.T) {
	cmds := rootCmd.Commands()
	names := make(map[string]bool)
	for _, c := range cmds {
		names[c.Name()] = true
	}
	if !names["check"] {
		t.Error("expected check subcommand")
	}
	if !names["rules"] {
		t.Error("expected rules subcommand")
	}
}
