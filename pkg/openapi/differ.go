package openapi

import (
	"fmt"

	"github.com/JSLEEKR/schemadiff/pkg/schema"
)

// PathChange represents a change at the OpenAPI path level.
type PathChange struct {
	Path     string
	Method   string
	Type     PathChangeType
	Severity schema.Severity
	Message  string
}

// PathChangeType identifies the kind of path-level change.
type PathChangeType string

const (
	PathRemoved          PathChangeType = "path-removed"
	PathAdded            PathChangeType = "path-added"
	MethodRemoved        PathChangeType = "method-removed"
	MethodAdded          PathChangeType = "method-added"
	ParameterAdded       PathChangeType = "parameter-added"
	ParameterRemoved     PathChangeType = "parameter-removed"
	ParameterTypeChanged PathChangeType = "parameter-type-changed"
	RequiredParamAdded   PathChangeType = "required-parameter-added"
	ResponseRemoved      PathChangeType = "response-removed"
	ResponseAdded        PathChangeType = "response-added"
)

// DiffPaths detects path-level changes between two OpenAPI specs.
func DiffPaths(old, new *Spec) []PathChange {
	var changes []PathChange

	// Paths removed
	for path := range old.Paths {
		if _, exists := new.Paths[path]; !exists {
			changes = append(changes, PathChange{
				Path:     path,
				Type:     PathRemoved,
				Severity: schema.SeverityBreaking,
				Message:  fmt.Sprintf("Path %q removed", path),
			})
		}
	}

	// Paths added
	for path := range new.Paths {
		if _, exists := old.Paths[path]; !exists {
			changes = append(changes, PathChange{
				Path:     path,
				Type:     PathAdded,
				Severity: schema.SeverityInfo,
				Message:  fmt.Sprintf("Path %q added", path),
			})
		}
	}

	// Compare matching paths
	for path, oldItem := range old.Paths {
		newItem, exists := new.Paths[path]
		if !exists {
			continue
		}

		changes = append(changes, diffPathItem(path, oldItem, newItem)...)
	}

	return changes
}

func diffPathItem(path string, old, new *PathItem) []PathChange {
	var changes []PathChange

	ops := map[string]struct {
		old *Operation
		new *Operation
	}{
		"GET":    {old.Get, new.Get},
		"POST":   {old.Post, new.Post},
		"PUT":    {old.Put, new.Put},
		"PATCH":  {old.Patch, new.Patch},
		"DELETE": {old.Delete, new.Delete},
	}

	for method, pair := range ops {
		if pair.old != nil && pair.new == nil {
			changes = append(changes, PathChange{
				Path:     path,
				Method:   method,
				Type:     MethodRemoved,
				Severity: schema.SeverityBreaking,
				Message:  fmt.Sprintf("%s %s removed", method, path),
			})
		} else if pair.old == nil && pair.new != nil {
			changes = append(changes, PathChange{
				Path:     path,
				Method:   method,
				Type:     MethodAdded,
				Severity: schema.SeverityInfo,
				Message:  fmt.Sprintf("%s %s added", method, path),
			})
		} else if pair.old != nil && pair.new != nil {
			changes = append(changes, diffOperation(path, method, pair.old, pair.new)...)
		}
	}

	return changes
}

func diffOperation(path, method string, old, new *Operation) []PathChange {
	var changes []PathChange

	// Compare parameters
	oldParams := paramMap(old.Parameters)
	newParams := paramMap(new.Parameters)

	for name, oldParam := range oldParams {
		newParam, exists := newParams[name]
		if !exists {
			changes = append(changes, PathChange{
				Path:     path,
				Method:   method,
				Type:     ParameterRemoved,
				Severity: schema.SeverityBreaking,
				Message:  fmt.Sprintf("Parameter %q removed from %s %s", name, method, path),
			})
			continue
		}

		// Check required change
		if !oldParam.Required && newParam.Required {
			changes = append(changes, PathChange{
				Path:     path,
				Method:   method,
				Type:     RequiredParamAdded,
				Severity: schema.SeverityBreaking,
				Message:  fmt.Sprintf("Parameter %q made required in %s %s", name, method, path),
			})
		}

		// Check type change
		if oldParam.Schema != nil && newParam.Schema != nil {
			oldType := fmt.Sprintf("%v", oldParam.Schema.Type)
			newType := fmt.Sprintf("%v", newParam.Schema.Type)
			if oldType != newType && oldType != "" && newType != "" {
				changes = append(changes, PathChange{
					Path:     path,
					Method:   method,
					Type:     ParameterTypeChanged,
					Severity: schema.SeverityBreaking,
					Message:  fmt.Sprintf("Parameter %q type changed from %s to %s in %s %s", name, oldType, newType, method, path),
				})
			}
		}
	}

	for name := range newParams {
		if _, exists := oldParams[name]; !exists {
			newParam := newParams[name]
			sev := schema.SeverityInfo
			if newParam.Required {
				sev = schema.SeverityBreaking
			}
			changes = append(changes, PathChange{
				Path:     path,
				Method:   method,
				Type:     ParameterAdded,
				Severity: sev,
				Message:  fmt.Sprintf("Parameter %q added to %s %s", name, method, path),
			})
		}
	}

	// Compare responses
	for code := range old.Responses {
		if _, exists := new.Responses[code]; !exists {
			changes = append(changes, PathChange{
				Path:     path,
				Method:   method,
				Type:     ResponseRemoved,
				Severity: schema.SeverityBreaking,
				Message:  fmt.Sprintf("Response %s removed from %s %s", code, method, path),
			})
		}
	}
	for code := range new.Responses {
		if _, exists := old.Responses[code]; !exists {
			changes = append(changes, PathChange{
				Path:     path,
				Method:   method,
				Type:     ResponseAdded,
				Severity: schema.SeverityInfo,
				Message:  fmt.Sprintf("Response %s added to %s %s", code, method, path),
			})
		}
	}

	return changes
}

func paramMap(params []Parameter) map[string]Parameter {
	m := make(map[string]Parameter, len(params))
	for _, p := range params {
		m[p.Name] = p
	}
	return m
}
