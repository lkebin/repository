package parser

import (
	"testing"
)

func TestNewPartTree(t *testing.T) {
	pt, err := NewPartTree("FindByIdAndName")
	if err != nil {
		t.Error(err)
	}

	if pt.Subject == nil {
		t.Error("Subject should not be nil")
	}

	if len(pt.Predicate.Nodes) != 1 {
		t.Errorf("Expected 2 nodes in predicate, got %d", len(pt.Predicate.Nodes))
	}

	if len(pt.Predicate.Nodes[0].Children) != 2 {
		t.Errorf("Expected 2 children in first node, got %d", len(pt.Predicate.Nodes[0].Children))
	}
}
