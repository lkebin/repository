package parser

import "testing"

func TestNewOrderBySource(t *testing.T) {
	tests := []struct {
		name       string
		clause     string
		wantOrders []Order
	}{
		{
			name:   "single ascending",
			clause: "NameAsc",
			wantOrders: []Order{
				{Property: "Name", Direction: "Asc"},
			},
		},
		{
			name:   "single descending",
			clause: "NameDesc",
			wantOrders: []Order{
				{Property: "Name", Direction: "Desc"},
			},
		},
		{
			name:   "multiple properties",
			clause: "PropertyAscPropertyaDescPropertybAsc",
			wantOrders: []Order{
				{Property: "Property", Direction: "Asc"},
				{Property: "Propertya", Direction: "Desc"},
				{Property: "Propertyb", Direction: "Asc"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obs, err := NewOrderBySource(tt.clause)
			if err != nil {
				t.Fatalf("NewOrderBySource(%q) error: %v", tt.clause, err)
			}

			if len(obs.Orders) != len(tt.wantOrders) {
				t.Fatalf("expected %d orders, got %d", len(tt.wantOrders), len(obs.Orders))
			}

			for i, want := range tt.wantOrders {
				got := obs.Orders[i]
				if got.Property != want.Property {
					t.Errorf("order[%d].Property: expected %q, got %q", i, want.Property, got.Property)
				}
				if got.Direction != want.Direction {
					t.Errorf("order[%d].Direction: expected %q, got %q", i, want.Direction, got.Direction)
				}
			}
		})
	}
}

func TestNewOrderBySourceEmpty(t *testing.T) {
	obs, err := NewOrderBySource("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obs != nil {
		t.Error("expected nil for empty clause")
	}
}

func TestNewOrderBySourceErrors(t *testing.T) {
	tests := []struct {
		name   string
		clause string
	}{
		{"missing direction", "Name"},
		{"direction only", "Asc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewOrderBySource(tt.clause)
			if err == nil {
				t.Errorf("NewOrderBySource(%q): expected error, got nil", tt.clause)
			}
			if err != ErrInvalidOrderSyntax {
				t.Errorf("expected ErrInvalidOrderSyntax, got %v", err)
			}
		})
	}
}
