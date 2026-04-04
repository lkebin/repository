package parser

import (
	"testing"
)

func TestNewPartTree(t *testing.T) {
	tests := []struct {
		name             string
		source           string
		wantNodes        int
		wantChildren     []int
		wantProperties   [][]string
		wantTypes        [][]PartType
		wantSubject      func(*Subject) bool
		wantOrderBy      bool
		wantOrderByCount int
	}{
		// Basic And
		{
			name:           "single property",
			source:         "FindByName",
			wantNodes:      1,
			wantChildren:   []int{1},
			wantProperties: [][]string{{"Name"}},
			wantTypes:      [][]PartType{{SimpleProperty}},
		},
		{
			name:           "two properties with And",
			source:         "FindByIdAndName",
			wantNodes:      1,
			wantChildren:   []int{2},
			wantProperties: [][]string{{"Id", "Name"}},
			wantTypes:      [][]PartType{{SimpleProperty, SimpleProperty}},
		},
		{
			name:           "three properties with And",
			source:         "FindByNameAndBirthdayAndAge",
			wantNodes:      1,
			wantChildren:   []int{3},
			wantProperties: [][]string{{"Name", "Birthday", "Age"}},
			wantTypes:      [][]PartType{{SimpleProperty, SimpleProperty, SimpleProperty}},
		},

		// Or
		{
			name:           "Or clause",
			source:         "FindByNameOrBirthday",
			wantNodes:      2,
			wantChildren:   []int{1, 1},
			wantProperties: [][]string{{"Name"}, {"Birthday"}},
			wantTypes:      [][]PartType{{SimpleProperty}, {SimpleProperty}},
		},
		{
			name:           "mixed And and Or",
			source:         "FindByNameAndBirthdayOrAge",
			wantNodes:      2,
			wantChildren:   []int{2, 1},
			wantProperties: [][]string{{"Name", "Birthday"}, {"Age"}},
			wantTypes:      [][]PartType{{SimpleProperty, SimpleProperty}, {SimpleProperty}},
		},

		// Zero-arg operators
		{
			name:           "IsNull",
			source:         "FindByNameIsNull",
			wantNodes:      1,
			wantChildren:   []int{1},
			wantProperties: [][]string{{"Name"}},
			wantTypes:      [][]PartType{{IsNull}},
		},
		{
			name:           "IsNotNull",
			source:         "FindByNameIsNotNull",
			wantNodes:      1,
			wantChildren:   []int{1},
			wantProperties: [][]string{{"Name"}},
			wantTypes:      [][]PartType{{IsNotNull}},
		},
		{
			name:           "IsEmpty",
			source:         "FindByNameIsEmpty",
			wantNodes:      1,
			wantChildren:   []int{1},
			wantProperties: [][]string{{"Name"}},
			wantTypes:      [][]PartType{{IsEmpty}},
		},
		{
			name:           "IsNotEmpty",
			source:         "FindByNameIsNotEmpty",
			wantNodes:      1,
			wantChildren:   []int{1},
			wantProperties: [][]string{{"Name"}},
			wantTypes:      [][]PartType{{IsNotEmpty}},
		},

		// Mixed zero-arg with normal operators
		{
			name:           "IsNull And SimpleProperty",
			source:         "FindByNameIsNullAndBirthday",
			wantNodes:      1,
			wantChildren:   []int{2},
			wantProperties: [][]string{{"Name", "Birthday"}},
			wantTypes:      [][]PartType{{IsNull, SimpleProperty}},
		},
		{
			name:           "IsNull Or SimpleProperty",
			source:         "FindByNameIsNullOrBirthday",
			wantNodes:      2,
			wantChildren:   []int{1, 1},
			wantProperties: [][]string{{"Name"}, {"Birthday"}},
			wantTypes:      [][]PartType{{IsNull}, {SimpleProperty}},
		},
		{
			name:           "both sides IsNull Or",
			source:         "FindByNameIsNullOrBirthdayIsNull",
			wantNodes:      2,
			wantChildren:   []int{1, 1},
			wantProperties: [][]string{{"Name"}, {"Birthday"}},
			wantTypes:      [][]PartType{{IsNull}, {IsNull}},
		},

		// Parameterized operators
		{
			name:           "Between",
			source:         "FindByBirthdayBetween",
			wantNodes:      1,
			wantChildren:   []int{1},
			wantProperties: [][]string{{"Birthday"}},
			wantTypes:      [][]PartType{{Between}},
		},
		{
			name:           "LessThan",
			source:         "FindByAgeLessThan",
			wantNodes:      1,
			wantChildren:   []int{1},
			wantProperties: [][]string{{"Age"}},
			wantTypes:      [][]PartType{{LessThan}},
		},
		{
			name:           "GreaterThanEqual",
			source:         "FindByAgeGreaterThanEqual",
			wantNodes:      1,
			wantChildren:   []int{1},
			wantProperties: [][]string{{"Age"}},
			wantTypes:      [][]PartType{{GreaterThanEqual}},
		},
		{
			name:           "Like",
			source:         "FindByNameLike",
			wantNodes:      1,
			wantChildren:   []int{1},
			wantProperties: [][]string{{"Name"}},
			wantTypes:      [][]PartType{{Like}},
		},
		{
			name:           "In",
			source:         "FindByNameIn",
			wantNodes:      1,
			wantChildren:   []int{1},
			wantProperties: [][]string{{"Name"}},
			wantTypes:      [][]PartType{{In}},
		},
		{
			name:           "NotIn",
			source:         "FindByNameNotIn",
			wantNodes:      1,
			wantChildren:   []int{1},
			wantProperties: [][]string{{"Name"}},
			wantTypes:      [][]PartType{{NotIn}},
		},
		{
			name:           "NegatingSimpleProperty",
			source:         "FindByNameIsNot",
			wantNodes:      1,
			wantChildren:   []int{1},
			wantProperties: [][]string{{"Name"}},
			wantTypes:      [][]PartType{{NegatingSimpleProperty}},
		},

		// OrderBy
		{
			name:             "with OrderBy",
			source:           "FindByNameOrderByNameAsc",
			wantNodes:        1,
			wantChildren:     []int{1},
			wantProperties:   [][]string{{"Name"}},
			wantTypes:        [][]PartType{{SimpleProperty}},
			wantOrderBy:      true,
			wantOrderByCount: 1,
		},
		{
			name:             "with multiple OrderBy",
			source:           "FindByNameOrderByNameAscBirthdayDesc",
			wantNodes:        1,
			wantChildren:     []int{1},
			wantProperties:   [][]string{{"Name"}},
			wantTypes:        [][]PartType{{SimpleProperty}},
			wantOrderBy:      true,
			wantOrderByCount: 2,
		},

		// Subject types
		{
			name:         "Count subject",
			source:       "CountByName",
			wantNodes:    1,
			wantChildren: []int{1},
			wantSubject:  func(s *Subject) bool { return s.IsCount },
		},
		{
			name:         "Exists subject",
			source:       "ExistsByName",
			wantNodes:    1,
			wantChildren: []int{1},
			wantSubject:  func(s *Subject) bool { return s.IsExists },
		},
		{
			name:         "Delete subject",
			source:       "DeleteByName",
			wantNodes:    1,
			wantChildren: []int{1},
			wantSubject:  func(s *Subject) bool { return s.IsDelete },
		},
		{
			name:         "Distinct",
			source:       "FindDistinctByName",
			wantNodes:    1,
			wantChildren: []int{1},
			wantSubject:  func(s *Subject) bool { return s.IsDistinct },
		},
		{
			name:         "Limiting First",
			source:       "FindFirstByName",
			wantNodes:    1,
			wantChildren: []int{1},
			wantSubject:  func(s *Subject) bool { return s.IsLimiting && s.MaxResults == 1 },
		},
		{
			name:         "Limiting Top10",
			source:       "FindTop10ByName",
			wantNodes:    1,
			wantChildren: []int{1},
			wantSubject:  func(s *Subject) bool { return s.IsLimiting && s.MaxResults == 10 },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt, err := NewPartTree(tt.source)
			if err != nil {
				t.Fatalf("NewPartTree(%q) error: %v", tt.source, err)
			}

			if pt.Subject == nil {
				t.Fatal("Subject should not be nil")
			}

			if len(pt.Predicate.Nodes) != tt.wantNodes {
				t.Fatalf("expected %d nodes, got %d", tt.wantNodes, len(pt.Predicate.Nodes))
			}

			for i, node := range pt.Predicate.Nodes {
				if len(node.Children) != tt.wantChildren[i] {
					t.Errorf("node[%d]: expected %d children, got %d", i, tt.wantChildren[i], len(node.Children))
				}

				if tt.wantProperties != nil {
					for j, child := range node.Children {
						if child.Property != tt.wantProperties[i][j] {
							t.Errorf("node[%d].child[%d]: expected property %q, got %q", i, j, tt.wantProperties[i][j], child.Property)
						}
					}
				}

				if tt.wantTypes != nil {
					for j, child := range node.Children {
						if child.Type.Name != tt.wantTypes[i][j].Name {
							t.Errorf("node[%d].child[%d]: expected type %v, got %v", i, j, tt.wantTypes[i][j], child.Type)
						}
					}
				}
			}

			if tt.wantSubject != nil && !tt.wantSubject(pt.Subject) {
				t.Errorf("subject check failed for %q", tt.source)
			}

			if tt.wantOrderBy {
				if pt.Predicate.OrderBySource == nil {
					t.Error("expected OrderBySource to be non-nil")
				} else if len(pt.Predicate.OrderBySource.Orders) != tt.wantOrderByCount {
					t.Errorf("expected %d orders, got %d", tt.wantOrderByCount, len(pt.Predicate.OrderBySource.Orders))
				}
			}
		})
	}
}

func TestNewPartTreeNoPrefix(t *testing.T) {
	// Source doesn't match prefixTemplate but is a valid predicate.
	pt, err := NewPartTree("Name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pt.Predicate.Nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(pt.Predicate.Nodes))
	}
}

func TestNewPartTreeErrorNoPrefix(t *testing.T) {
	// Source doesn't match prefixTemplate; NewPredicate gets the full string
	// which contains multiple OrderBy, triggering an error.
	_, err := NewPartTree("NameOrderByAscOrderByDesc")
	if err == nil {
		t.Error("expected error for multiple OrderBy in non-prefix source")
	}
}

func TestNewPartTreeErrorWithPrefix(t *testing.T) {
	// Source matches prefixTemplate; predicate portion has multiple OrderBy.
	_, err := NewPartTree("FindByNameOrderByAscOrderByDesc")
	if err == nil {
		t.Error("expected error for multiple OrderBy in predicate after prefix")
	}
}
