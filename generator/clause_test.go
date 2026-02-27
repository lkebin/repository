package generator

import (
	"testing"

	"github.com/lkebin/repository/parser"
)

func TestGenFromClause(t *testing.T) {
	tests := []struct {
		tableName string
		want      string
	}{
		{"user", " FROM `user`"},
		{"user_profile", " FROM `user_profile`"},
	}

	for _, tt := range tests {
		t.Run(tt.tableName, func(t *testing.T) {
			got := GenFromClause(tt.tableName)
			if got != tt.want {
				t.Errorf("GenFromClause(%q) = %q, want %q", tt.tableName, got, tt.want)
			}
		})
	}
}

func TestGenDeleteClause(t *testing.T) {
	tests := []struct {
		tableName string
		want      string
	}{
		{"user", "DELETE FROM `user`"},
		{"order_item", "DELETE FROM `order_item`"},
	}

	for _, tt := range tests {
		t.Run(tt.tableName, func(t *testing.T) {
			got := GenDeleteClause(tt.tableName)
			if got != tt.want {
				t.Errorf("GenDeleteClause(%q) = %q, want %q", tt.tableName, got, tt.want)
			}
		})
	}
}

func TestGenWhereClause(t *testing.T) {
	tests := []struct {
		name    string
		columns []*Column
		want    string
	}{
		{
			name:    "single column",
			columns: []*Column{{Name: "id"}},
			want:    " WHERE `id` = ?",
		},
		{
			name:    "two columns",
			columns: []*Column{{Name: "name"}, {Name: "birthday"}},
			want:    " WHERE `name` = ? AND `birthday` = ?",
		},
		{
			name:    "three columns",
			columns: []*Column{{Name: "a"}, {Name: "b"}, {Name: "c"}},
			want:    " WHERE `a` = ? AND `b` = ? AND `c` = ?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenWhereClause(tt.columns)
			if got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestGenLimitClause(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   string
	}{
		{"no limit", "FindByName", ""},
		{"First defaults to 1", "FindFirstByName", " LIMIT 1"},
		{"Top10", "FindTop10ByName", " LIMIT 10"},
		{"Top defaults to 1", "FindTopByName", " LIMIT 1"},
		{"First20", "FindFirst20ByName", " LIMIT 20"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt, err := parser.NewPartTree(tt.source)
			if err != nil {
				t.Fatalf("NewPartTree(%q) error: %v", tt.source, err)
			}
			got := GenLimitClause(pt)
			if got != tt.want {
				t.Errorf("GenLimitClause for %q = %q, want %q", tt.source, got, tt.want)
			}
		})
	}
}

func TestIsQueryIn(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   bool
	}{
		{"simple property - no In", "FindByName", false},
		{"In operator", "FindByNameIn", true},
		{"NotIn operator", "FindByNameNotIn", true},
		{"IsNull - no In", "FindByNameIsNull", false},
		{"mixed with In", "FindByNameAndBirthdayIn", true},
		{"Between - no In", "FindByAgeBetween", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt, err := parser.NewPartTree(tt.source)
			if err != nil {
				t.Fatalf("NewPartTree(%q) error: %v", tt.source, err)
			}
			got := IsQueryIn(pt)
			if got != tt.want {
				t.Errorf("IsQueryIn for %q = %v, want %v", tt.source, got, tt.want)
			}
		})
	}
}
