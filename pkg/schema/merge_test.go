package schema

import "testing"

func TestMergeAllOfBasic(t *testing.T) {
	s := &Schema{
		AllOf: []*Schema{
			{
				Type: "object",
				Properties: map[string]*Schema{
					"name": {Type: "string"},
				},
				Required: []string{"name"},
			},
			{
				Properties: map[string]*Schema{
					"age": {Type: "integer"},
				},
				Required: []string{"age"},
			},
		},
	}

	merged := MergeAllOf(s)
	if merged.Type != "object" {
		t.Errorf("expected type=object, got %v", merged.Type)
	}
	if len(merged.Properties) != 2 {
		t.Errorf("expected 2 properties, got %d", len(merged.Properties))
	}
	if len(merged.Required) != 2 {
		t.Errorf("expected 2 required, got %d", len(merged.Required))
	}
	if merged.AllOf != nil {
		t.Error("expected allOf to be nil after merge")
	}
}

func TestMergeAllOfNested(t *testing.T) {
	s := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"address": {
				AllOf: []*Schema{
					{
						Type: "object",
						Properties: map[string]*Schema{
							"street": {Type: "string"},
						},
					},
					{
						Properties: map[string]*Schema{
							"city": {Type: "string"},
						},
					},
				},
			},
		},
	}

	merged := MergeAllOf(s)
	addr := merged.Properties["address"]
	if len(addr.Properties) != 2 {
		t.Errorf("expected 2 address properties, got %d", len(addr.Properties))
	}
}

func TestMergeAllOfConstraints(t *testing.T) {
	s := &Schema{
		AllOf: []*Schema{
			{Type: "string", MinLength: intPtr(1), MaxLength: intPtr(100)},
			{MinLength: intPtr(5), MaxLength: intPtr(50)},
		},
	}

	merged := MergeAllOf(s)
	if *merged.MinLength != 5 {
		t.Errorf("expected minLength=5 (stricter), got %d", *merged.MinLength)
	}
	if *merged.MaxLength != 50 {
		t.Errorf("expected maxLength=50 (stricter), got %d", *merged.MaxLength)
	}
}

func TestMergeAllOfNil(t *testing.T) {
	if MergeAllOf(nil) != nil {
		t.Error("expected nil for nil input")
	}
}

func TestMergeAllOfNoAllOf(t *testing.T) {
	s := &Schema{Type: "string"}
	merged := MergeAllOf(s)
	if merged.Type != "string" {
		t.Error("expected same schema when no allOf")
	}
}

func TestMergeAllOfWithItems(t *testing.T) {
	s := &Schema{
		AllOf: []*Schema{
			{Type: "array"},
			{Items: &Schema{Type: "string"}},
		},
	}
	merged := MergeAllOf(s)
	if merged.Items == nil {
		t.Error("expected items after merge")
	}
	if merged.Items.Type != "string" {
		t.Errorf("expected items type=string, got %v", merged.Items.Type)
	}
}

func TestMergeAllOfDuplicateRequired(t *testing.T) {
	s := &Schema{
		AllOf: []*Schema{
			{Required: []string{"name", "email"}},
			{Required: []string{"email", "age"}},
		},
	}
	merged := MergeAllOf(s)
	if len(merged.Required) != 3 {
		t.Errorf("expected 3 unique required, got %d: %v", len(merged.Required), merged.Required)
	}
}

func TestMergeAllOfMinimumMaximum(t *testing.T) {
	s := &Schema{
		AllOf: []*Schema{
			{Minimum: floatPtr(0), Maximum: floatPtr(100)},
			{Minimum: floatPtr(10), Maximum: floatPtr(50)},
		},
	}
	merged := MergeAllOf(s)
	if *merged.Minimum != 10 {
		t.Errorf("expected minimum=10, got %v", *merged.Minimum)
	}
	if *merged.Maximum != 50 {
		t.Errorf("expected maximum=50, got %v", *merged.Maximum)
	}
}

func TestMergeAllOfMinMaxItems(t *testing.T) {
	s := &Schema{
		AllOf: []*Schema{
			{MinItems: intPtr(1), MaxItems: intPtr(100)},
			{MinItems: intPtr(5), MaxItems: intPtr(20)},
		},
	}
	merged := MergeAllOf(s)
	if *merged.MinItems != 5 {
		t.Errorf("expected minItems=5, got %d", *merged.MinItems)
	}
	if *merged.MaxItems != 20 {
		t.Errorf("expected maxItems=20, got %d", *merged.MaxItems)
	}
}
