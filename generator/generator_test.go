package generator

import (
	"testing"

	"github.com/lkebin/repository/parser"
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

func TestParseOperator(t *testing.T) {
	tests := []struct {
		name     string
		partType parser.PartType
		want     string
		wantErr  bool
	}{
		{"Between", parser.Between, " BETWEEN ? AND ?", false},
		{"IsNotNull", parser.IsNotNull, " IS NOT NULL", false},
		{"IsNull", parser.IsNull, " IS NULL", false},
		{"LessThan", parser.LessThan, " < ?", false},
		{"LessThanEqual", parser.LessThanEqual, " <= ?", false},
		{"GreaterThan", parser.GreaterThan, " > ?", false},
		{"GreaterThanEqual", parser.GreaterThanEqual, " >= ?", false},
		{"Before", parser.Before, " < ?", false},
		{"After", parser.After, " > ?", false},
		{"NotLike", parser.NotLike, " NOT LIKE ?", false},
		{"Like", parser.Like, " LIKE ?", false},
		{"StartingWith", parser.StartingWith, " LIKE ?", false},
		{"EndingWith", parser.EndingWith, " LIKE ?", false},
		{"IsNotEmpty", parser.IsNotEmpty, " IS NOT NULL", false},
		{"IsEmpty", parser.IsEmpty, " IS NULL", false},
		{"NotContaining", parser.NotContaining, " NOT LIKE ?", false},
		{"Containing", parser.Containing, " LIKE ?", false},
		{"NotIn", parser.NotIn, " NOT IN (?)", false},
		{"In", parser.In, " IN (?)", false},
		{"NegatingSimpleProperty", parser.NegatingSimpleProperty, " != ?", false},
		{"SimpleProperty", parser.SimpleProperty, " = ?", false},
		{"unsupported operator", parser.PartType{NumberOfArguments: 99, Keywords: []string{"Unknown"}}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseOperator(tt.partType)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestParseTag(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		wantName string
		wantOpts TagOptions
	}{
		{"simple name", "name", "name", ""},
		{"name with pk", "id,pk", "id", "pk"},
		{"name with multiple opts", "id,pk,autoincrement", "id", "pk,autoincrement"},
		{"empty tag", "", "", ""},
		{"name with unsafe", "created_at,unsafe", "created_at", "unsafe"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, opts := ParseTag(tt.tag)
			if name != tt.wantName {
				t.Errorf("name: expected %q, got %q", tt.wantName, name)
			}
			if opts != tt.wantOpts {
				t.Errorf("opts: expected %q, got %q", tt.wantOpts, opts)
			}
		})
	}
}

func TestTagOptionsContains(t *testing.T) {
	tests := []struct {
		name   string
		opts   TagOptions
		option string
		want   bool
	}{
		{"contains pk", "pk,autoincrement", "pk", true},
		{"contains autoincrement", "pk,autoincrement", "autoincrement", true},
		{"does not contain unsafe", "pk,autoincrement", "unsafe", false},
		{"single option match", "pk", "pk", true},
		{"single option no match", "pk", "autoincrement", false},
		{"empty options", "", "pk", false},
		{"contains unsafe", "unsafe", "unsafe", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opts.Contains(tt.option)
			if got != tt.want {
				t.Errorf("TagOptions(%q).Contains(%q) = %v, want %v", tt.opts, tt.option, got, tt.want)
			}
		})
	}
}

func TestTagOptionsGet(t *testing.T) {
	tests := []struct {
		name       string
		opts       TagOptions
		optionName string
		want       string
	}{
		{"get table value", "pk,table=users", "table", "users"},
		{"get missing key", "pk,table=users", "missing", ""},
		{"empty options", "", "table", ""},
		{"key without value", "pk,autoincrement", "pk", ""},
		{"multiple key-value pairs", "pk,table=users,charset=utf8", "charset", "utf8"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opts.Get(tt.optionName)
			if got != tt.want {
				t.Errorf("TagOptions(%q).Get(%q) = %q, want %q", tt.opts, tt.optionName, got, tt.want)
			}
		})
	}
}
