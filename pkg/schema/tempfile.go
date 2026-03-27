package schema

import (
	"fmt"
	"os"
	"path/filepath"
)

// SecureTempFile creates a temporary file with restricted permissions.
func SecureTempFile(dir, pattern string) (*os.File, error) {
	if dir == "" {
		dir = os.TempDir()
	}

	// Ensure temp directory exists
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("creating temp directory: %w", err)
	}

	f, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return nil, fmt.Errorf("creating temp file: %w", err)
	}

	// Set restrictive permissions (owner read/write only)
	if err := os.Chmod(f.Name(), 0600); err != nil {
		f.Close()
		os.Remove(f.Name())
		return nil, fmt.Errorf("setting temp file permissions: %w", err)
	}

	return f, nil
}

// SecureWriteFile writes data to a file with restricted permissions.
func SecureWriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

// CleanupTempFiles removes temporary files matching a pattern.
func CleanupTempFiles(dir, pattern string) error {
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return fmt.Errorf("globbing temp files: %w", err)
	}

	var lastErr error
	for _, m := range matches {
		if err := os.Remove(m); err != nil {
			lastErr = err
		}
	}
	return lastErr
}
