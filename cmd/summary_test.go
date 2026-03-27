package cmd

import "testing"

func TestSummaryCmdArgs(t *testing.T) {
	err := summaryCmd.Args(summaryCmd, []string{})
	if err == nil {
		t.Error("expected error for no args")
	}
	err = summaryCmd.Args(summaryCmd, []string{"a", "b"})
	if err != nil {
		t.Errorf("expected no error for 2 args, got %v", err)
	}
}

func TestSummaryCmdFlags(t *testing.T) {
	f := summaryCmd.Flags()
	if f.Lookup("format") == nil {
		t.Error("expected --format flag")
	}
	if f.Lookup("json") == nil {
		t.Error("expected --json flag")
	}
}

func TestRootHasSummary(t *testing.T) {
	cmds := rootCmd.Commands()
	found := false
	for _, c := range cmds {
		if c.Name() == "summary" {
			found = true
		}
	}
	if !found {
		t.Error("expected summary subcommand")
	}
}
