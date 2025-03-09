package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDistinct(t *testing.T) {
	s := NewSubject("FindDistinctByLastName")
	assert.True(t, s.IsDistinct)
}

func TestIsCount(t *testing.T) {
	s := NewSubject("CountByLastName")
	assert.True(t, s.IsCount)
}

func TestIsExists(t *testing.T) {
	s := NewSubject("ExistsByLastName")
	assert.True(t, s.IsExists)
}

func TestIsDelete(t *testing.T) {
	s := NewSubject("DeleteByLastName")
	assert.True(t, s.IsDelete)
}

func TestIsLimiting(t *testing.T) {
	s := NewSubject("FindFirstByLastName")
	assert.True(t, s.IsLimiting)
	assert.Equal(t, 1, s.MaxResults)

	s1 := NewSubject("FindFirst10ByLastName")
	assert.True(t, s1.IsLimiting)
	assert.Equal(t, 10, s1.MaxResults)

	s2 := NewSubject("FindTopByLastName")
	assert.True(t, s2.IsLimiting)
	assert.Equal(t, 1, s2.MaxResults)

	s3 := NewSubject("FindTop10ByLastName")
	assert.True(t, s3.IsLimiting)
	assert.Equal(t, 10, s3.MaxResults)
}
