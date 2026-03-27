package openapi

import (
	"testing"

	"github.com/JSLEEKR/schemadiff/pkg/schema"
)

func TestDiffPathsPathRemoved(t *testing.T) {
	old := &Spec{
		OpenAPI: "3.0.0",
		Paths: map[string]*PathItem{
			"/users": {Get: &Operation{}},
			"/posts": {Get: &Operation{}},
		},
	}
	new := &Spec{
		OpenAPI: "3.0.0",
		Paths: map[string]*PathItem{
			"/users": {Get: &Operation{}},
		},
	}

	changes := DiffPaths(old, new)
	found := false
	for _, c := range changes {
		if c.Type == PathRemoved && c.Path == "/posts" {
			found = true
			if c.Severity != schema.SeverityBreaking {
				t.Error("expected breaking for path removal")
			}
		}
	}
	if !found {
		t.Error("expected path-removed change")
	}
}

func TestDiffPathsPathAdded(t *testing.T) {
	old := &Spec{Paths: map[string]*PathItem{"/users": {Get: &Operation{}}}}
	new := &Spec{Paths: map[string]*PathItem{
		"/users": {Get: &Operation{}},
		"/posts": {Get: &Operation{}},
	}}

	changes := DiffPaths(old, new)
	found := false
	for _, c := range changes {
		if c.Type == PathAdded {
			found = true
		}
	}
	if !found {
		t.Error("expected path-added change")
	}
}

func TestDiffPathsMethodRemoved(t *testing.T) {
	old := &Spec{Paths: map[string]*PathItem{
		"/users": {Get: &Operation{}, Post: &Operation{}},
	}}
	new := &Spec{Paths: map[string]*PathItem{
		"/users": {Get: &Operation{}},
	}}

	changes := DiffPaths(old, new)
	found := false
	for _, c := range changes {
		if c.Type == MethodRemoved && c.Method == "POST" {
			found = true
		}
	}
	if !found {
		t.Error("expected method-removed change")
	}
}

func TestDiffPathsMethodAdded(t *testing.T) {
	old := &Spec{Paths: map[string]*PathItem{
		"/users": {Get: &Operation{}},
	}}
	new := &Spec{Paths: map[string]*PathItem{
		"/users": {Get: &Operation{}, Put: &Operation{}},
	}}

	changes := DiffPaths(old, new)
	found := false
	for _, c := range changes {
		if c.Type == MethodAdded {
			found = true
		}
	}
	if !found {
		t.Error("expected method-added change")
	}
}

func TestDiffPathsParameterRemoved(t *testing.T) {
	old := &Spec{Paths: map[string]*PathItem{
		"/users": {Get: &Operation{
			Parameters: []Parameter{
				{Name: "page", In: "query"},
				{Name: "limit", In: "query"},
			},
		}},
	}}
	new := &Spec{Paths: map[string]*PathItem{
		"/users": {Get: &Operation{
			Parameters: []Parameter{
				{Name: "page", In: "query"},
			},
		}},
	}}

	changes := DiffPaths(old, new)
	found := false
	for _, c := range changes {
		if c.Type == ParameterRemoved {
			found = true
		}
	}
	if !found {
		t.Error("expected parameter-removed change")
	}
}

func TestDiffPathsRequiredParamAdded(t *testing.T) {
	old := &Spec{Paths: map[string]*PathItem{
		"/users": {Get: &Operation{
			Parameters: []Parameter{
				{Name: "filter", In: "query", Required: false},
			},
		}},
	}}
	new := &Spec{Paths: map[string]*PathItem{
		"/users": {Get: &Operation{
			Parameters: []Parameter{
				{Name: "filter", In: "query", Required: true},
			},
		}},
	}}

	changes := DiffPaths(old, new)
	found := false
	for _, c := range changes {
		if c.Type == RequiredParamAdded {
			found = true
			if c.Severity != schema.SeverityBreaking {
				t.Error("expected breaking")
			}
		}
	}
	if !found {
		t.Error("expected required-parameter-added change")
	}
}

func TestDiffPathsParameterTypeChanged(t *testing.T) {
	old := &Spec{Paths: map[string]*PathItem{
		"/users": {Get: &Operation{
			Parameters: []Parameter{
				{Name: "id", In: "path", Schema: &schema.Schema{Type: "string"}},
			},
		}},
	}}
	new := &Spec{Paths: map[string]*PathItem{
		"/users": {Get: &Operation{
			Parameters: []Parameter{
				{Name: "id", In: "path", Schema: &schema.Schema{Type: "integer"}},
			},
		}},
	}}

	changes := DiffPaths(old, new)
	found := false
	for _, c := range changes {
		if c.Type == ParameterTypeChanged {
			found = true
		}
	}
	if !found {
		t.Error("expected parameter-type-changed change")
	}
}

func TestDiffPathsNewRequiredParam(t *testing.T) {
	old := &Spec{Paths: map[string]*PathItem{
		"/users": {Get: &Operation{
			Parameters: []Parameter{},
		}},
	}}
	new := &Spec{Paths: map[string]*PathItem{
		"/users": {Get: &Operation{
			Parameters: []Parameter{
				{Name: "apiKey", In: "header", Required: true},
			},
		}},
	}}

	changes := DiffPaths(old, new)
	found := false
	for _, c := range changes {
		if c.Type == ParameterAdded && c.Severity == schema.SeverityBreaking {
			found = true
		}
	}
	if !found {
		t.Error("expected breaking for new required parameter")
	}
}

func TestDiffPathsResponseRemoved(t *testing.T) {
	old := &Spec{Paths: map[string]*PathItem{
		"/users": {Get: &Operation{
			Responses: map[string]*Response{
				"200": {Description: "OK"},
				"404": {Description: "Not Found"},
			},
		}},
	}}
	new := &Spec{Paths: map[string]*PathItem{
		"/users": {Get: &Operation{
			Responses: map[string]*Response{
				"200": {Description: "OK"},
			},
		}},
	}}

	changes := DiffPaths(old, new)
	found := false
	for _, c := range changes {
		if c.Type == ResponseRemoved {
			found = true
		}
	}
	if !found {
		t.Error("expected response-removed change")
	}
}

func TestDiffPathsNoChanges(t *testing.T) {
	spec := &Spec{Paths: map[string]*PathItem{
		"/users": {Get: &Operation{
			Parameters: []Parameter{{Name: "id", In: "path"}},
			Responses:  map[string]*Response{"200": {Description: "OK"}},
		}},
	}}

	changes := DiffPaths(spec, spec)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestDiffPathsEmptySpecs(t *testing.T) {
	old := &Spec{Paths: map[string]*PathItem{}}
	new := &Spec{Paths: map[string]*PathItem{}}

	changes := DiffPaths(old, new)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}
