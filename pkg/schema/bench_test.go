package schema

import (
	"encoding/json"
	"fmt"
	"testing"
)

func BenchmarkDiffSimple(b *testing.B) {
	old := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"name": {Type: "string"},
			"age":  {Type: "integer"},
		},
		Required: []string{"name"},
	}
	new := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"name": {Type: "string"},
			"age":  {Type: "string"}, // changed
		},
		Required: []string{"name", "age"}, // added
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Diff(old, new)
	}
}

func BenchmarkDiffLargeSchema(b *testing.B) {
	old := generateLargeSchema(100)
	new := generateLargeSchema(100)
	// Modify some properties to create changes
	new.Properties["prop_50"] = &Schema{Type: "integer"}
	new.Required = append(new.Required, "prop_99")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Diff(old, new)
	}
}

func BenchmarkDiffDeepNesting(b *testing.B) {
	old := generateDeepSchema(20)
	new := generateDeepSchema(20)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Diff(old, new)
	}
}

func BenchmarkParseJSON(b *testing.B) {
	s := generateLargeSchema(100)
	data, _ := json.Marshal(s)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseJSON(data)
	}
}

func BenchmarkResolveRefs(b *testing.B) {
	s := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"addr1": {Ref: "#/definitions/Address"},
			"addr2": {Ref: "#/definitions/Address"},
			"addr3": {Ref: "#/definitions/Address"},
		},
		Definitions: map[string]*Schema{
			"Address": {
				Type: "object",
				Properties: map[string]*Schema{
					"street": {Type: "string"},
					"city":   {Type: "string"},
					"zip":    {Type: "string"},
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ResolveRefs(s)
	}
}

func BenchmarkMergeAllOf(b *testing.B) {
	s := &Schema{
		AllOf: []*Schema{
			{Type: "object", Properties: map[string]*Schema{"a": {Type: "string"}}},
			{Properties: map[string]*Schema{"b": {Type: "string"}}},
			{Properties: map[string]*Schema{"c": {Type: "string"}}},
			{Required: []string{"a", "b"}},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MergeAllOf(s)
	}
}

func BenchmarkValidateSchema(b *testing.B) {
	s := generateLargeSchema(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateSchema(s)
	}
}

func BenchmarkFilterBySeverity(b *testing.B) {
	changes := make([]Change, 1000)
	for i := range changes {
		changes[i] = Change{Severity: Severity(i % 3)}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FilterBySeverity(changes, SeverityBreaking)
	}
}

func generateLargeSchema(numProps int) *Schema {
	props := make(map[string]*Schema, numProps)
	for i := 0; i < numProps; i++ {
		props[fmt.Sprintf("prop_%d", i)] = &Schema{
			Type:      "string",
			MinLength: intPtr(1),
			MaxLength: intPtr(100),
		}
	}
	return &Schema{
		Type:       "object",
		Properties: props,
		Required:   []string{"prop_0", "prop_1"},
	}
}

func generateDeepSchema(depth int) *Schema {
	if depth == 0 {
		return &Schema{Type: "string"}
	}
	return &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"nested": generateDeepSchema(depth - 1),
		},
	}
}
