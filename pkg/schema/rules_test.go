package schema

import "testing"

func TestDefaultSeverity(t *testing.T) {
	breakingTypes := []ChangeType{
		ChangeRequiredAdded, ChangePropertyRemoved, ChangeTypeChanged,
		ChangeEnumValueRemoved, ChangeMinLengthIncreased, ChangeMaxLengthDecreased,
		ChangeMinimumIncreased, ChangeMaximumDecreased, ChangeMinItemsIncreased,
		ChangeMaxItemsDecreased, ChangePatternChanged, ChangeAdditionalPropsFalse,
		ChangeItemsChanged,
	}
	for _, ct := range breakingTypes {
		if DefaultSeverity(ct) != SeverityBreaking {
			t.Errorf("expected BREAKING for %q", ct)
		}
	}

	if DefaultSeverity(ChangeFormatChanged) != SeverityWarning {
		t.Error("expected WARNING for format-changed")
	}

	infoTypes := []ChangeType{
		ChangeRequiredRemoved, ChangePropertyAdded, ChangeEnumValueAdded,
	}
	for _, ct := range infoTypes {
		if DefaultSeverity(ct) != SeverityInfo {
			t.Errorf("expected INFO for %q", ct)
		}
	}

	if DefaultSeverity(ChangeType("unknown")) != SeverityWarning {
		t.Error("expected WARNING for unknown type")
	}
}

func TestRuleDescriptions(t *testing.T) {
	allTypes := []ChangeType{
		ChangeRequiredAdded, ChangeRequiredRemoved, ChangePropertyRemoved,
		ChangePropertyAdded, ChangeTypeChanged, ChangeEnumValueRemoved,
		ChangeEnumValueAdded, ChangeMinLengthIncreased, ChangeMaxLengthDecreased,
		ChangeMinimumIncreased, ChangeMaximumDecreased, ChangeMinItemsIncreased,
		ChangeMaxItemsDecreased, ChangePatternChanged, ChangeFormatChanged,
		ChangeAdditionalPropsFalse, ChangeItemsChanged,
	}
	for _, ct := range allTypes {
		desc, ok := RuleDescription[ct]
		if !ok {
			t.Errorf("missing description for %q", ct)
		}
		if desc == "" {
			t.Errorf("empty description for %q", ct)
		}
	}
}
