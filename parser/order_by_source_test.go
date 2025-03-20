package parser

import "testing"

func TestOrderBySource(t *testing.T) {
	orderBySource, err := NewOrderBySource("PropertyAscPropertyaDescPropertybAsc")
	if err != nil {
		t.Fatal(err)
	}

	if len(orderBySource.Orders) != 3 {
		t.Errorf("expect 1, got %d", len(orderBySource.Orders))
	}

	if orderBySource.Orders[0].Property != "Property" {
		t.Errorf("expect Property, got %s", orderBySource.Orders[0].Property)
	}

	if orderBySource.Orders[0].Direction != "Asc" {
		t.Errorf("expect Asc, got %s", orderBySource.Orders[0].Direction)
	}

	if orderBySource.Orders[1].Property != "Propertya" {
		t.Errorf("expect Propertya, got %s", orderBySource.Orders[1].Property)
	}

	if orderBySource.Orders[1].Direction != "Desc" {
		t.Errorf("expect Desc, got %s", orderBySource.Orders[1].Direction)
	}

	if orderBySource.Orders[2].Property != "Propertyb" {
		t.Errorf("expect Propertyb, got %s", orderBySource.Orders[2].Property)
	}

	if orderBySource.Orders[2].Direction != "Asc" {
		t.Errorf("expect Asc, got %s", orderBySource.Orders[2].Direction)
	}
}
