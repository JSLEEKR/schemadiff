package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.FailOn != "breaking" {
		t.Errorf("expected fail_on=breaking, got %s", cfg.FailOn)
	}
	if cfg.Output != "text" {
		t.Errorf("expected output=text, got %s", cfg.Output)
	}
	if cfg.Format != "auto" {
		t.Errorf("expected format=auto, got %s", cfg.Format)
	}
}

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".schemadiff.json")
	data := []byte(`{
		"version": "1",
		"fail_on": "warning",
		"output": "json",
		"ignore": ["format-changed"],
		"paths": [
			{"name": "API", "old": "old.json", "new": "new.json"}
		]
	}`)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.FailOn != "warning" {
		t.Errorf("expected fail_on=warning, got %s", cfg.FailOn)
	}
	if cfg.Output != "json" {
		t.Errorf("expected output=json, got %s", cfg.Output)
	}
	if len(cfg.Ignore) != 1 || cfg.Ignore[0] != "format-changed" {
		t.Errorf("expected ignore=[format-changed], got %v", cfg.Ignore)
	}
	if len(cfg.Paths) != 1 {
		t.Errorf("expected 1 path config, got %d", len(cfg.Paths))
	}
}

func TestLoadConfigNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/.schemadiff.json")
	if err == nil {
		t.Error("expected error for nonexistent config")
	}
}

func TestLoadConfigInvalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".schemadiff.json")
	if err := os.WriteFile(path, []byte(`{invalid`), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadConfig(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestFindConfig(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, ".schemadiff.json")
	if err := os.WriteFile(cfgPath, []byte(`{"version":"1"}`), 0644); err != nil {
		t.Fatal(err)
	}

	found, err := FindConfig(dir)
	if err != nil {
		t.Fatal(err)
	}
	if found != cfgPath {
		t.Errorf("expected %s, got %s", cfgPath, found)
	}
}

func TestFindConfigNotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := FindConfig(dir)
	if err == nil {
		t.Error("expected error when no config found")
	}
}

func TestFindConfigWalksUp(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "a", "b")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(dir, ".schemadiff.json")
	if err := os.WriteFile(cfgPath, []byte(`{"version":"1"}`), 0644); err != nil {
		t.Fatal(err)
	}

	found, err := FindConfig(subdir)
	if err != nil {
		t.Fatal(err)
	}
	if found != cfgPath {
		t.Errorf("expected %s, got %s", cfgPath, found)
	}
}

func TestLoadConfigOrDefault(t *testing.T) {
	dir := t.TempDir()
	cfg := LoadConfigOrDefault(dir)
	if cfg.FailOn != "breaking" {
		t.Error("expected default fail_on")
	}
}

func TestLoadConfigOrDefaultWithFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".schemadiff.json")
	if err := os.WriteFile(path, []byte(`{"fail_on":"any"}`), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := LoadConfigOrDefault(dir)
	if cfg.FailOn != "any" {
		t.Errorf("expected fail_on=any, got %s", cfg.FailOn)
	}
	// Check defaults filled in
	if cfg.Output != "text" {
		t.Errorf("expected output=text default, got %s", cfg.Output)
	}
}

func TestIsIgnored(t *testing.T) {
	cfg := &Config{Ignore: []string{"format-changed", "pattern-changed"}}
	if !cfg.IsIgnored("format-changed") {
		t.Error("expected format-changed to be ignored")
	}
	if cfg.IsIgnored("type-changed") {
		t.Error("expected type-changed not to be ignored")
	}
}

func TestIsIgnoredEmpty(t *testing.T) {
	cfg := &Config{}
	if cfg.IsIgnored("anything") {
		t.Error("expected nothing ignored with empty list")
	}
}

func TestConfigFileNames(t *testing.T) {
	if len(ConfigFileNames) == 0 {
		t.Error("expected config file names")
	}
}

func TestPathConfig(t *testing.T) {
	cfg := &Config{
		Paths: []PathConfig{
			{Name: "API", Old: "old.json", New: "new.json", Format: "openapi"},
		},
	}
	if cfg.Paths[0].Name != "API" {
		t.Error("expected path name")
	}
	if cfg.Paths[0].Format != "openapi" {
		t.Error("expected path format override")
	}
}
