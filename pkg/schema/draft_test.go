package schema

import "testing"

func TestDetectDraft(t *testing.T) {
	tests := []struct {
		name string
		data string
		want Draft
	}{
		{"draft4", `{"$schema":"http://json-schema.org/draft-04/schema#"}`, Draft4},
		{"draft6", `{"$schema":"http://json-schema.org/draft-06/schema#"}`, Draft6},
		{"draft7", `{"$schema":"http://json-schema.org/draft-07/schema#"}`, Draft7},
		{"2019-09", `{"$schema":"https://json-schema.org/draft/2019-09/schema"}`, Draft2019},
		{"2020-12", `{"$schema":"https://json-schema.org/draft/2020-12/schema"}`, Draft2020},
		{"no schema", `{"type":"string"}`, DraftUnknown},
		{"invalid", `{invalid`, DraftUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectDraft([]byte(tt.data))
			if got != tt.want {
				t.Errorf("DetectDraft() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDraftURI(t *testing.T) {
	tests := []struct {
		uri  string
		want Draft
	}{
		{"http://json-schema.org/draft-04/schema#", Draft4},
		{"http://json-schema.org/draft-06/schema#", Draft6},
		{"http://json-schema.org/draft-07/schema#", Draft7},
		{"https://json-schema.org/draft/2019-09/schema", Draft2019},
		{"https://json-schema.org/draft/2020-12/schema", Draft2020},
		{"http://example.com/custom-schema", DraftUnknown},
		{"", DraftUnknown},
	}

	for _, tt := range tests {
		got := parseDraftURI(tt.uri)
		if got != tt.want {
			t.Errorf("parseDraftURI(%q) = %v, want %v", tt.uri, got, tt.want)
		}
	}
}

func TestDraftInfo(t *testing.T) {
	drafts := []Draft{Draft4, Draft6, Draft7, Draft2019, Draft2020, DraftUnknown}
	for _, d := range drafts {
		info := DraftInfo(d)
		if info == "" {
			t.Errorf("expected non-empty info for %v", d)
		}
	}
}

func TestIsDraftCompatible(t *testing.T) {
	// Same draft
	ok, msg := IsDraftCompatible(Draft7, Draft7)
	if !ok {
		t.Error("expected compatible for same draft")
	}
	if msg != "" {
		t.Error("expected no message for same draft")
	}

	// Unknown
	ok, msg = IsDraftCompatible(DraftUnknown, Draft7)
	if !ok {
		t.Error("expected compatible with unknown")
	}
	if msg == "" {
		t.Error("expected message for unknown draft")
	}

	// Different major
	ok, msg = IsDraftCompatible(Draft4, Draft2020)
	if !ok {
		t.Error("expected compatible across versions")
	}
	if msg == "" {
		t.Error("expected warning for different major versions")
	}

	// Same major
	ok, msg = IsDraftCompatible(Draft6, Draft7)
	if !ok {
		t.Error("expected compatible within major")
	}
}

func TestDraftMajor(t *testing.T) {
	if draftMajor(Draft4) != 1 {
		t.Error("expected major=1 for draft4")
	}
	if draftMajor(Draft6) != 2 {
		t.Error("expected major=2 for draft6")
	}
	if draftMajor(Draft7) != 2 {
		t.Error("expected major=2 for draft7")
	}
	if draftMajor(Draft2019) != 3 {
		t.Error("expected major=3 for 2019")
	}
	if draftMajor(Draft2020) != 3 {
		t.Error("expected major=3 for 2020")
	}
	if draftMajor(DraftUnknown) != 0 {
		t.Error("expected major=0 for unknown")
	}
}
