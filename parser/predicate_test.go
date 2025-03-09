package parser

import (
	"testing"
)

func TestNewPredicate(t *testing.T) {
	p, err := NewPredicate("AndroidEqualsAllIgnoreCase")
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if len(p.Nodes) != 1 {
		t.Errorf("expect 1, got %d", len(p.Nodes))
	}
}
