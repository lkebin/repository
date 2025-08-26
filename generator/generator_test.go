package generator

import (
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"UserRepository", "user_repository"},
		{"UserUuidRepository", "user_uuid_repository"},
		{"My2025Database", "my_2025_database"},
		{"My2025DatabaseA2", "my_2025_database_a_2"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := ToSnakeCase(test.input)
			if result != test.expected {
				t.Errorf("expected %s, got %s", test.expected, result)
			}
		})
	}
}
