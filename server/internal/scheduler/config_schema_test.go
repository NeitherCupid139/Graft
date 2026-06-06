package scheduler

import "testing"

func TestValidateConfigSchemaRejectsInvertedBounds(t *testing.T) {
	testCases := []struct {
		name   string
		schema string
	}{
		{
			name:   "string length bounds",
			schema: `{"type":"object","properties":{"name":{"type":"string","minLength":5,"maxLength":3}}}`,
		},
		{
			name:   "number bounds",
			schema: `{"type":"object","properties":{"count":{"type":"integer","minimum":10,"maximum":2}}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := ValidateConfigSchema(tc.schema); err == nil {
				t.Fatal("expected invalid schema bounds")
			}
		})
	}
}

func TestValidateConfigJSONEnumKeepsValueTypesStrict(t *testing.T) {
	schema := `{"type":"object","properties":{"mode":{"type":"string","enum":["42"]},"count":{"type":"integer","enum":[42]}}}`
	if err := ValidateConfigJSON(schema, `{"mode":"42","count":42}`); err != nil {
		t.Fatalf("expected matching enum values: %v", err)
	}
	if err := ValidateConfigJSON(schema, `{"mode":42,"count":42}`); err == nil {
		t.Fatal("expected numeric value to fail string enum/type validation")
	}
	if err := ValidateConfigJSON(schema, `{"mode":"42","count":"42"}`); err == nil {
		t.Fatal("expected string value to fail integer enum/type validation")
	}
}
