package schema

import "strings"

// IgnoreConfig holds ignore rules for filtering changes.
type IgnoreConfig struct {
	// Rules lists rule IDs to ignore.
	Rules []string
	// Paths lists path prefixes to ignore.
	Paths []string
	// Severities lists severity levels to ignore.
	Severities []Severity
}

// NewIgnoreConfig creates an IgnoreConfig from rule IDs.
func NewIgnoreConfig(rules []string) *IgnoreConfig {
	return &IgnoreConfig{Rules: rules}
}

// ShouldIgnore checks if a change should be filtered out.
func (ic *IgnoreConfig) ShouldIgnore(c Change) bool {
	if ic == nil {
		return false
	}

	// Check rule ID
	for _, rule := range ic.Rules {
		if string(c.Type) == rule {
			return true
		}
	}

	// Check path prefix
	for _, prefix := range ic.Paths {
		if strings.HasPrefix(c.Path, prefix) {
			return true
		}
	}

	// Check severity
	for _, sev := range ic.Severities {
		if c.Severity == sev {
			return true
		}
	}

	return false
}

// ApplyIgnoreRules filters changes according to the ignore config.
func ApplyIgnoreRules(changes []Change, ic *IgnoreConfig) []Change {
	if ic == nil {
		return changes
	}

	var filtered []Change
	for _, c := range changes {
		if !ic.ShouldIgnore(c) {
			filtered = append(filtered, c)
		}
	}
	return filtered
}
