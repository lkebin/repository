package filter

import (
	"strings"
	"testing"
)

func TestNewFilter(t *testing.T) {
	f := NewFilter()
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
	w, v, err := f.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w != "true" {
		t.Errorf("empty filter should return %q, got %q", "true", w)
	}
	if len(v) != 0 {
		t.Errorf("empty filter should have no values, got %d", len(v))
	}
}

func TestRulesGetSet(t *testing.T) {
	f := NewFilter()
	rules := []*Rule{{Column: "a", Type: "term", Operator: "eq", Query: 1}}
	f.SetRules(rules)
	if len(f.Rules()) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(f.Rules()))
	}
	if f.Rules()[0].Column != "a" {
		t.Errorf("expected column %q, got %q", "a", f.Rules()[0].Column)
	}
}

func TestRuleConstructors(t *testing.T) {
	tests := []struct {
		name     string
		rule     *Rule
		wantType string
		wantOp   string
		wantCol  string
	}{
		{"In", In("col", 1, 2), "terms", "in", "col"},
		{"NotIn", NotIn("col", 3), "terms", "notin", "col"},
		{"InJSON", InJSON("col", "a"), "terms", "injson", "col"},
		{"NotInJSON", NotInJSON("col", "b"), "terms", "notinjson", "col"},
		{"Is", Is("col", "x"), "term", "is", "col"},
		{"Not", Not("col", "y"), "term", "not", "col"},
		{"Eq", Eq("col", 1), "term", "eq", "col"},
		{"Ne", Ne("col", 2), "term", "ne", "col"},
		{"Gt", Gt("col", 3), "range", "gt", "col"},
		{"Gte", Gte("col", 4), "range", "gte", "col"},
		{"Lt", Lt("col", 5), "range", "lt", "col"},
		{"Lte", Lte("col", 6), "range", "lte", "col"},
		{"Between", Between("col", 1, 10), "range", "between", "col"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.rule.Type != tt.wantType {
				t.Errorf("Type: expected %q, got %q", tt.wantType, tt.rule.Type)
			}
			if tt.rule.Operator != tt.wantOp {
				t.Errorf("Operator: expected %q, got %q", tt.wantOp, tt.rule.Operator)
			}
			if tt.rule.Column != tt.wantCol {
				t.Errorf("Column: expected %q, got %q", tt.wantCol, tt.rule.Column)
			}
		})
	}
}

func TestExistNotExistConstructors(t *testing.T) {
	e := Exist("col")
	if e.Type != "exists" || e.Inverse {
		t.Errorf("Exist: expected type=exists inverse=false, got type=%s inverse=%v", e.Type, e.Inverse)
	}
	ne := NotExist("col")
	if ne.Type != "exists" || !ne.Inverse {
		t.Errorf("NotExist: expected type=exists inverse=true, got type=%s inverse=%v", ne.Type, ne.Inverse)
	}
}

func TestAndOrConstructors(t *testing.T) {
	a := And(Eq("a", 1), Eq("b", 2))
	if a.Type != "and" || len(a.Rules) != 2 {
		t.Errorf("And: expected type=and rules=2, got type=%s rules=%d", a.Type, len(a.Rules))
	}
	o := Or(Eq("a", 1), Eq("b", 2))
	if o.Type != "or" || len(o.Rules) != 2 {
		t.Errorf("Or: expected type=or rules=2, got type=%s rules=%d", o.Type, len(o.Rules))
	}
}

func TestInterfaceMethodsAppendRules(t *testing.T) {
	f := NewFilter()
	f.In("a", 1)
	f.NotIn("b", 2)
	f.InJSON("c", "x")
	f.NotInJSON("d", "y")
	f.Is("e", 1)
	f.Not("f", 2)
	f.Exist("g")
	f.NotExist("h")
	f.Eq("i", 1)
	f.Ne("j", 2)
	f.Gt("k", 3)
	f.Gte("l", 4)
	f.Lt("m", 5)
	f.Lte("n", 6)
	f.Between("o", 1, 10)
	f.And(Eq("p", 1))
	f.Or(Eq("q", 2))

	if len(f.Rules()) != 17 {
		t.Fatalf("expected 17 rules, got %d", len(f.Rules()))
	}
}

func TestBuildSingleOperators(t *testing.T) {
	tests := []struct {
		name      string
		rule      *Rule
		wantWhere string
		wantVals  int
	}{
		{"eq", Eq("name", "alice"), "(name = ?)", 1},
		{"ne", Ne("name", "bob"), "(name != ?)", 1},
		{"is", Is("status", "active"), "(status = ?)", 1},
		{"not", Not("status", "inactive"), "(status != ?)", 1},
		{"gt", Gt("age", 18), "(age > ?)", 1},
		{"gte", Gte("age", 18), "(age >= ?)", 1},
		{"lt", Lt("age", 65), "(age < ?)", 1},
		{"lte", Lte("age", 65), "(age <= ?)", 1},
		{"between", Between("age", 18, 65), "(age BETWEEN ? AND ?)", 2},
		{"exist", Exist("email"), "(email IS NOT NULL)", 0},
		{"not exist", NotExist("email"), "(email IS NULL)", 0},
		{"in", In("id", 1, 2, 3), "", 3},       // sqlx.In rewrites the query
		{"notin", NotIn("id", 1, 2, 3), "", 3}, // sqlx.In rewrites the query
		{"injson", InJSON("tags", "a", "b"), "", 1},
		{"notinjson", NotInJSON("tags", "a", "b"), "", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewFilter()
			f.SetRules([]*Rule{tt.rule})
			w, v, err := f.Build()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantWhere != "" && w != tt.wantWhere {
				t.Errorf("where: expected %q, got %q", tt.wantWhere, w)
			}
			if len(v) != tt.wantVals {
				t.Errorf("values: expected %d, got %d", tt.wantVals, len(v))
			}
		})
	}
}

func TestBuildInQueries(t *testing.T) {
	f := NewFilter()
	f.In("id", 1, 2, 3)
	w, v, err := f.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(w, "IN") {
		t.Errorf("expected IN clause, got %q", w)
	}
	if len(v) != 3 {
		t.Errorf("expected 3 values, got %d", len(v))
	}
}

func TestBuildNotInQuery(t *testing.T) {
	f := NewFilter()
	f.NotIn("id", 4, 5)
	w, v, err := f.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(w, "NOT IN") {
		t.Errorf("expected NOT IN clause, got %q", w)
	}
	if len(v) != 2 {
		t.Errorf("expected 2 values, got %d", len(v))
	}
}

func TestBuildInJSONQuery(t *testing.T) {
	f := NewFilter()
	f.InJSON("tags", "a", "b")
	w, v, err := f.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(w, "JSON_CONTAINS") {
		t.Errorf("expected JSON_CONTAINS clause, got %q", w)
	}
	if len(v) != 1 {
		t.Errorf("expected 1 value (JSON string), got %d", len(v))
	}
}

func TestBuildNotInJSONQuery(t *testing.T) {
	f := NewFilter()
	f.NotInJSON("tags", "x")
	w, v, err := f.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(w, "!JSON_CONTAINS") {
		t.Errorf("expected !JSON_CONTAINS clause, got %q", w)
	}
	if len(v) != 1 {
		t.Errorf("expected 1 value, got %d", len(v))
	}
}

func TestBuildMultipleRules(t *testing.T) {
	f := NewFilter()
	f.Eq("name", "alice")
	f.Gt("age", 18)
	w, v, err := f.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(w, " AND ") {
		t.Errorf("expected AND between rules, got %q", w)
	}
	if len(v) != 2 {
		t.Errorf("expected 2 values, got %d", len(v))
	}
}

func TestBuildAndComposite(t *testing.T) {
	f := NewFilter()
	f.And(Eq("a", 1), Gt("b", 2))
	w, v, err := f.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(w, "(") {
		t.Errorf("expected parenthesized AND group, got %q", w)
	}
	if len(v) != 2 {
		t.Errorf("expected 2 values, got %d", len(v))
	}
}

func TestBuildOrComposite(t *testing.T) {
	f := NewFilter()
	f.Or(Eq("a", 1), Eq("b", 2))
	w, v, err := f.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(w, " OR ") {
		t.Errorf("expected OR in result, got %q", w)
	}
	if len(v) != 2 {
		t.Errorf("expected 2 values, got %d", len(v))
	}
}

func TestBuildErrorMissingOperator(t *testing.T) {
	f := NewFilter()
	f.SetRules([]*Rule{{Column: "a", Type: "term", Operator: "", Query: 1}})
	_, _, err := f.Build()
	if err == nil {
		t.Error("expected error for missing operator")
	}
}

func TestBuildErrorMissingQuery(t *testing.T) {
	f := NewFilter()
	f.SetRules([]*Rule{{Column: "a", Type: "term", Operator: "eq", Query: nil}})
	_, _, err := f.Build()
	if err == nil {
		t.Error("expected error for missing query")
	}
}

func TestBuildErrorInvalidRuleType(t *testing.T) {
	f := NewFilter()
	f.SetRules([]*Rule{{Column: "a", Type: "unknown"}})
	_, _, err := f.Build()
	if err == nil {
		t.Error("expected error for invalid rule type")
	}
}

func TestBuildBetweenMissingGte(t *testing.T) {
	f := NewFilter()
	f.SetRules([]*Rule{{
		Column:   "age",
		Type:     "range",
		Operator: "between",
		Query:    map[string]any{"lte": 65},
	}})
	_, _, err := f.Build()
	if err == nil {
		t.Error("expected error for between missing gte key")
	}
}

func TestBuildBetweenMissingLte(t *testing.T) {
	f := NewFilter()
	f.SetRules([]*Rule{{
		Column:   "age",
		Type:     "range",
		Operator: "between",
		Query:    map[string]any{"gte": 18},
	}})
	_, _, err := f.Build()
	if err == nil {
		t.Error("expected error for between missing lte key")
	}
}

func TestBuildAndCompositeErrorPropagation(t *testing.T) {
	f := NewFilter()
	f.And(Eq("a", 1), &Rule{Column: "b", Type: "term", Operator: "", Query: 2})
	_, _, err := f.Build()
	if err == nil {
		t.Error("expected error to propagate from And child")
	}
}

func TestBuildOrCompositeErrorPropagation(t *testing.T) {
	f := NewFilter()
	f.Or(Eq("a", 1), &Rule{Column: "b", Type: "term", Operator: "", Query: 2})
	_, _, err := f.Build()
	if err == nil {
		t.Error("expected error to propagate from Or child")
	}
}

func TestBuildInEmptySliceError(t *testing.T) {
	f := NewFilter()
	f.SetRules([]*Rule{{Column: "id", Type: "terms", Operator: "in", Query: []any{}}})
	_, _, err := f.Build()
	if err == nil {
		t.Error("expected error for IN with empty slice")
	}
}

func TestBuildNotInEmptySliceError(t *testing.T) {
	f := NewFilter()
	f.SetRules([]*Rule{{Column: "id", Type: "terms", Operator: "notin", Query: []any{}}})
	_, _, err := f.Build()
	if err == nil {
		t.Error("expected error for NOT IN with empty slice")
	}
}
