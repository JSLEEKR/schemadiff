package internal

import (
	"encoding/json"
	"testing"

	"github.com/JSLEEKR/schemadiff/pkg/openapi"
	"github.com/JSLEEKR/schemadiff/pkg/report"
	"github.com/JSLEEKR/schemadiff/pkg/schema"
)

func TestIntegrationJSONSchemaFromFile(t *testing.T) {
	old, err := schema.ParseFile("testdata/old_schema.json")
	if err != nil {
		t.Fatal(err)
	}
	new, err := schema.ParseFile("testdata/new_schema.json")
	if err != nil {
		t.Fatal(err)
	}

	result := schema.Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking changes")
	}

	// Verify specific breaking changes
	expectedBreaking := map[schema.ChangeType]bool{
		schema.ChangeMinLengthIncreased:  false,
		schema.ChangeMaxLengthDecreased:  false,
		schema.ChangePropertyRemoved:     false,
		schema.ChangeTypeChanged:         false,
		schema.ChangeEnumValueRemoved:    false,
		schema.ChangeRequiredAdded:       false,
		schema.ChangeAdditionalPropsFalse: false,
		schema.ChangeMinItemsIncreased:   false,
		schema.ChangeMaxItemsDecreased:   false,
	}

	for _, c := range result.Changes {
		if c.Severity == schema.SeverityBreaking {
			expectedBreaking[c.Type] = true
		}
	}

	for ct, found := range expectedBreaking {
		if !found {
			t.Errorf("expected breaking change type %q not found", ct)
		}
	}
}

func TestIntegrationOpenAPIFromFile(t *testing.T) {
	old, err := openapi.ParseFile("testdata/old_api.json")
	if err != nil {
		t.Fatal(err)
	}
	new, err := openapi.ParseFile("testdata/new_api.json")
	if err != nil {
		t.Fatal(err)
	}

	results := openapi.DiffSpecs(old, new)
	if len(results) == 0 {
		t.Fatal("expected changes in OpenAPI diff")
	}

	hasBreaking := false
	for _, r := range results {
		if r.HasBreaking {
			hasBreaking = true
		}
	}
	if !hasBreaking {
		t.Error("expected breaking changes in OpenAPI diff")
	}
}

func TestIntegrationJSONOutput(t *testing.T) {
	old := &schema.Schema{
		Type: "object",
		Properties: map[string]*schema.Schema{
			"name": {Type: "string"},
		},
		Required: []string{"name"},
	}
	new := &schema.Schema{
		Type: "object",
		Properties: map[string]*schema.Schema{
			"name": {Type: "integer"},
		},
		Required: []string{"name"},
	}

	result := schema.Diff(old, new)

	// Verify JSON output is valid
	var buf []byte
	buf, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}

	var parsed schema.DiffResult
	if err := json.Unmarshal(buf, &parsed); err != nil {
		t.Fatal(err)
	}
	if !parsed.HasBreaking {
		t.Error("expected breaking in parsed JSON")
	}
}

func TestIntegrationSARIFOutput(t *testing.T) {
	result := schema.NewDiffResult([]schema.Change{
		{
			Path:     "type",
			Type:     schema.ChangeTypeChanged,
			Severity: schema.SeverityBreaking,
			Message:  "Type changed",
		},
	})

	var buf []byte
	var sarifReport report.SARIFReport

	r := report.NewSARIFReporter(nil, "test.json")
	_ = r // We just verify SARIF structure manually
	_ = buf
	_ = sarifReport

	// Build the SARIF via report package
	var out = new(testBuffer)
	sarif := report.NewSARIFReporter(out, "test.json")
	if err := sarif.Report(result); err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal(out.Bytes(), &sarifReport); err != nil {
		t.Fatalf("invalid SARIF output: %v", err)
	}
	if sarifReport.Version != "2.1.0" {
		t.Errorf("expected SARIF 2.1.0, got %s", sarifReport.Version)
	}
}

func TestIntegrationNoBreakingChanges(t *testing.T) {
	old := &schema.Schema{
		Type: "object",
		Properties: map[string]*schema.Schema{
			"name": {Type: "string"},
		},
		Required: []string{"name"},
	}
	new := &schema.Schema{
		Type: "object",
		Properties: map[string]*schema.Schema{
			"name":  {Type: "string"},
			"email": {Type: "string"}, // new optional property
		},
		Required: []string{"name"},
	}

	result := schema.Diff(old, new)
	if result.HasBreaking {
		t.Error("expected no breaking changes")
	}
	if result.Summary.Info != 1 {
		t.Errorf("expected 1 info change, got %d", result.Summary.Info)
	}
}

func TestIntegrationOpenAPIExtraction(t *testing.T) {
	spec, err := openapi.ParseFile("testdata/old_api.json")
	if err != nil {
		t.Fatal(err)
	}

	schemas := openapi.ExtractSchemas(spec)
	if len(schemas) == 0 {
		t.Error("expected extracted schemas")
	}

	// Should have component schema + endpoint schemas
	locations := map[string]bool{}
	for _, s := range schemas {
		locations[s.Location] = true
	}
	if !locations["component"] {
		t.Error("expected component schemas")
	}
	if !locations["request"] {
		t.Error("expected request schemas")
	}
	if !locations["response"] {
		t.Error("expected response schemas")
	}
}

func TestIntegrationMultipleOutputFormats(t *testing.T) {
	oldS := &schema.Schema{Type: "string"}
	newS := &schema.Schema{Type: "number"}
	result := schema.Diff(oldS, newS)

	// Text output
	textBuf := &testBuffer{}
	text := report.NewTextReporter(textBuf, false)
	if err := text.Report("Test", result); err != nil {
		t.Fatal(err)
	}
	if textBuf.Len() == 0 {
		t.Error("expected non-empty text output")
	}

	// JSON output
	jsonBuf := &testBuffer{}
	jr := report.NewJSONReporter(jsonBuf, true)
	if err := jr.Report(result); err != nil {
		t.Fatal(err)
	}
	var jsonReport report.JSONReport
	if err := json.Unmarshal(jsonBuf.Bytes(), &jsonReport); err != nil {
		t.Fatal(err)
	}

	// SARIF output
	sarifBuf := &testBuffer{}
	sr := report.NewSARIFReporter(sarifBuf, "test.json")
	if err := sr.Report(result); err != nil {
		t.Fatal(err)
	}
	var sarifReport report.SARIFReport
	if err := json.Unmarshal(sarifBuf.Bytes(), &sarifReport); err != nil {
		t.Fatal(err)
	}
}

func TestIntegrationRefResolution(t *testing.T) {
	old := &schema.Schema{
		Type: "object",
		Properties: map[string]*schema.Schema{
			"address": {Ref: "#/definitions/Address"},
		},
		Definitions: map[string]*schema.Schema{
			"Address": {
				Type: "object",
				Properties: map[string]*schema.Schema{
					"street": {Type: "string"},
					"city":   {Type: "string"},
				},
				Required: []string{"street"},
			},
		},
	}

	new := &schema.Schema{
		Type: "object",
		Properties: map[string]*schema.Schema{
			"address": {Ref: "#/definitions/Address"},
		},
		Definitions: map[string]*schema.Schema{
			"Address": {
				Type: "object",
				Properties: map[string]*schema.Schema{
					"street": {Type: "integer"}, // type changed
					"city":   {Type: "string"},
				},
				Required: []string{"street", "city"}, // city added
			},
		},
	}

	result := schema.Diff(old, new)
	if !result.HasBreaking {
		t.Error("expected breaking changes through ref resolution")
	}
}

// testBuffer is a simple bytes.Buffer-like writer for tests.
type testBuffer struct {
	data []byte
}

func (b *testBuffer) Write(p []byte) (n int, err error) {
	b.data = append(b.data, p...)
	return len(p), nil
}

func (b *testBuffer) Bytes() []byte {
	return b.data
}

func (b *testBuffer) Len() int {
	return len(b.data)
}
