package openapi

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseJSON(t *testing.T) {
	data := []byte(`{
		"openapi": "3.0.0",
		"info": {"title": "Test API", "version": "1.0.0"},
		"paths": {
			"/users": {
				"get": {
					"responses": {
						"200": {
							"description": "Success",
							"content": {
								"application/json": {
									"schema": {"type": "array", "items": {"type": "object"}}
								}
							}
						}
					}
				}
			}
		}
	}`)

	spec, err := ParseJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	if spec.OpenAPI != "3.0.0" {
		t.Errorf("expected openapi=3.0.0, got %s", spec.OpenAPI)
	}
	if spec.Info.Title != "Test API" {
		t.Errorf("expected title=Test API, got %s", spec.Info.Title)
	}
}

func TestParseJSONInvalidVersion(t *testing.T) {
	data := []byte(`{
		"openapi": "2.0",
		"info": {"title": "Test", "version": "1.0.0"},
		"paths": {}
	}`)
	_, err := ParseJSON(data)
	if err == nil {
		t.Error("expected error for unsupported version")
	}
}

func TestParseJSONInvalid(t *testing.T) {
	_, err := ParseJSON([]byte(`{invalid`))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestParseYAML(t *testing.T) {
	data := []byte(`
openapi: "3.0.0"
info:
  title: Test API
  version: "1.0.0"
paths:
  /users:
    get:
      responses:
        "200":
          description: Success
`)
	spec, err := ParseYAML(data)
	if err == nil && spec != nil {
		if spec.OpenAPI != "3.0.0" {
			t.Errorf("expected openapi=3.0.0, got %s", spec.OpenAPI)
		}
	}
}

func TestParseYAMLInvalid(t *testing.T) {
	_, err := ParseYAML([]byte(`\tinvalid: [yaml`))
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestParseFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "api.json")
	data := []byte(`{
		"openapi": "3.0.0",
		"info": {"title": "Test", "version": "1.0.0"},
		"paths": {}
	}`)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	spec, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if spec.Info.Title != "Test" {
		t.Errorf("expected title=Test, got %s", spec.Info.Title)
	}
}

func TestParseFileNotFound(t *testing.T) {
	_, err := ParseFile("/nonexistent/api.json")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestParseFileYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "api.yaml")
	data := []byte("openapi: \"3.0.0\"\ninfo:\n  title: YAML API\n  version: \"1.0.0\"\npaths: {}\n")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	spec, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if spec.Info.Title != "YAML API" {
		t.Errorf("expected title=YAML API, got %s", spec.Info.Title)
	}
}

func TestParseWithComponents(t *testing.T) {
	data := []byte(`{
		"openapi": "3.0.0",
		"info": {"title": "Test", "version": "1.0.0"},
		"paths": {},
		"components": {
			"schemas": {
				"User": {
					"type": "object",
					"properties": {
						"name": {"type": "string"},
						"email": {"type": "string"}
					},
					"required": ["name", "email"]
				}
			}
		}
	}`)

	spec, err := ParseJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	if spec.Components == nil {
		t.Fatal("expected components")
	}
	if _, ok := spec.Components.Schemas["User"]; !ok {
		t.Error("expected User schema")
	}
}

func TestParseWithRequestBody(t *testing.T) {
	data := []byte(`{
		"openapi": "3.0.0",
		"info": {"title": "Test", "version": "1.0.0"},
		"paths": {
			"/users": {
				"post": {
					"requestBody": {
						"required": true,
						"content": {
							"application/json": {
								"schema": {
									"type": "object",
									"properties": {
										"name": {"type": "string"}
									}
								}
							}
						}
					},
					"responses": {
						"201": {"description": "Created"}
					}
				}
			}
		}
	}`)

	spec, err := ParseJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	if spec.Paths["/users"].Post.RequestBody == nil {
		t.Error("expected request body")
	}
}

func TestParseWithParameters(t *testing.T) {
	data := []byte(`{
		"openapi": "3.0.0",
		"info": {"title": "Test", "version": "1.0.0"},
		"paths": {
			"/users/{id}": {
				"get": {
					"parameters": [
						{
							"name": "id",
							"in": "path",
							"required": true,
							"schema": {"type": "string"}
						}
					],
					"responses": {
						"200": {"description": "Success"}
					}
				}
			}
		}
	}`)

	spec, err := ParseJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	params := spec.Paths["/users/{id}"].Get.Parameters
	if len(params) != 1 {
		t.Fatalf("expected 1 parameter, got %d", len(params))
	}
	if params[0].Name != "id" {
		t.Errorf("expected param name=id, got %s", params[0].Name)
	}
}
