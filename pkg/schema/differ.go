package schema

import (
	"fmt"
	"strings"
)

// Diff compares two schemas and returns all detected changes.
func Diff(old, new *Schema) *DiffResult {
	// Resolve refs before diffing
	old = ResolveRefs(old)
	new = ResolveRefs(new)

	// Merge allOf compositions
	old = MergeAllOf(old)
	new = MergeAllOf(new)

	var changes []Change
	changes = diffSchemas(old, new, "", changes)
	return NewDiffResult(changes)
}

func diffSchemas(old, new *Schema, path string, changes []Change) []Change {
	if old == nil && new == nil {
		return changes
	}
	if old == nil {
		// Entire schema is new — not breaking
		return changes
	}
	if new == nil {
		// Entire schema removed — breaking
		changes = append(changes, Change{
			Path:     pathOrRoot(path),
			Type:     ChangePropertyRemoved,
			Severity: SeverityBreaking,
			Message:  "Schema removed",
		})
		return changes
	}

	// Compare type
	changes = diffType(old, new, path, changes)

	// Compare required fields
	changes = diffRequired(old, new, path, changes)

	// Compare properties
	changes = diffProperties(old, new, path, changes)

	// Compare enum values
	changes = diffEnum(old, new, path, changes)

	// Compare constraints
	changes = diffConstraints(old, new, path, changes)

	// Compare pattern
	changes = diffPattern(old, new, path, changes)

	// Compare format
	changes = diffFormat(old, new, path, changes)

	// Compare additionalProperties
	changes = diffAdditionalProperties(old, new, path, changes)

	// Compare items (array schema)
	changes = diffItems(old, new, path, changes)

	return changes
}

func diffType(old, new *Schema, path string, changes []Change) []Change {
	oldType := normalizeType(old.Type)
	newType := normalizeType(new.Type)

	if oldType == "" || newType == "" {
		return changes
	}

	if oldType != newType {
		changes = append(changes, Change{
			Path:     joinPath(path, "type"),
			Type:     ChangeTypeChanged,
			Severity: SeverityBreaking,
			Message:  fmt.Sprintf("Type changed from %q to %q", oldType, newType),
			OldValue: oldType,
			NewValue: newType,
		})
	}
	return changes
}

func diffRequired(old, new *Schema, path string, changes []Change) []Change {
	oldSet := toStringSet(old.Required)
	newSet := toStringSet(new.Required)

	// Fields added to required
	for field := range newSet {
		if !oldSet[field] {
			changes = append(changes, Change{
				Path:     joinPath(path, "required"),
				Type:     ChangeRequiredAdded,
				Severity: SeverityBreaking,
				Message:  fmt.Sprintf("Field %q added to required", field),
				NewValue: field,
			})
		}
	}

	// Fields removed from required
	for field := range oldSet {
		if !newSet[field] {
			changes = append(changes, Change{
				Path:     joinPath(path, "required"),
				Type:     ChangeRequiredRemoved,
				Severity: SeverityInfo,
				Message:  fmt.Sprintf("Field %q removed from required", field),
				OldValue: field,
			})
		}
	}

	return changes
}

func diffProperties(old, new *Schema, path string, changes []Change) []Change {
	// Properties removed
	for name := range old.Properties {
		if _, exists := new.Properties[name]; !exists {
			changes = append(changes, Change{
				Path:     joinPath(path, "properties", name),
				Type:     ChangePropertyRemoved,
				Severity: SeverityBreaking,
				Message:  fmt.Sprintf("Property %q removed", name),
				OldValue: name,
			})
		}
	}

	// Properties added
	for name := range new.Properties {
		if _, exists := old.Properties[name]; !exists {
			changes = append(changes, Change{
				Path:     joinPath(path, "properties", name),
				Type:     ChangePropertyAdded,
				Severity: SeverityInfo,
				Message:  fmt.Sprintf("Property %q added", name),
				NewValue: name,
			})
		}
	}

	// Properties that exist in both — recurse
	for name, oldProp := range old.Properties {
		if newProp, exists := new.Properties[name]; exists {
			changes = diffSchemas(oldProp, newProp, joinPath(path, "properties", name), changes)
		}
	}

	return changes
}

func diffEnum(old, new *Schema, path string, changes []Change) []Change {
	if old.Enum == nil && new.Enum == nil {
		return changes
	}

	oldSet := toInterfaceSet(old.Enum)
	newSet := toInterfaceSet(new.Enum)

	// Enum values removed
	for val := range oldSet {
		if !newSet[val] {
			changes = append(changes, Change{
				Path:     joinPath(path, "enum"),
				Type:     ChangeEnumValueRemoved,
				Severity: SeverityBreaking,
				Message:  fmt.Sprintf("Enum value %q removed", val),
				OldValue: val,
			})
		}
	}

	// Enum values added
	for val := range newSet {
		if !oldSet[val] {
			changes = append(changes, Change{
				Path:     joinPath(path, "enum"),
				Type:     ChangeEnumValueAdded,
				Severity: SeverityInfo,
				Message:  fmt.Sprintf("Enum value %q added", val),
				NewValue: val,
			})
		}
	}

	return changes
}

func diffConstraints(old, new *Schema, path string, changes []Change) []Change {
	// minLength increased = stricter = breaking
	if old.MinLength != nil && new.MinLength != nil && *new.MinLength > *old.MinLength {
		changes = append(changes, Change{
			Path:     joinPath(path, "minLength"),
			Type:     ChangeMinLengthIncreased,
			Severity: SeverityBreaking,
			Message:  fmt.Sprintf("minLength increased from %d to %d", *old.MinLength, *new.MinLength),
			OldValue: *old.MinLength,
			NewValue: *new.MinLength,
		})
	} else if old.MinLength == nil && new.MinLength != nil {
		changes = append(changes, Change{
			Path:     joinPath(path, "minLength"),
			Type:     ChangeMinLengthIncreased,
			Severity: SeverityBreaking,
			Message:  fmt.Sprintf("minLength constraint added: %d", *new.MinLength),
			NewValue: *new.MinLength,
		})
	}

	// maxLength decreased = stricter = breaking
	if old.MaxLength != nil && new.MaxLength != nil && *new.MaxLength < *old.MaxLength {
		changes = append(changes, Change{
			Path:     joinPath(path, "maxLength"),
			Type:     ChangeMaxLengthDecreased,
			Severity: SeverityBreaking,
			Message:  fmt.Sprintf("maxLength decreased from %d to %d", *old.MaxLength, *new.MaxLength),
			OldValue: *old.MaxLength,
			NewValue: *new.MaxLength,
		})
	} else if old.MaxLength == nil && new.MaxLength != nil {
		changes = append(changes, Change{
			Path:     joinPath(path, "maxLength"),
			Type:     ChangeMaxLengthDecreased,
			Severity: SeverityBreaking,
			Message:  fmt.Sprintf("maxLength constraint added: %d", *new.MaxLength),
			NewValue: *new.MaxLength,
		})
	}

	// minimum increased = stricter = breaking
	if old.Minimum != nil && new.Minimum != nil && *new.Minimum > *old.Minimum {
		changes = append(changes, Change{
			Path:     joinPath(path, "minimum"),
			Type:     ChangeMinimumIncreased,
			Severity: SeverityBreaking,
			Message:  fmt.Sprintf("minimum increased from %v to %v", *old.Minimum, *new.Minimum),
			OldValue: *old.Minimum,
			NewValue: *new.Minimum,
		})
	} else if old.Minimum == nil && new.Minimum != nil {
		changes = append(changes, Change{
			Path:     joinPath(path, "minimum"),
			Type:     ChangeMinimumIncreased,
			Severity: SeverityBreaking,
			Message:  fmt.Sprintf("minimum constraint added: %v", *new.Minimum),
			NewValue: *new.Minimum,
		})
	}

	// maximum decreased = stricter = breaking
	if old.Maximum != nil && new.Maximum != nil && *new.Maximum < *old.Maximum {
		changes = append(changes, Change{
			Path:     joinPath(path, "maximum"),
			Type:     ChangeMaximumDecreased,
			Severity: SeverityBreaking,
			Message:  fmt.Sprintf("maximum decreased from %v to %v", *old.Maximum, *new.Maximum),
			OldValue: *old.Maximum,
			NewValue: *new.Maximum,
		})
	} else if old.Maximum == nil && new.Maximum != nil {
		changes = append(changes, Change{
			Path:     joinPath(path, "maximum"),
			Type:     ChangeMaximumDecreased,
			Severity: SeverityBreaking,
			Message:  fmt.Sprintf("maximum constraint added: %v", *new.Maximum),
			NewValue: *new.Maximum,
		})
	}

	// minItems increased = breaking
	if old.MinItems != nil && new.MinItems != nil && *new.MinItems > *old.MinItems {
		changes = append(changes, Change{
			Path:     joinPath(path, "minItems"),
			Type:     ChangeMinItemsIncreased,
			Severity: SeverityBreaking,
			Message:  fmt.Sprintf("minItems increased from %d to %d", *old.MinItems, *new.MinItems),
			OldValue: *old.MinItems,
			NewValue: *new.MinItems,
		})
	} else if old.MinItems == nil && new.MinItems != nil {
		changes = append(changes, Change{
			Path:     joinPath(path, "minItems"),
			Type:     ChangeMinItemsIncreased,
			Severity: SeverityBreaking,
			Message:  fmt.Sprintf("minItems constraint added: %d", *new.MinItems),
			NewValue: *new.MinItems,
		})
	}

	// maxItems decreased = breaking
	if old.MaxItems != nil && new.MaxItems != nil && *new.MaxItems < *old.MaxItems {
		changes = append(changes, Change{
			Path:     joinPath(path, "maxItems"),
			Type:     ChangeMaxItemsDecreased,
			Severity: SeverityBreaking,
			Message:  fmt.Sprintf("maxItems decreased from %d to %d", *old.MaxItems, *new.MaxItems),
			OldValue: *old.MaxItems,
			NewValue: *new.MaxItems,
		})
	} else if old.MaxItems == nil && new.MaxItems != nil {
		changes = append(changes, Change{
			Path:     joinPath(path, "maxItems"),
			Type:     ChangeMaxItemsDecreased,
			Severity: SeverityBreaking,
			Message:  fmt.Sprintf("maxItems constraint added: %d", *new.MaxItems),
			NewValue: *new.MaxItems,
		})
	}

	return changes
}

func diffPattern(old, new *Schema, path string, changes []Change) []Change {
	if old.Pattern != new.Pattern {
		if old.Pattern != "" && new.Pattern != "" {
			changes = append(changes, Change{
				Path:     joinPath(path, "pattern"),
				Type:     ChangePatternChanged,
				Severity: SeverityBreaking,
				Message:  fmt.Sprintf("Pattern changed from %q to %q", old.Pattern, new.Pattern),
				OldValue: old.Pattern,
				NewValue: new.Pattern,
			})
		} else if old.Pattern == "" && new.Pattern != "" {
			changes = append(changes, Change{
				Path:     joinPath(path, "pattern"),
				Type:     ChangePatternChanged,
				Severity: SeverityBreaking,
				Message:  fmt.Sprintf("Pattern constraint added: %q", new.Pattern),
				NewValue: new.Pattern,
			})
		}
	}
	return changes
}

func diffFormat(old, new *Schema, path string, changes []Change) []Change {
	if old.Format != new.Format && old.Format != "" && new.Format != "" {
		changes = append(changes, Change{
			Path:     joinPath(path, "format"),
			Type:     ChangeFormatChanged,
			Severity: SeverityWarning,
			Message:  fmt.Sprintf("Format changed from %q to %q", old.Format, new.Format),
			OldValue: old.Format,
			NewValue: new.Format,
		})
	}
	return changes
}

func diffAdditionalProperties(old, new *Schema, path string, changes []Change) []Change {
	oldAllow := old.AdditionalProperties == nil || *old.AdditionalProperties
	newAllow := new.AdditionalProperties == nil || *new.AdditionalProperties

	if oldAllow && !newAllow {
		changes = append(changes, Change{
			Path:     joinPath(path, "additionalProperties"),
			Type:     ChangeAdditionalPropsFalse,
			Severity: SeverityBreaking,
			Message:  "additionalProperties changed from allowed to restricted",
			OldValue: true,
			NewValue: false,
		})
	}
	return changes
}

func diffItems(old, new *Schema, path string, changes []Change) []Change {
	if old.Items != nil && new.Items != nil {
		changes = diffSchemas(old.Items, new.Items, joinPath(path, "items"), changes)
	} else if old.Items != nil && new.Items == nil {
		changes = append(changes, Change{
			Path:     joinPath(path, "items"),
			Type:     ChangeItemsChanged,
			Severity: SeverityBreaking,
			Message:  "Array items schema removed",
		})
	}
	return changes
}

// Helper functions

func joinPath(parts ...string) string {
	var nonEmpty []string
	for _, p := range parts {
		if p != "" {
			nonEmpty = append(nonEmpty, p)
		}
	}
	return strings.Join(nonEmpty, ".")
}

func pathOrRoot(path string) string {
	if path == "" {
		return "(root)"
	}
	return path
}

func normalizeType(t interface{}) string {
	if t == nil {
		return ""
	}
	switch v := t.(type) {
	case string:
		return v
	case []interface{}:
		if len(v) == 1 {
			if s, ok := v[0].(string); ok {
				return s
			}
		}
		strs := make([]string, len(v))
		for i, item := range v {
			if s, ok := item.(string); ok {
				strs[i] = s
			}
		}
		return strings.Join(strs, ",")
	default:
		return fmt.Sprintf("%v", v)
	}
}

func toStringSet(items []string) map[string]bool {
	set := make(map[string]bool, len(items))
	for _, item := range items {
		set[item] = true
	}
	return set
}

func toInterfaceSet(items []interface{}) map[string]bool {
	set := make(map[string]bool, len(items))
	for _, item := range items {
		set[fmt.Sprintf("%v", item)] = true
	}
	return set
}
