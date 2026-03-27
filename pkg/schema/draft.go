package schema

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Draft represents a JSON Schema draft version.
type Draft string

const (
	DraftUnknown Draft = "unknown"
	Draft4       Draft = "draft-04"
	Draft6       Draft = "draft-06"
	Draft7       Draft = "draft-07"
	Draft2019    Draft = "2019-09"
	Draft2020    Draft = "2020-12"
)

// DetectDraft detects the JSON Schema draft version from raw JSON.
func DetectDraft(data []byte) Draft {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return DraftUnknown
	}

	schemaField, ok := raw["$schema"]
	if !ok {
		return DraftUnknown
	}

	var schemaURI string
	if err := json.Unmarshal(schemaField, &schemaURI); err != nil {
		return DraftUnknown
	}

	return parseDraftURI(schemaURI)
}

func parseDraftURI(uri string) Draft {
	uri = strings.ToLower(uri)

	if strings.Contains(uri, "draft-04") || strings.Contains(uri, "draft/4") {
		return Draft4
	}
	if strings.Contains(uri, "draft-06") || strings.Contains(uri, "draft/6") {
		return Draft6
	}
	if strings.Contains(uri, "draft-07") || strings.Contains(uri, "draft/7") {
		return Draft7
	}
	if strings.Contains(uri, "2019-09") {
		return Draft2019
	}
	if strings.Contains(uri, "2020-12") {
		return Draft2020
	}

	return DraftUnknown
}

// DraftInfo provides human-readable info about a draft version.
func DraftInfo(d Draft) string {
	switch d {
	case Draft4:
		return "JSON Schema Draft 4 (2013)"
	case Draft6:
		return "JSON Schema Draft 6 (2017)"
	case Draft7:
		return "JSON Schema Draft 7 (2018)"
	case Draft2019:
		return "JSON Schema 2019-09 (Draft 8)"
	case Draft2020:
		return "JSON Schema 2020-12 (latest)"
	default:
		return "Unknown JSON Schema version"
	}
}

// IsDraftCompatible checks if two draft versions are compatible for comparison.
func IsDraftCompatible(a, b Draft) (bool, string) {
	if a == DraftUnknown || b == DraftUnknown {
		return true, "unknown draft versions, comparison may be imprecise"
	}
	if a == b {
		return true, ""
	}

	// Major version jumps may have different semantics
	majorA := draftMajor(a)
	majorB := draftMajor(b)

	if majorA != majorB {
		return true, fmt.Sprintf("different schema drafts: %s vs %s, some rules may not apply", DraftInfo(a), DraftInfo(b))
	}

	return true, ""
}

func draftMajor(d Draft) int {
	switch d {
	case Draft4:
		return 1
	case Draft6, Draft7:
		return 2
	case Draft2019, Draft2020:
		return 3
	default:
		return 0
	}
}
