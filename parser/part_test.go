package parser

import (
	"reflect"
	"testing"
)

func TestNewPart(t *testing.T) {
	tests := []struct {
		name           string
		source         string
		wantProperty   string
		wantType       PartType
		wantIgnoreCase IgnoreCaseType
	}{
		// SimpleProperty and variants
		{"Is keyword", "NameIs", "Name", SimpleProperty, Never},
		{"Equals keyword", "NameEquals", "Name", SimpleProperty, Never},
		{"plain property falls back to SimpleProperty", "PlainProperty", "PlainProperty", SimpleProperty, Never},

		// NegatingSimpleProperty
		{"IsNot keyword", "NameIsNot", "Name", NegatingSimpleProperty, Never},
		{"Not keyword", "NameNot", "Name", NegatingSimpleProperty, Never},

		// Null/NotNull
		{"IsNull", "NameIsNull", "Name", IsNull, Never},
		{"Null keyword variant", "NameNull", "Name", IsNull, Never},
		{"IsNotNull", "NameIsNotNull", "Name", IsNotNull, Never},
		{"NotNull keyword variant", "NameNotNull", "Name", IsNotNull, Never},

		// Empty/NotEmpty
		{"IsEmpty", "NameIsEmpty", "Name", IsEmpty, Never},
		{"Empty keyword variant", "NameEmpty", "Name", IsEmpty, Never},
		{"IsNotEmpty", "NameIsNotEmpty", "Name", IsNotEmpty, Never},
		{"NotEmpty keyword variant", "NameNotEmpty", "Name", IsNotEmpty, Never},

		// Comparison
		{"Between", "AgeBetween", "Age", Between, Never},
		{"IsBetween keyword variant", "AgeIsBetween", "Age", Between, Never},
		{"LessThan", "AgeLessThan", "Age", LessThan, Never},
		{"LessThanEqual", "AgeLessThanEqual", "Age", LessThanEqual, Never},
		{"GreaterThan", "AgeGreaterThan", "Age", GreaterThan, Never},
		{"GreaterThanEqual", "AgeGreaterThanEqual", "Age", GreaterThanEqual, Never},
		{"Before", "DateBefore", "Date", Before, Never},
		{"After", "DateAfter", "Date", After, Never},

		// String matching
		{"Like", "NameLike", "Name", Like, Never},
		{"NotLike", "NameNotLike", "Name", NotLike, Never},
		{"StartingWith", "NameStartingWith", "Name", StartingWith, Never},
		{"EndingWith", "NameEndingWith", "Name", EndingWith, Never},
		{"Containing", "NameContaining", "Name", Containing, Never},
		{"NotContaining", "NameNotContaining", "Name", NotContaining, Never},

		// Collection
		{"In", "NameIn", "Name", In, Never},
		{"NotIn", "NameNotIn", "Name", NotIn, Never},

		// Boolean
		{"True", "ActiveTrue", "Active", True, Never},
		{"IsTrue keyword variant", "ActiveIsTrue", "Active", True, Never},
		{"False", "ActiveFalse", "Active", False, Never},
		{"IsFalse keyword variant", "ActiveIsFalse", "Active", False, Never},

		// Exists
		{"Exists", "NameExists", "Name", Exists, Never},

		// IgnoreCase
		{"IgnoreCase", "NameIgnoreCase", "Name", SimpleProperty, Always},
		{"IgnoringCase", "NameIgnoringCase", "Name", SimpleProperty, Always},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			part := NewPart(tt.source, false)

			if part.Property != tt.wantProperty {
				t.Errorf("Property: expected %q, got %q", tt.wantProperty, part.Property)
			}

			if !reflect.DeepEqual(part.Type, tt.wantType) {
				t.Errorf("Type: expected %v, got %v", tt.wantType, part.Type)
			}

			if part.IgnoreCase != tt.wantIgnoreCase {
				t.Errorf("IgnoreCase: expected %v, got %v", tt.wantIgnoreCase, part.IgnoreCase)
			}
		})
	}
}

func TestNewPartAlwaysIgnoreCase(t *testing.T) {
	part := NewPart("Name", true)
	if part.IgnoreCase != WhenPossible {
		t.Errorf("expected WhenPossible when isAlwaysIgnoreCase=true, got %v", part.IgnoreCase)
	}

	part2 := NewPart("NameIgnoreCase", true)
	if part2.IgnoreCase != Always {
		t.Errorf("explicit IgnoreCase should remain Always even with isAlwaysIgnoreCase=true, got %v", part2.IgnoreCase)
	}
}
