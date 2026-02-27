package parser

import (
	"reflect"
	"testing"
)

func TestNewPredicate(t *testing.T) {
	tests := []struct {
		name           string
		predicate      string
		wantNodes      int
		wantChildren   []int
		wantProperties [][]string
		wantTypes      [][]PartType
		wantIgnoreCase bool
	}{
		{
			name:           "simple property",
			predicate:      "Name",
			wantNodes:      1,
			wantChildren:   []int{1},
			wantProperties: [][]string{{"Name"}},
			wantTypes:      [][]PartType{{SimpleProperty}},
		},
		{
			name:           "And splits into children",
			predicate:      "NameAndBirthday",
			wantNodes:      1,
			wantChildren:   []int{2},
			wantProperties: [][]string{{"Name", "Birthday"}},
			wantTypes:      [][]PartType{{SimpleProperty, SimpleProperty}},
		},
		{
			name:           "Or splits into nodes",
			predicate:      "NameOrBirthday",
			wantNodes:      2,
			wantChildren:   []int{1, 1},
			wantProperties: [][]string{{"Name"}, {"Birthday"}},
			wantTypes:      [][]PartType{{SimpleProperty}, {SimpleProperty}},
		},
		{
			name:           "And binds tighter than Or",
			predicate:      "NameAndBirthdayOrAge",
			wantNodes:      2,
			wantChildren:   []int{2, 1},
			wantProperties: [][]string{{"Name", "Birthday"}, {"Age"}},
			wantTypes:      [][]PartType{{SimpleProperty, SimpleProperty}, {SimpleProperty}},
		},
		{
			name:           "zero-arg in predicate",
			predicate:      "NameIsNullAndBirthday",
			wantNodes:      1,
			wantChildren:   []int{2},
			wantProperties: [][]string{{"Name", "Birthday"}},
			wantTypes:      [][]PartType{{IsNull, SimpleProperty}},
		},
		{
			name:           "all zero-arg predicates",
			predicate:      "NameIsNullAndBirthdayIsNull",
			wantNodes:      1,
			wantChildren:   []int{2},
			wantProperties: [][]string{{"Name", "Birthday"}},
			wantTypes:      [][]PartType{{IsNull, IsNull}},
		},
		{
			name:           "AllIgnoreCase flag",
			predicate:      "NameAllIgnoreCase",
			wantNodes:      1,
			wantChildren:   []int{1},
			wantIgnoreCase: true,
		},
		{
			name:           "AllIgnoringCase flag",
			predicate:      "NameAllIgnoringCase",
			wantNodes:      1,
			wantChildren:   []int{1},
			wantIgnoreCase: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewPredicate(tt.predicate)
			if err != nil {
				t.Fatalf("NewPredicate(%q) error: %v", tt.predicate, err)
			}

			if len(p.Nodes) != tt.wantNodes {
				t.Fatalf("expected %d nodes, got %d", tt.wantNodes, len(p.Nodes))
			}

			for i, node := range p.Nodes {
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
						if !reflect.DeepEqual(child.Type, tt.wantTypes[i][j]) {
							t.Errorf("node[%d].child[%d]: expected type %v, got %v", i, j, tt.wantTypes[i][j], child.Type)
						}
					}
				}
			}

			if tt.wantIgnoreCase && !p.IsAlwaysIgnoreCase {
				t.Error("expected IsAlwaysIgnoreCase to be true")
			}
		})
	}
}

func TestNewPredicateOrderBy(t *testing.T) {
	tests := []struct {
		name       string
		predicate  string
		wantOrders []Order
	}{
		{
			name:      "single order",
			predicate: "NameOrderByNameAsc",
			wantOrders: []Order{
				{Property: "Name", Direction: "Asc"},
			},
		},
		{
			name:      "multiple orders",
			predicate: "NameOrderByNameAscBirthdayDesc",
			wantOrders: []Order{
				{Property: "Name", Direction: "Asc"},
				{Property: "Birthday", Direction: "Desc"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewPredicate(tt.predicate)
			if err != nil {
				t.Fatalf("NewPredicate(%q) error: %v", tt.predicate, err)
			}

			if p.OrderBySource == nil {
				t.Fatal("expected OrderBySource to be non-nil")
			}

			if len(p.OrderBySource.Orders) != len(tt.wantOrders) {
				t.Fatalf("expected %d orders, got %d", len(tt.wantOrders), len(p.OrderBySource.Orders))
			}

			for i, want := range tt.wantOrders {
				got := p.OrderBySource.Orders[i]
				if got.Property != want.Property {
					t.Errorf("order[%d]: expected property %q, got %q", i, want.Property, got.Property)
				}
				if got.Direction != want.Direction {
					t.Errorf("order[%d]: expected direction %q, got %q", i, want.Direction, got.Direction)
				}
			}
		})
	}
}

func TestNewPredicateMultipleOrderByError(t *testing.T) {
	_, err := NewPredicate("NameOrderByAscOrderByDesc")
	if err == nil {
		t.Error("expected error for multiple OrderBy clauses")
	}
}
