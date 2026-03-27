package report

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/JSLEEKR/schemadiff/pkg/schema"
)

// ReportMetadata provides context about a diff report.
type ReportMetadata struct {
	// Timestamp is when the report was generated.
	Timestamp time.Time `json:"timestamp"`
	// OldFile is the path to the old schema file.
	OldFile string `json:"old_file"`
	// NewFile is the path to the new schema file.
	NewFile string `json:"new_file"`
	// OldHash is the SHA-256 hash of the old file.
	OldHash string `json:"old_hash,omitempty"`
	// NewHash is the SHA-256 hash of the new file.
	NewHash string `json:"new_hash,omitempty"`
	// Tool is the tool name.
	Tool string `json:"tool"`
	// Version is the tool version.
	Version string `json:"version"`
}

// NewReportMetadata creates metadata for a report.
func NewReportMetadata(oldFile, newFile, version string) *ReportMetadata {
	meta := &ReportMetadata{
		Timestamp: time.Now().UTC(),
		OldFile:   oldFile,
		NewFile:   newFile,
		Tool:      "schemadiff",
		Version:   version,
	}

	if hash, err := fileHash(oldFile); err == nil {
		meta.OldHash = hash
	}
	if hash, err := fileHash(newFile); err == nil {
		meta.NewHash = hash
	}

	return meta
}

func fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// ValidateSARIFOutput checks that SARIF output is well-formed.
func ValidateSARIFOutput(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("empty SARIF output")
	}
	if data[0] != '{' {
		return fmt.Errorf("SARIF output must be a JSON object")
	}

	// Basic structure check
	var report SARIFReport
	if err := parseJSON(data, &report); err != nil {
		return fmt.Errorf("invalid SARIF JSON: %w", err)
	}

	if report.Version != "2.1.0" {
		return fmt.Errorf("unsupported SARIF version: %s", report.Version)
	}
	if len(report.Runs) == 0 {
		return fmt.Errorf("SARIF report must have at least one run")
	}

	return nil
}

func parseJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// SanitizeForSARIF ensures change messages are safe for SARIF output.
func SanitizeForSARIF(changes []schema.Change) []schema.Change {
	return schema.SanitizeOutput(changes)
}
