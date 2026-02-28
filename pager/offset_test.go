package pager

import "testing"

func TestNewOffsetPager(t *testing.T) {
	p := NewOffsetPager()
	orderBy, page, err := p.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page != 1 {
		t.Errorf("default page: expected 1, got %d", page)
	}
	if orderBy != "" {
		t.Errorf("default orderBy: expected empty, got %q", orderBy)
	}
}

func TestSetPage(t *testing.T) {
	tests := []struct {
		name     string
		page     int64
		wantPage int64
	}{
		{"positive page", 5, 5},
		{"page 1", 1, 1},
		{"page 0", 0, 0},
		{"negative page resets to 1", -1, 1},
		{"large negative resets to 1", -100, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewOffsetPager()
			p.SetPage(tt.page)
			_, page, _ := p.Build()
			if page != tt.wantPage {
				t.Errorf("expected page %d, got %d", tt.wantPage, page)
			}
		})
	}
}

func TestSetOrder(t *testing.T) {
	tests := []struct {
		name      string
		pairs     []string
		wantOrder string
	}{
		{"single pair", []string{"name", Asc}, "`name` ASC"},
		{"single pair desc", []string{"age", Desc}, "`age` DESC"},
		{"multiple pairs", []string{"name", Asc, "age", Desc}, "`name` ASC,`age` DESC"},
		{"odd count auto-appends Asc", []string{"name"}, "`name` ASC"},
		{"three items auto-appends Asc", []string{"name", Desc, "age"}, "`name` DESC,`age` ASC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewOffsetPager()
			p.SetOrder(tt.pairs...)
			orderBy, _, err := p.Build()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if orderBy != tt.wantOrder {
				t.Errorf("expected %q, got %q", tt.wantOrder, orderBy)
			}
		})
	}
}

func TestBuild(t *testing.T) {
	p := NewOffsetPager()
	p.SetPage(3)
	p.SetOrder("name", Asc, "birthday", Desc)

	orderBy, page, err := p.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page != 3 {
		t.Errorf("expected page 3, got %d", page)
	}
	want := "`name` ASC,`birthday` DESC"
	if orderBy != want {
		t.Errorf("expected %q, got %q", want, orderBy)
	}
}
