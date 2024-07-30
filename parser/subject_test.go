package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDistinct(t *testing.T) {
	s := NewSubject("FindDistinctByLastName")
	assert.True(t, s.isDistinct)
}

func TestIsCount(t *testing.T) {
	s := NewSubject("CountByLastName")
	assert.True(t, s.isCount)
}

func TestIsExists(t *testing.T) {
	s := NewSubject("ExistsByLastName")
	assert.True(t, s.isExists)
}

func TestIsDelete(t *testing.T) {
	s := NewSubject("DeleteByLastName")
	assert.True(t, s.isDelete)
}

func TestIsLimiting(t *testing.T) {
	s := NewSubject("FindFirstByLastName")
	assert.True(t, s.isLimiting)
	assert.Equal(t, 1, s.maxResults)

	s1 := NewSubject("FindFirst10ByLastName")
	assert.True(t, s1.isLimiting)
	assert.Equal(t, 10, s1.maxResults)

	s2 := NewSubject("FindTopByLastName")
	assert.True(t, s2.isLimiting)
	assert.Equal(t, 1, s2.maxResults)

	s3 := NewSubject("FindTop10ByLastName")
	assert.True(t, s3.isLimiting)
	assert.Equal(t, 10, s3.maxResults)
}
