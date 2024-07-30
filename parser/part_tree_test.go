package parser

import "testing"

func TestNewPartTree(t *testing.T) {
	_, err := NewPartTree("FindByIdAndName")
	if err != nil {
		t.Error(err)
	}
}
