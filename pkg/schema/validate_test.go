package schema

import "testing"

func TestValidateSchemaValid(t *testing.T) {
	s := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"name": {Type: "string", MinLength: intPtr(1), MaxLength: intPtr(100)},
			"age":  {Type: "integer", Minimum: floatPtr(0), Maximum: floatPtr(150)},
		},
	}
	errs := ValidateSchema(s)
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
}

func TestValidateSchemaInvalidType(t *testing.T) {
	s := &Schema{Type: "invalid"}
	errs := ValidateSchema(s)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if errs[0].Field != "type" {
		t.Errorf("expected field=type, got %s", errs[0].Field)
	}
}

func TestValidateSchemaMinMaxLength(t *testing.T) {
	s := &Schema{Type: "string", MinLength: intPtr(10), MaxLength: intPtr(5)}
	errs := ValidateSchema(s)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestValidateSchemaMinMaxNum(t *testing.T) {
	s := &Schema{Type: "number", Minimum: floatPtr(100), Maximum: floatPtr(10)}
	errs := ValidateSchema(s)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestValidateSchemaMinMaxItems(t *testing.T) {
	s := &Schema{Type: "array", MinItems: intPtr(10), MaxItems: intPtr(5)}
	errs := ValidateSchema(s)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestValidateSchemaNil(t *testing.T) {
	errs := ValidateSchema(nil)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error for nil, got %d", len(errs))
	}
}

func TestValidateSchemaNestedInvalid(t *testing.T) {
	s := &Schema{
		Type: "object",
		Properties: map[string]*Schema{
			"bad": {Type: "string", MinLength: intPtr(10), MaxLength: intPtr(5)},
		},
	}
	errs := ValidateSchema(s)
	if len(errs) == 0 {
		t.Error("expected errors for nested invalid schema")
	}
}

func TestValidateSchemaItemsInvalid(t *testing.T) {
	s := &Schema{
		Type:  "array",
		Items: &Schema{Type: "notreal"},
	}
	errs := ValidateSchema(s)
	if len(errs) == 0 {
		t.Error("expected errors for invalid items schema")
	}
}

func TestValidationErrorString(t *testing.T) {
	e := ValidationError{Field: "type", Message: "invalid"}
	if e.Error() != "type: invalid" {
		t.Errorf("unexpected error string: %s", e.Error())
	}
}

func TestValidateSchemaAllValidTypes(t *testing.T) {
	validTypes := []string{"string", "number", "integer", "boolean", "object", "array", "null"}
	for _, typ := range validTypes {
		s := &Schema{Type: typ}
		errs := ValidateSchema(s)
		if len(errs) != 0 {
			t.Errorf("expected no errors for type=%s, got %v", typ, errs)
		}
	}
}

func TestValidateSchemaNoType(t *testing.T) {
	s := &Schema{Properties: map[string]*Schema{"x": {Type: "string"}}}
	errs := ValidateSchema(s)
	if len(errs) != 0 {
		t.Error("schema without type should be valid")
	}
}
