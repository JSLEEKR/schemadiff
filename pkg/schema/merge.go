package schema

// MergeAllOf merges allOf schemas into a single flattened schema.
// This is needed because many schemas use allOf for composition.
func MergeAllOf(s *Schema) *Schema {
	if s == nil {
		return nil
	}

	result := *s

	// If there's an allOf, merge all sub-schemas
	if len(s.AllOf) > 0 {
		for _, sub := range s.AllOf {
			merged := MergeAllOf(sub)
			result = mergeInto(result, *merged)
		}
		result.AllOf = nil
	}

	// Recurse into properties
	if result.Properties != nil {
		newProps := make(map[string]*Schema, len(result.Properties))
		for k, v := range result.Properties {
			newProps[k] = MergeAllOf(v)
		}
		result.Properties = newProps
	}

	// Recurse into items
	if result.Items != nil {
		result.Items = MergeAllOf(result.Items)
	}

	return &result
}

func mergeInto(base, overlay Schema) Schema {
	// Merge type
	if overlay.Type != nil && base.Type == nil {
		base.Type = overlay.Type
	}

	// Merge properties
	if overlay.Properties != nil {
		if base.Properties == nil {
			base.Properties = make(map[string]*Schema)
		}
		for k, v := range overlay.Properties {
			base.Properties[k] = v
		}
	}

	// Merge required (union)
	if len(overlay.Required) > 0 {
		existing := toStringSet(base.Required)
		for _, r := range overlay.Required {
			if !existing[r] {
				base.Required = append(base.Required, r)
			}
		}
	}

	// Merge enum (intersection if both present, overlay if only overlay)
	if overlay.Enum != nil {
		if base.Enum == nil {
			base.Enum = overlay.Enum
		}
	}

	// Merge constraints (take stricter)
	if overlay.MinLength != nil {
		if base.MinLength == nil || *overlay.MinLength > *base.MinLength {
			base.MinLength = overlay.MinLength
		}
	}
	if overlay.MaxLength != nil {
		if base.MaxLength == nil || *overlay.MaxLength < *base.MaxLength {
			base.MaxLength = overlay.MaxLength
		}
	}
	if overlay.Minimum != nil {
		if base.Minimum == nil || *overlay.Minimum > *base.Minimum {
			base.Minimum = overlay.Minimum
		}
	}
	if overlay.Maximum != nil {
		if base.Maximum == nil || *overlay.Maximum < *base.Maximum {
			base.Maximum = overlay.Maximum
		}
	}
	if overlay.MinItems != nil {
		if base.MinItems == nil || *overlay.MinItems > *base.MinItems {
			base.MinItems = overlay.MinItems
		}
	}
	if overlay.MaxItems != nil {
		if base.MaxItems == nil || *overlay.MaxItems < *base.MaxItems {
			base.MaxItems = overlay.MaxItems
		}
	}

	// Merge pattern
	if overlay.Pattern != "" && base.Pattern == "" {
		base.Pattern = overlay.Pattern
	}

	// Merge format
	if overlay.Format != "" && base.Format == "" {
		base.Format = overlay.Format
	}

	// Merge additionalProperties
	if overlay.AdditionalProperties != nil && base.AdditionalProperties == nil {
		base.AdditionalProperties = overlay.AdditionalProperties
	}

	// Merge items
	if overlay.Items != nil && base.Items == nil {
		base.Items = overlay.Items
	}

	return base
}
