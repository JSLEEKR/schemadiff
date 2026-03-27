package schema

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSecureTempFile(t *testing.T) {
	dir := t.TempDir()
	f, err := SecureTempFile(dir, "schemadiff-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	if f == nil {
		t.Fatal("expected non-nil file")
	}

	// Write some data
	_, err = f.Write([]byte("test data"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestSecureTempFileDefaultDir(t *testing.T) {
	f, err := SecureTempFile("", "schemadiff-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	if f == nil {
		t.Fatal("expected non-nil file")
	}
}

func TestSecureWriteFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "output", "report.json")

	err := SecureWriteFile(path, []byte(`{"test": true}`))
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"test": true}` {
		t.Errorf("unexpected data: %s", data)
	}
}

func TestSecureWriteFileExistingDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.json")

	err := SecureWriteFile(path, []byte("test"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCleanupTempFiles(t *testing.T) {
	dir := t.TempDir()

	// Create some temp files
	for i := 0; i < 3; i++ {
		f, err := os.CreateTemp(dir, "schemadiff-*")
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
	}

	err := CleanupTempFiles(dir, "schemadiff-*")
	if err != nil {
		t.Fatal(err)
	}

	// Verify cleanup
	matches, _ := filepath.Glob(filepath.Join(dir, "schemadiff-*"))
	if len(matches) != 0 {
		t.Errorf("expected 0 remaining files, got %d", len(matches))
	}
}

func TestCleanupTempFilesNoMatch(t *testing.T) {
	dir := t.TempDir()
	err := CleanupTempFiles(dir, "nonexistent-*")
	if err != nil {
		t.Errorf("expected no error for no matches: %v", err)
	}
}
