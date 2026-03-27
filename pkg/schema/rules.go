package schema

// RuleDescription provides human-readable explanations for each change type.
var RuleDescription = map[ChangeType]string{
	ChangeRequiredAdded:        "Adding a required field breaks existing consumers that don't include it.",
	ChangeRequiredRemoved:      "Removing a required field is backward-compatible; existing payloads still work.",
	ChangePropertyRemoved:      "Removing a property breaks consumers that depend on it.",
	ChangePropertyAdded:        "Adding an optional property is backward-compatible.",
	ChangeTypeChanged:          "Changing a field's type breaks consumers expecting the original type.",
	ChangeEnumValueRemoved:     "Removing an enum value breaks consumers that send or expect that value.",
	ChangeEnumValueAdded:       "Adding an enum value is backward-compatible for producers.",
	ChangeMinLengthIncreased:   "Increasing minLength rejects previously valid shorter strings.",
	ChangeMaxLengthDecreased:   "Decreasing maxLength rejects previously valid longer strings.",
	ChangeMinimumIncreased:     "Increasing minimum rejects previously valid lower numbers.",
	ChangeMaximumDecreased:     "Decreasing maximum rejects previously valid higher numbers.",
	ChangeMinItemsIncreased:    "Increasing minItems rejects previously valid shorter arrays.",
	ChangeMaxItemsDecreased:    "Decreasing maxItems rejects previously valid longer arrays.",
	ChangePatternChanged:       "Changing the pattern rejects previously valid strings.",
	ChangeFormatChanged:        "Changing the format may break consumers with format validation.",
	ChangeAdditionalPropsFalse: "Restricting additional properties rejects previously valid payloads with extra fields.",
	ChangeItemsChanged:         "Changing the array items schema may break consumers.",
}

// DefaultSeverity returns the default severity for a change type.
func DefaultSeverity(ct ChangeType) Severity {
	switch ct {
	case ChangeRequiredAdded,
		ChangePropertyRemoved,
		ChangeTypeChanged,
		ChangeEnumValueRemoved,
		ChangeMinLengthIncreased,
		ChangeMaxLengthDecreased,
		ChangeMinimumIncreased,
		ChangeMaximumDecreased,
		ChangeMinItemsIncreased,
		ChangeMaxItemsDecreased,
		ChangePatternChanged,
		ChangeAdditionalPropsFalse,
		ChangeItemsChanged:
		return SeverityBreaking

	case ChangeFormatChanged:
		return SeverityWarning

	case ChangeRequiredRemoved,
		ChangePropertyAdded,
		ChangeEnumValueAdded:
		return SeverityInfo

	default:
		return SeverityWarning
	}
}
