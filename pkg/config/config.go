// Package config provides project-level configuration for schemadiff.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents project-level schemadiff configuration.
type Config struct {
	// Version is the config file version.
	Version string `json:"version"`
	// FailOn sets the default fail condition.
	FailOn string `json:"fail_on,omitempty"`
	// Output sets the default output format.
	Output string `json:"output,omitempty"`
	// Format sets the default input format.
	Format string `json:"format,omitempty"`
	// NoColor disables colored output.
	NoColor bool `json:"no_color,omitempty"`
	// Ignore lists rule IDs to skip.
	Ignore []string `json:"ignore,omitempty"`
	// Paths defines schema file pairs to check.
	Paths []PathConfig `json:"paths,omitempty"`
}

// PathConfig defines a pair of schema files to compare.
type PathConfig struct {
	// Name is a human-readable label.
	Name string `json:"name,omitempty"`
	// Old is the path to the old schema file.
	Old string `json:"old"`
	// New is the path to the new schema file.
	New string `json:"new"`
	// Format overrides the input format for this pair.
	Format string `json:"format,omitempty"`
}

// DefaultConfig returns a config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Version: "1",
		FailOn:  "breaking",
		Output:  "text",
		Format:  "auto",
	}
}

// ConfigFileNames lists possible config file names in priority order.
var ConfigFileNames = []string{
	".schemadiff.json",
	".schemadiff.yaml",
	".schemadiff.yml",
	"schemadiff.json",
}

// FindConfig searches for a config file starting from dir and walking up.
func FindConfig(dir string) (string, error) {
	for {
		for _, name := range ConfigFileNames {
			path := filepath.Join(dir, name)
			if _, err := os.Stat(path); err == nil {
				return path, nil
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("no config file found")
}

// LoadConfig reads and parses a config file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return &cfg, nil
}

// LoadConfigOrDefault tries to find and load a config file, falling back to defaults.
func LoadConfigOrDefault(dir string) *Config {
	path, err := FindConfig(dir)
	if err != nil {
		return DefaultConfig()
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		return DefaultConfig()
	}

	// Fill in defaults for empty fields
	if cfg.FailOn == "" {
		cfg.FailOn = "breaking"
	}
	if cfg.Output == "" {
		cfg.Output = "text"
	}
	if cfg.Format == "" {
		cfg.Format = "auto"
	}

	return cfg
}

// IsIgnored checks if a rule ID is in the ignore list.
func (c *Config) IsIgnored(ruleID string) bool {
	for _, id := range c.Ignore {
		if id == ruleID {
			return true
		}
	}
	return false
}
