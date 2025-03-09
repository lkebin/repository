package parser

import (
	"testing"
)

func TestIsDistinct(t *testing.T) {
	s := NewSubject("FindDistinctByLastName")
	if s.IsDistinct == false {
		t.Errorf("expect true, got false")
	}
}

func TestIsCount(t *testing.T) {
	s := NewSubject("CountByLastName")
	if s.IsCount == false {
		t.Errorf("expect true, got false")
	}
}

func TestIsExists(t *testing.T) {
	s := NewSubject("ExistsByLastName")
	if s.IsExists == false {
		t.Errorf("expect true, got false")
	}
}

func TestIsDelete(t *testing.T) {
	s := NewSubject("DeleteByLastName")
	if s.IsDelete == false {
		t.Errorf("expect true, got false")
	}
}

func TestIsLimiting(t *testing.T) {
	s := NewSubject("FindFirstByLastName")
	if s.IsLimiting == false {
		t.Errorf("expect true, got false")
	}

	if s.MaxResults != 1 {
		t.Errorf("expect 1, got %d", s.MaxResults)
	}

	s1 := NewSubject("FindFirst10ByLastName")
	if s1.IsLimiting == false {
		t.Errorf("expect true, got false")
	}

	if s1.MaxResults != 10 {
		t.Errorf("expect 10, got %d", s.MaxResults)
	}

	s2 := NewSubject("FindTopByLastName")
	if s2.IsLimiting == false {
		t.Errorf("expect true, got false")
	}

	if s2.MaxResults != 1 {
		t.Errorf("expect 1, got %d", s.MaxResults)
	}

	s3 := NewSubject("FindTop10ByLastName")
	if s3.IsLimiting == false {
		t.Errorf("expect true, got false")
	}

	if s3.MaxResults != 10 {
		t.Errorf("expect 10, got %d", s.MaxResults)
	}
}
