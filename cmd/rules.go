package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/JSLEEKR/schemadiff/pkg/schema"
	"github.com/spf13/cobra"
)

var rulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "List all breaking change detection rules",
	Long:  "Display all rules used to detect breaking, warning, and informational changes in schemas.",
	Run:   runRules,
}

func runRules(_ *cobra.Command, _ []string) {
	allTypes := []schema.ChangeType{
		schema.ChangeRequiredAdded,
		schema.ChangeRequiredRemoved,
		schema.ChangePropertyRemoved,
		schema.ChangePropertyAdded,
		schema.ChangeTypeChanged,
		schema.ChangeEnumValueRemoved,
		schema.ChangeEnumValueAdded,
		schema.ChangeMinLengthIncreased,
		schema.ChangeMaxLengthDecreased,
		schema.ChangeMinimumIncreased,
		schema.ChangeMaximumDecreased,
		schema.ChangeMinItemsIncreased,
		schema.ChangeMaxItemsDecreased,
		schema.ChangePatternChanged,
		schema.ChangeFormatChanged,
		schema.ChangeAdditionalPropsFalse,
		schema.ChangeItemsChanged,
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "RULE ID\tSEVERITY\tDESCRIPTION")
	fmt.Fprintln(w, "-------\t--------\t-----------")

	for _, ct := range allTypes {
		severity := schema.DefaultSeverity(ct)
		desc := schema.RuleDescription[ct]
		fmt.Fprintf(w, "%s\t%s\t%s\n", ct, severity, desc)
	}

	w.Flush()
}
