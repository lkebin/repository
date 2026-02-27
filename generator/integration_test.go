package generator

import (
	"strings"
	"testing"

	"github.com/lkebin/repository/parser"
)

func loadTestModel(t *testing.T) *ModelSpecs {
	t.Helper()
	specs := ParseModel([]string{"User"}, []string{"../testdata"}, []string{})
	if len(specs) != 1 {
		t.Fatalf("expected 1 model spec, got %d", len(specs))
	}
	return &specs[0]
}

func loadTestRepository(t *testing.T) *RepositorySpecs {
	t.Helper()
	specs := ParseRepository([]string{"UserRepository"}, []string{"../testdata"}, []string{})
	if len(specs) != 1 {
		t.Fatalf("expected 1 repository spec, got %d", len(specs))
	}
	return &specs[0]
}

func TestLookupColumns(t *testing.T) {
	model := loadTestModel(t)
	columns := lookupColumns(model)

	expected := []struct {
		name     string
		property string
	}{
		{"id", "Id"},
		{"name", "Name"},
		{"birthday", "Birthday"},
		{"created_at", "CreatedAt"},
		{"updated_at", "UpdatedAt"},
	}

	if len(columns) != len(expected) {
		t.Fatalf("expected %d columns, got %d", len(expected), len(columns))
	}

	for i, want := range expected {
		if columns[i].Name != want.name {
			t.Errorf("column[%d].Name: expected %q, got %q", i, want.name, columns[i].Name)
		}
		if columns[i].Property != want.property {
			t.Errorf("column[%d].Property: expected %q, got %q", i, want.property, columns[i].Property)
		}
	}
}

func TestLookupPkColumn(t *testing.T) {
	model := loadTestModel(t)
	pk := lookupPkColumn(model)

	if pk == nil {
		t.Fatal("expected non-nil pk column")
	}
	if pk.Name != "id" {
		t.Errorf("expected pk column name %q, got %q", "id", pk.Name)
	}
	if pk.Property != "Id" {
		t.Errorf("expected pk property %q, got %q", "Id", pk.Property)
	}
}

func TestLookupColumnByProperty(t *testing.T) {
	model := loadTestModel(t)

	col := lookupColumnByProperty("Name", model)
	if col == nil {
		t.Fatal("expected non-nil column for property Name")
	}
	if col.Name != "name" {
		t.Errorf("expected column name %q, got %q", "name", col.Name)
	}

	missing := lookupColumnByProperty("NonExistent", model)
	if missing != nil {
		t.Error("expected nil for non-existent property")
	}
}

func TestPkFieldName(t *testing.T) {
	model := loadTestModel(t)
	name := PkFieldName(model)
	if name != "Id" {
		t.Errorf("expected %q, got %q", "Id", name)
	}
}

func TestIsPkAutoIncrement(t *testing.T) {
	model := loadTestModel(t)
	if !IsPkAutoIncrement(model) {
		t.Error("expected IsPkAutoIncrement to be true for testdata User")
	}
}

func TestGenSelectClause(t *testing.T) {
	model := loadTestModel(t)

	got := GenSelectClause(model, false)
	want := "SELECT `id`, `name`, `birthday`, `created_at`, `updated_at`"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}

	gotDistinct := GenSelectClause(model, true)
	wantDistinct := "SELECT DISTINCT `id`, `name`, `birthday`, `created_at`, `updated_at`"
	if gotDistinct != wantDistinct {
		t.Errorf("expected %q, got %q", wantDistinct, gotDistinct)
	}
}

func TestGenCountClause(t *testing.T) {
	model := loadTestModel(t)
	got := GenCountClause(model)
	want := "SELECT COUNT(`id`)"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestGenInsertClause(t *testing.T) {
	model := loadTestModel(t)
	got := GenInsertClause("user", model)
	want := "INSERT INTO `user` (`name`, `birthday`, `updated_at`) VALUES (?, ?, ?)"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestGenUpdateClause(t *testing.T) {
	model := loadTestModel(t)
	got := GenUpdateClause("user", model)
	want := "UPDATE `user` SET `name` = ?, `birthday` = ?, `updated_at` = ?"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestGenInsertColumns(t *testing.T) {
	model := loadTestModel(t)
	columns := GenInsertColumns(model)

	names := make([]string, len(columns))
	for i, c := range columns {
		names[i] = c.Name
	}

	if len(columns) != 3 {
		t.Fatalf("expected 3 insert columns (excluding pk autoincrement and unsafe), got %d: %v", len(columns), names)
	}

	expected := []string{"name", "birthday", "updated_at"}
	for i, want := range expected {
		if columns[i].Name != want {
			t.Errorf("column[%d]: expected %q, got %q", i, want, columns[i].Name)
		}
	}
}

func TestGenUpdateColumns(t *testing.T) {
	model := loadTestModel(t)
	columns := GenUpdateColumns(model)

	names := make([]string, len(columns))
	for i, c := range columns {
		names[i] = c.Name
	}

	if len(columns) != 3 {
		t.Fatalf("expected 3 update columns (excluding pk and autoincrement and unsafe), got %d: %v", len(columns), names)
	}

	expected := []string{"name", "birthday", "updated_at"}
	for i, want := range expected {
		if columns[i].Name != want {
			t.Errorf("column[%d]: expected %q, got %q", i, want, columns[i].Name)
		}
	}
}

func TestGenOrderByClause(t *testing.T) {
	model := loadTestModel(t)

	tests := []struct {
		name   string
		source string
		want   string
	}{
		{
			name:   "no order by",
			source: "FindByName",
			want:   "",
		},
		{
			name:   "single order by",
			source: "FindByNameOrderByNameAsc",
			want:   " ORDER BY name ASC",
		},
		{
			name:   "multiple order by",
			source: "FindByNameOrderByNameAscBirthdayDesc",
			want:   " ORDER BY name ASC, birthday DESC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt, err := parser.NewPartTree(tt.source)
			if err != nil {
				t.Fatalf("NewPartTree(%q) error: %v", tt.source, err)
			}
			got := GenOrderByClause(pt, model)
			if got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestGenWhereClausePredicate(t *testing.T) {
	model := loadTestModel(t)
	repo := loadTestRepository(t)

	methodByName := make(map[string]int)
	for i, m := range repo.Methods {
		methodByName[m.Name()] = i
	}

	tests := []struct {
		methodName string
		want       string
	}{
		{"FindByName", " WHERE `name` = ?"},
		{"FindByNameIsNull", " WHERE `name` IS NULL"},
		{"FindByNameAndBirthday", " WHERE (`name` = ?) AND (`birthday` = ?)"},
		{"FindByNameIsNullAndBirthday", " WHERE (`name` IS NULL) AND (`birthday` = ?)"},
		{"FindByNameIn", " WHERE `name` IN (?)"},
		{"FindByBirthdayBetween", " WHERE `birthday` BETWEEN ? AND ?"},
		{"FindByNameAndBirthdayIsNull", " WHERE (`name` = ?) AND (`birthday` IS NULL)"},
	}

	for _, tt := range tests {
		t.Run(tt.methodName, func(t *testing.T) {
			idx, ok := methodByName[tt.methodName]
			if !ok {
				t.Fatalf("method %q not found in repository", tt.methodName)
			}

			m := repo.Methods[idx]
			pt, err := parser.NewPartTree(tt.methodName)
			if err != nil {
				t.Fatalf("NewPartTree(%q) error: %v", tt.methodName, err)
			}

			got := GenWhereClausePredicate(pt, m.Signature().Params(), model)
			if got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestGenerateRepositoryImplements(t *testing.T) {
	repo := loadTestRepository(t)

	output, err := GenerateRepositoryImplements(repo)
	if err != nil {
		t.Fatalf("GenerateRepositoryImplements error: %v", err)
	}

	code := string(output)

	checks := []struct {
		desc    string
		content string
	}{
		{"package declaration", "package testdata"},
		{"struct definition", "type userRepositoryImpl struct"},
		{"constructor", "func NewUserRepository(db sqlx.ExtContext) UserRepository"},

		// CRUD methods
		{"Create method", "func (r *userRepositoryImpl) Create("},
		{"Update method", "func (r *userRepositoryImpl) Update("},
		{"FindById method", "func (r *userRepositoryImpl) FindById("},
		{"FindAll method", "func (r *userRepositoryImpl) FindAll("},
		{"DeleteById method", "func (r *userRepositoryImpl) DeleteById("},
		{"ExistsById method", "func (r *userRepositoryImpl) ExistsById("},
		{"Count method", "func (r *userRepositoryImpl) Count("},

		// Custom methods
		{"FindByName method", "func (r *userRepositoryImpl) FindByName("},
		{"FindByNameIsNull method", "func (r *userRepositoryImpl) FindByNameIsNull("},
		{"FindByNameAndBirthday method", "func (r *userRepositoryImpl) FindByNameAndBirthday("},
		{"FindByNameIsNullAndBirthday method", "func (r *userRepositoryImpl) FindByNameIsNullAndBirthday("},
		{"FindByNameIn method", "func (r *userRepositoryImpl) FindByNameIn("},
		{"FindByBirthdayBetween method", "func (r *userRepositoryImpl) FindByBirthdayBetween("},
		{"CountByName method", "func (r *userRepositoryImpl) CountByName("},
		{"ExistsByName method", "func (r *userRepositoryImpl) ExistsByName("},
		{"DeleteByName method", "func (r *userRepositoryImpl) DeleteByName("},
		{"FindByNameAndBirthdayIsNull method", "func (r *userRepositoryImpl) FindByNameAndBirthdayIsNull("},

		// SQL correctness for IsNull (the original bug)
		{"IsNull SQL has no trailing comma", "WHERE `name` IS NULL\")"}, // closes with ") not ", )"
		{"IsNull+And SQL", "WHERE (`name` IS NULL) AND (`birthday` = ?)"},
		{"And+IsNull SQL", "WHERE (`name` = ?) AND (`birthday` IS NULL)"},

		// SQL for other clauses
		{"BETWEEN SQL", "BETWEEN ? AND ?"},
		{"IN SQL", "IN (?)"},
		{"LIMIT clause", "LIMIT 10"},
		{"ORDER BY clause", "ORDER BY name ASC"},
	}

	for _, c := range checks {
		t.Run(c.desc, func(t *testing.T) {
			if !strings.Contains(code, c.content) {
				t.Errorf("generated code missing %q", c.content)
			}
		})
	}
}
