package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPredicate(t *testing.T) {
	p, err := NewPredicate("AndroidEqualsAllIgnoreCase")
	assert.Nil(t, err)
	assert.Len(t, p.Nodes, 1)
}
