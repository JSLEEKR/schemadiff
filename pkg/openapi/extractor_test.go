package openapi

import (
	"testing"

	"github.com/JSLEEKR/schemadiff/pkg/schema"
)

func makeSpec(schemas map[string]*schema.Schema) *Spec {
	return &Spec{
		OpenAPI:    "3.0.0",
		Info:       Info{Title: "Test", Version: "1.0.0"},
		Paths:      map[string]*PathItem{},
		Components: &Components{Schemas: schemas},
	}
}

func TestExtractComponentSchemas(t *testing.T) {
	spec := makeSpec(map[string]*schema.Schema{
		"User": {Type: "object"},
		"Post": {Type: "object"},
	})

	schemas := ExtractSchemas(spec)
	if len(schemas) != 2 {
		t.Fatalf("expected 2 schemas, got %d", len(schemas))
	}
	for _, s := range schemas {
		if s.Location != "component" {
			t.Errorf("expected location=component, got %s", s.Location)
		}
	}
}

func TestExtractRequestSchemas(t *testing.T) {
	spec := &Spec{
		OpenAPI: "3.0.0",
		Info:    Info{Title: "Test", Version: "1.0.0"},
		Paths: map[string]*PathItem{
			"/users": {
				Post: &Operation{
					RequestBody: &RequestBody{
						Content: map[string]*Content{
							"application/json": {
								Schema: &schema.Schema{Type: "object"},
							},
						},
					},
					Responses: map[string]*Response{},
				},
			},
		},
	}

	schemas := ExtractSchemas(spec)
	found := false
	for _, s := range schemas {
		if s.Location == "request" {
			found = true
			if s.Method != "POST" {
				t.Errorf("expected method=POST, got %s", s.Method)
			}
		}
	}
	if !found {
		t.Error("expected request schema")
	}
}

func TestExtractResponseSchemas(t *testing.T) {
	spec := &Spec{
		OpenAPI: "3.0.0",
		Info:    Info{Title: "Test", Version: "1.0.0"},
		Paths: map[string]*PathItem{
			"/users": {
				Get: &Operation{
					Responses: map[string]*Response{
						"200": {
							Description: "Success",
							Content: map[string]*Content{
								"application/json": {
									Schema: &schema.Schema{Type: "array"},
								},
							},
						},
					},
				},
			},
		},
	}

	schemas := ExtractSchemas(spec)
	found := false
	for _, s := range schemas {
		if s.Location == "response" {
			found = true
		}
	}
	if !found {
		t.Error("expected response schema")
	}
}

func TestExtractParameterSchemas(t *testing.T) {
	spec := &Spec{
		OpenAPI: "3.0.0",
		Info:    Info{Title: "Test", Version: "1.0.0"},
		Paths: map[string]*PathItem{
			"/users/{id}": {
				Get: &Operation{
					Parameters: []Parameter{
						{Name: "id", In: "path", Schema: &schema.Schema{Type: "string"}},
					},
					Responses: map[string]*Response{},
				},
			},
		},
	}

	schemas := ExtractSchemas(spec)
	found := false
	for _, s := range schemas {
		if s.Location == "parameter" {
			found = true
			if s.Path != "/users/{id}" {
				t.Errorf("expected path=/users/{id}, got %s", s.Path)
			}
		}
	}
	if !found {
		t.Error("expected parameter schema")
	}
}

func TestExtractSchemasEmpty(t *testing.T) {
	spec := &Spec{
		OpenAPI: "3.0.0",
		Info:    Info{Title: "Test", Version: "1.0.0"},
		Paths:   map[string]*PathItem{},
	}

	schemas := ExtractSchemas(spec)
	if len(schemas) != 0 {
		t.Errorf("expected 0 schemas, got %d", len(schemas))
	}
}

func TestDiffSpecsComponentChanged(t *testing.T) {
	old := makeSpec(map[string]*schema.Schema{
		"User": {
			Type: "object",
			Properties: map[string]*schema.Schema{
				"name":  {Type: "string"},
				"email": {Type: "string"},
			},
			Required: []string{"name"},
		},
	})

	new := makeSpec(map[string]*schema.Schema{
		"User": {
			Type: "object",
			Properties: map[string]*schema.Schema{
				"name": {Type: "string"},
			},
			Required: []string{"name", "email"}, // email made required but also removed from props
		},
	})

	results := DiffSpecs(old, new)
	if _, ok := results["User"]; !ok {
		t.Fatal("expected changes for User schema")
	}
	if !results["User"].HasBreaking {
		t.Error("expected breaking changes")
	}
}

func TestDiffSpecsSchemaRemoved(t *testing.T) {
	old := makeSpec(map[string]*schema.Schema{
		"User": {Type: "object"},
	})
	new := makeSpec(map[string]*schema.Schema{})

	results := DiffSpecs(old, new)
	if _, ok := results["User"]; !ok {
		t.Fatal("expected changes for removed User schema")
	}
	if !results["User"].HasBreaking {
		t.Error("expected breaking for schema removal")
	}
}

func TestDiffSpecsSchemaAdded(t *testing.T) {
	old := makeSpec(map[string]*schema.Schema{})
	new := makeSpec(map[string]*schema.Schema{
		"User": {Type: "object"},
	})

	results := DiffSpecs(old, new)
	if _, ok := results["User"]; !ok {
		t.Fatal("expected changes for added User schema")
	}
	if results["User"].HasBreaking {
		t.Error("expected no breaking for schema addition")
	}
}

func TestDiffSpecsNoChanges(t *testing.T) {
	spec := makeSpec(map[string]*schema.Schema{
		"User": {Type: "object", Properties: map[string]*schema.Schema{"name": {Type: "string"}}},
	})

	results := DiffSpecs(spec, spec)
	if len(results) != 0 {
		t.Errorf("expected no changes, got %d", len(results))
	}
}

func TestExtractSchemasSorted(t *testing.T) {
	spec := makeSpec(map[string]*schema.Schema{
		"Zebra":    {Type: "object"},
		"Apple":    {Type: "object"},
		"Mango":    {Type: "object"},
	})

	schemas := ExtractSchemas(spec)
	for i := 1; i < len(schemas); i++ {
		if schemas[i].Name < schemas[i-1].Name {
			t.Errorf("schemas not sorted: %s before %s", schemas[i-1].Name, schemas[i].Name)
		}
	}
}

func TestDiffSpecsMultipleMethods(t *testing.T) {
	old := &Spec{
		OpenAPI: "3.0.0",
		Info:    Info{Title: "Test", Version: "1.0.0"},
		Paths: map[string]*PathItem{
			"/items": {
				Get: &Operation{
					Responses: map[string]*Response{
						"200": {
							Description: "OK",
							Content: map[string]*Content{
								"application/json": {
									Schema: &schema.Schema{Type: "array", Items: &schema.Schema{Type: "string"}},
								},
							},
						},
					},
				},
				Post: &Operation{
					RequestBody: &RequestBody{
						Content: map[string]*Content{
							"application/json": {
								Schema: &schema.Schema{
									Type: "object",
									Properties: map[string]*schema.Schema{
										"name": {Type: "string"},
									},
								},
							},
						},
					},
					Responses: map[string]*Response{},
				},
			},
		},
	}

	new := &Spec{
		OpenAPI: "3.0.0",
		Info:    Info{Title: "Test", Version: "1.0.0"},
		Paths: map[string]*PathItem{
			"/items": {
				Get: &Operation{
					Responses: map[string]*Response{
						"200": {
							Description: "OK",
							Content: map[string]*Content{
								"application/json": {
									Schema: &schema.Schema{Type: "array", Items: &schema.Schema{Type: "integer"}},
								},
							},
						},
					},
				},
				Post: &Operation{
					RequestBody: &RequestBody{
						Content: map[string]*Content{
							"application/json": {
								Schema: &schema.Schema{
									Type: "object",
									Properties: map[string]*schema.Schema{
										"name": {Type: "string"},
									},
									Required: []string{"name"},
								},
							},
						},
					},
					Responses: map[string]*Response{},
				},
			},
		},
	}

	results := DiffSpecs(old, new)
	if len(results) == 0 {
		t.Error("expected changes for multiple method spec diff")
	}

	hasBreaking := false
	for _, r := range results {
		if r.HasBreaking {
			hasBreaking = true
		}
	}
	if !hasBreaking {
		t.Error("expected at least one breaking change")
	}
}
