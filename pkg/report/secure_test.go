package report

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/JSLEEKR/schemadiff/pkg/schema"
)

func TestNewReportMetadata(t *testing.T) {
	dir := t.TempDir()
	oldPath := filepath.Join(dir, "old.json")
	newPath := filepath.Join(dir, "new.json")
	os.WriteFile(oldPath, []byte(`{"type":"string"}`), 0644)
	os.WriteFile(newPath, []byte(`{"type":"number"}`), 0644)

	meta := NewReportMetadata(oldPath, newPath, "1.0.0")
	if meta.Tool != "schemadiff" {
		t.Errorf("expected tool=schemadiff, got %s", meta.Tool)
	}
	if meta.OldHash == "" {
		t.Error("expected non-empty old hash")
	}
	if meta.NewHash == "" {
		t.Error("expected non-empty new hash")
	}
	if meta.OldHash == meta.NewHash {
		t.Error("expected different hashes for different files")
	}
	if meta.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestNewReportMetadataNoFiles(t *testing.T) {
	meta := NewReportMetadata("/nonexistent/old.json", "/nonexistent/new.json", "1.0.0")
	if meta.OldHash != "" {
		t.Error("expected empty hash for nonexistent file")
	}
}

func TestValidateSARIFOutputValid(t *testing.T) {
	var buf bytes.Buffer
	r := NewSARIFReporter(&buf, "test.json")
	result := schema.NewDiffResult([]schema.Change{
		{Type: schema.ChangeTypeChanged, Severity: schema.SeverityBreaking, Message: "test"},
	})
	if err := r.Report(result); err != nil {
		t.Fatal(err)
	}

	if err := ValidateSARIFOutput(buf.Bytes()); err != nil {
		t.Errorf("expected valid SARIF: %v", err)
	}
}

func TestValidateSARIFOutputEmpty(t *testing.T) {
	if err := ValidateSARIFOutput([]byte{}); err == nil {
		t.Error("expected error for empty output")
	}
}

func TestValidateSARIFOutputNotJSON(t *testing.T) {
	if err := ValidateSARIFOutput([]byte("not json")); err == nil {
		t.Error("expected error for non-JSON")
	}
}

func TestValidateSARIFOutputInvalidJSON(t *testing.T) {
	if err := ValidateSARIFOutput([]byte(`{"version":"1.0"}`)); err == nil {
		t.Error("expected error for wrong version")
	}
}

func TestValidateSARIFOutputNoRuns(t *testing.T) {
	if err := ValidateSARIFOutput([]byte(`{"version":"2.1.0","runs":[]}`)); err == nil {
		t.Error("expected error for no runs")
	}
}

func TestSanitizeForSARIF(t *testing.T) {
	changes := []schema.Change{
		{Message: "normal message"},
		{Message: "message\x00with\x01control"},
	}
	sanitized := SanitizeForSARIF(changes)
	if sanitized[0].Message != "normal message" {
		t.Error("expected unchanged normal message")
	}
	if sanitized[1].Message != "messagewithcontrol" {
		t.Errorf("expected control chars removed, got %q", sanitized[1].Message)
	}
}

func TestFileHash(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("hello"), 0644)

	hash, err := fileHash(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(hash) != 64 { // SHA-256 hex string
		t.Errorf("expected 64-char hash, got %d", len(hash))
	}
}

func TestFileHashNotFound(t *testing.T) {
	_, err := fileHash("/nonexistent/file")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestFileHashDeterministic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("deterministic"), 0644)

	hash1, _ := fileHash(path)
	hash2, _ := fileHash(path)
	if hash1 != hash2 {
		t.Error("expected same hash for same file")
	}
}
