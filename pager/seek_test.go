package pager

import "testing"

func TestNewSeekPager(t *testing.T) {
	p := NewSeekPager()
	orderBy, seekWhere, seekValues, err := p.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if orderBy != "`id` ASC" {
		t.Errorf("default orderBy: expected `id` ASC, got %q", orderBy)
	}
	if seekWhere != "true" {
		t.Errorf("default seekWhere: expected true, got %q", seekWhere)
	}
	if len(seekValues) != 0 {
		t.Errorf("default seekValues: expected empty, got %v", seekValues)
	}
}

func TestSeekPagerSetOrder(t *testing.T) {
	tests := []struct {
		name      string
		pairs     []string
		wantOrder string
	}{
		{"single pair", []string{"name", OrderAsc}, "`name` ASC,`id` ASC"},
		{"single pair desc", []string{"name", OrderDesc}, "`name` DESC,`id` ASC"},
		{"multiple pairs", []string{"name", OrderAsc, "age", OrderDesc}, "`name` ASC,`age` DESC,`id` ASC"},
		{"odd count auto-appends asc", []string{"name"}, "`name` ASC,`id` ASC"},
		{"three items auto-appends asc", []string{"name", OrderDesc, "age"}, "`name` DESC,`age` ASC,`id` ASC"},
		{"four items", []string{"name", OrderAsc, "age", OrderDesc}, "`name` ASC,`age` DESC,`id` ASC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewSeekPager()
			p.SetOrder(tt.pairs...)
			orderBy, _, _, err := p.Build()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if orderBy != tt.wantOrder {
				t.Errorf("expected %q, got %q", tt.wantOrder, orderBy)
			}
		})
	}
}

func TestSeekPagerSetLastIDColumn(t *testing.T) {
	p := NewSeekPager()
	p.SetLastIDColumn("user_id")
	orderBy, _, _, err := p.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "`user_id` ASC"
	if orderBy != want {
		t.Errorf("expected %q, got %q", want, orderBy)
	}
}

func TestSeekPagerSetNext(t *testing.T) {
	tests := []struct {
		name          string
		orderPairs    []string
		lastIDColumn  string
		next          Next
		wantWhere     string
		wantValuesLen int
	}{
		{
			name:          "only id cursor - no order columns",
			orderPairs:    []string{},
			lastIDColumn:  "id",
			next:          NewNext(100),
			wantWhere:     "((`id` > ?))",
			wantValuesLen: 1,
		},
		{
			name:          "single order column asc",
			orderPairs:    []string{"created_at", OrderAsc},
			lastIDColumn:  "id",
			next:          NewNext("2024-01-01", 100),
			wantWhere:     "((`created_at` > ?) OR (`created_at` = ? AND `id` > ?))",
			wantValuesLen: 3,
		},
		{
			name:          "single order column desc",
			orderPairs:    []string{"created_at", OrderDesc},
			lastIDColumn:  "id",
			next:          NewNext("2024-01-01", 100),
			wantWhere:     "((`created_at` < ?) OR (`created_at` = ? AND `id` > ?))",
			wantValuesLen: 3,
		},
		{
			name:          "two order columns asc",
			orderPairs:    []string{"name", OrderAsc, "age", OrderAsc},
			lastIDColumn:  "id",
			next:          NewNext("Alice", 30, 100),
			wantWhere:     "((`name` > ?) OR (`name` = ? AND `age` > ?) OR (`name` = ? AND `age` = ? AND `id` > ?))",
			wantValuesLen: 6,
		},
		{
			name:          "two order columns mixed",
			orderPairs:    []string{"name", OrderAsc, "age", OrderDesc},
			lastIDColumn:  "id",
			next:          NewNext("Alice", 30, 100),
			wantWhere:     "((`name` > ?) OR (`name` = ? AND `age` < ?) OR (`name` = ? AND `age` = ? AND `id` > ?))",
			wantValuesLen: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewSeekPager()
			if len(tt.orderPairs) > 0 {
				p.SetOrder(tt.orderPairs...)
			}
			if tt.lastIDColumn != "" {
				p.SetLastIDColumn(tt.lastIDColumn)
			}
			p.SetNext(tt.next)

			_, seekWhere, seekValues, err := p.Build()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if seekWhere != tt.wantWhere {
				t.Errorf("expected where %q, got %q", tt.wantWhere, seekWhere)
			}
			if len(seekValues) != tt.wantValuesLen {
				t.Errorf("expected %d values, got %d: %v", tt.wantValuesLen, len(seekValues), seekValues)
			}
		})
	}
}

func TestSeekPagerBuild(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(SeekPager)
		wantOrderBy   string
		wantWhere     string
		wantValuesLen int
	}{
		{
			name:          "default",
			setup:         func(p SeekPager) {},
			wantOrderBy:   "`id` ASC",
			wantWhere:     "true",
			wantValuesLen: 0,
		},
		{
			name: "with single order",
			setup: func(p SeekPager) {
				p.SetOrder("name", OrderAsc)
			},
			wantOrderBy:   "`name` ASC,`id` ASC",
			wantWhere:     "true",
			wantValuesLen: 0,
		},
		{
			name: "with multiple orders",
			setup: func(p SeekPager) {
				p.SetOrder("name", OrderAsc, "age", OrderDesc)
			},
			wantOrderBy:   "`name` ASC,`age` DESC,`id` ASC",
			wantWhere:     "true",
			wantValuesLen: 0,
		},
		{
			name: "with order and next",
			setup: func(p SeekPager) {
				p.SetOrder("name", OrderAsc)
				p.SetNext(NewNext("test", 100))
			},
			wantOrderBy:   "`name` ASC,`id` ASC",
			wantWhere:     "((`name` > ?) OR (`name` = ? AND `id` > ?))",
			wantValuesLen: 3,
		},
		{
			name: "with multiple orders and next",
			setup: func(p SeekPager) {
				p.SetOrder("name", OrderAsc, "age", OrderDesc)
				p.SetNext(NewNext("test", 25, 100))
			},
			wantOrderBy:   "`name` ASC,`age` DESC,`id` ASC",
			wantWhere:     "((`name` > ?) OR (`name` = ? AND `age` < ?) OR (`name` = ? AND `age` = ? AND `id` > ?))",
			wantValuesLen: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewSeekPager()
			tt.setup(p)

			orderBy, seekWhere, seekValues, err := p.Build()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if orderBy != tt.wantOrderBy {
				t.Errorf("expected orderBy %q, got %q", tt.wantOrderBy, orderBy)
			}
			if seekWhere != tt.wantWhere {
				t.Errorf("expected where %q, got %q", tt.wantWhere, seekWhere)
			}
			if len(seekValues) != tt.wantValuesLen {
				t.Errorf("expected %d values, got %d", tt.wantValuesLen, len(seekValues))
			}
		})
	}
}

func TestNewNext(t *testing.T) {
	n := NewNext(1, "test", 3.14)
	if len(n) != 3 {
		t.Errorf("expected length 3, got %d", len(n))
	}
	if n[0] != 1 {
		t.Errorf("expected n[0]=1, got %v", n[0])
	}
	if n[1] != "test" {
		t.Errorf("expected n[1]=test, got %v", n[1])
	}
	if n[2] != 3.14 {
		t.Errorf("expected n[2]=3.14, got %v", n[2])
	}
}
