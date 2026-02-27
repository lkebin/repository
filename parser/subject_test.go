package parser

import (
	"testing"
)

func TestNewSubject(t *testing.T) {
	tests := []struct {
		name       string
		subject    string
		isDistinct bool
		isCount    bool
		isExists   bool
		isDelete   bool
		isLimiting bool
		maxResults int
	}{
		{
			name:       "Find subject",
			subject:    "FindBy",
			isDistinct: false,
		},
		{
			name:       "Distinct",
			subject:    "FindDistinctBy",
			isDistinct: true,
		},
		{
			name:    "Count",
			subject: "CountBy",
			isCount: true,
		},
		{
			name:     "Exists",
			subject:  "ExistsBy",
			isExists: true,
		},
		{
			name:     "Delete",
			subject:  "DeleteBy",
			isDelete: true,
		},
		{
			name:     "Remove",
			subject:  "RemoveBy",
			isDelete: true,
		},
		{
			name:       "First (default 1)",
			subject:    "FindFirstBy",
			isLimiting: true,
			maxResults: 1,
		},
		{
			name:       "First10",
			subject:    "FindFirst10By",
			isLimiting: true,
			maxResults: 10,
		},
		{
			name:       "Top (default 1)",
			subject:    "FindTopBy",
			isLimiting: true,
			maxResults: 1,
		},
		{
			name:       "Top10",
			subject:    "FindTop10By",
			isLimiting: true,
			maxResults: 10,
		},
		{
			name:       "empty subject",
			subject:    "",
			isDistinct: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSubject(tt.subject)

			if s.IsDistinct != tt.isDistinct {
				t.Errorf("IsDistinct: expected %v, got %v", tt.isDistinct, s.IsDistinct)
			}
			if s.IsCount != tt.isCount {
				t.Errorf("IsCount: expected %v, got %v", tt.isCount, s.IsCount)
			}
			if s.IsExists != tt.isExists {
				t.Errorf("IsExists: expected %v, got %v", tt.isExists, s.IsExists)
			}
			if s.IsDelete != tt.isDelete {
				t.Errorf("IsDelete: expected %v, got %v", tt.isDelete, s.IsDelete)
			}
			if s.IsLimiting != tt.isLimiting {
				t.Errorf("IsLimiting: expected %v, got %v", tt.isLimiting, s.IsLimiting)
			}
			if tt.isLimiting && s.MaxResults != tt.maxResults {
				t.Errorf("MaxResults: expected %d, got %d", tt.maxResults, s.MaxResults)
			}
		})
	}
}
