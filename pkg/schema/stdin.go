package schema

import (
	"fmt"
	"io"
	"os"
)

// ReadFromStdinOrFile reads schema data from stdin (if path is "-") or from a file.
func ReadFromStdinOrFile(path string) ([]byte, error) {
	if path == "-" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("reading stdin: %w", err)
		}
		if len(data) == 0 {
			return nil, fmt.Errorf("no data received from stdin")
		}
		return data, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", path, err)
	}
	return data, nil
}

// ParseFromReader parses a schema from an io.Reader.
func ParseFromReader(r io.Reader) (*Schema, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("empty input")
	}

	// Try JSON first, then YAML
	s, err := ParseJSON(data)
	if err != nil {
		return ParseYAML(data)
	}
	return s, nil
}
