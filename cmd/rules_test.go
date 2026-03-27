package cmd

import (
	"bytes"
	"testing"
)

func TestRulesCmdRuns(t *testing.T) {
	var buf bytes.Buffer
	rulesCmd.SetOut(&buf)
	rulesCmd.SetArgs([]string{})
	if err := rulesCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	// The rules command writes to os.Stdout via tabwriter, not the cobra output
	// so we just verify it doesn't error
}
