package schema

// FilterBySeverity returns only changes at or above the given severity level.
func FilterBySeverity(changes []Change, minSeverity Severity) []Change {
	var filtered []Change
	for _, c := range changes {
		if c.Severity >= minSeverity {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

// FilterByType returns only changes of the specified types.
func FilterByType(changes []Change, types ...ChangeType) []Change {
	typeSet := make(map[ChangeType]bool, len(types))
	for _, t := range types {
		typeSet[t] = true
	}

	var filtered []Change
	for _, c := range changes {
		if typeSet[c.Type] {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

// FilterByPath returns only changes matching a path prefix.
func FilterByPath(changes []Change, pathPrefix string) []Change {
	var filtered []Change
	for _, c := range changes {
		if len(c.Path) >= len(pathPrefix) && c.Path[:len(pathPrefix)] == pathPrefix {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

// GroupBySeverity groups changes by severity level.
func GroupBySeverity(changes []Change) map[Severity][]Change {
	groups := make(map[Severity][]Change)
	for _, c := range changes {
		groups[c.Severity] = append(groups[c.Severity], c)
	}
	return groups
}

// GroupByPath groups changes by their top-level path segment.
func GroupByPath(changes []Change) map[string][]Change {
	groups := make(map[string][]Change)
	for _, c := range changes {
		topLevel := c.Path
		for i, ch := range c.Path {
			if ch == '.' {
				topLevel = c.Path[:i]
				break
			}
		}
		groups[topLevel] = append(groups[topLevel], c)
	}
	return groups
}
