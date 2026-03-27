package openapi

import (
	"fmt"
	"sort"

	"github.com/JSLEEKR/schemadiff/pkg/schema"
)

// SchemaLocation identifies where a schema was found in an OpenAPI spec.
type SchemaLocation struct {
	// Name is the schema identifier (e.g., "User", "POST /users request").
	Name string
	// Path is the OpenAPI path (e.g., "/users").
	Path string
	// Method is the HTTP method (e.g., "GET", "POST").
	Method string
	// Location is where in the operation the schema was found.
	Location string // "component", "request", "response", "parameter"
	// Schema is the extracted schema.
	Schema *schema.Schema
}

// ExtractSchemas extracts all schemas from an OpenAPI spec.
func ExtractSchemas(spec *Spec) []SchemaLocation {
	var schemas []SchemaLocation

	// Component schemas
	if spec.Components != nil {
		for name, s := range spec.Components.Schemas {
			schemas = append(schemas, SchemaLocation{
				Name:     name,
				Location: "component",
				Schema:   s,
			})
		}
	}

	// Path schemas
	for path, item := range spec.Paths {
		ops := map[string]*Operation{
			"GET":    item.Get,
			"POST":   item.Post,
			"PUT":    item.Put,
			"PATCH":  item.Patch,
			"DELETE": item.Delete,
		}

		for method, op := range ops {
			if op == nil {
				continue
			}

			// Request body schemas
			if op.RequestBody != nil {
				for contentType, content := range op.RequestBody.Content {
					if content.Schema != nil {
						name := fmt.Sprintf("%s %s request (%s)", method, path, contentType)
						schemas = append(schemas, SchemaLocation{
							Name:     name,
							Path:     path,
							Method:   method,
							Location: "request",
							Schema:   content.Schema,
						})
					}
				}
			}

			// Response schemas
			for statusCode, resp := range op.Responses {
				for contentType, content := range resp.Content {
					if content.Schema != nil {
						name := fmt.Sprintf("%s %s %s response (%s)", method, path, statusCode, contentType)
						schemas = append(schemas, SchemaLocation{
							Name:     name,
							Path:     path,
							Method:   method,
							Location: "response",
							Schema:   content.Schema,
						})
					}
				}
			}

			// Parameter schemas
			for _, param := range op.Parameters {
				if param.Schema != nil {
					name := fmt.Sprintf("%s %s param:%s", method, path, param.Name)
					schemas = append(schemas, SchemaLocation{
						Name:     name,
						Path:     path,
						Method:   method,
						Location: "parameter",
						Schema:   param.Schema,
					})
				}
			}
		}
	}

	// Sort for deterministic output
	sort.Slice(schemas, func(i, j int) bool {
		return schemas[i].Name < schemas[j].Name
	})

	return schemas
}

// DiffSpecs compares two OpenAPI specs and returns changes for each schema.
func DiffSpecs(old, new *Spec) map[string]*schema.DiffResult {
	oldSchemas := extractSchemaMap(old)
	newSchemas := extractSchemaMap(new)

	results := make(map[string]*schema.DiffResult)

	// Compare matching schemas
	for name, oldSchema := range oldSchemas {
		if newSchema, exists := newSchemas[name]; exists {
			result := schema.Diff(oldSchema, newSchema)
			if result.Summary.Total > 0 {
				results[name] = result
			}
		} else {
			// Schema removed entirely
			results[name] = schema.NewDiffResult([]schema.Change{
				{
					Path:     name,
					Type:     schema.ChangePropertyRemoved,
					Severity: schema.SeverityBreaking,
					Message:  fmt.Sprintf("Schema %q removed from spec", name),
				},
			})
		}
	}

	// Check for new schemas (informational)
	for name := range newSchemas {
		if _, exists := oldSchemas[name]; !exists {
			results[name] = schema.NewDiffResult([]schema.Change{
				{
					Path:     name,
					Type:     schema.ChangePropertyAdded,
					Severity: schema.SeverityInfo,
					Message:  fmt.Sprintf("Schema %q added to spec", name),
				},
			})
		}
	}

	return results
}

func extractSchemaMap(spec *Spec) map[string]*schema.Schema {
	schemas := make(map[string]*schema.Schema)
	for _, loc := range ExtractSchemas(spec) {
		schemas[loc.Name] = loc.Schema
	}
	return schemas
}
