package generator

import (
	"go/types"
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

func loadNoPkModel(t *testing.T) *ModelSpecs {
	t.Helper()
	specs := ParseModel([]string{"NoPkModel"}, []string{"../testdata"}, []string{})
	if len(specs) != 1 {
		t.Fatalf("expected 1 model spec, got %d", len(specs))
	}
	return &specs[0]
}

func loadNoPkRepository(t *testing.T) *RepositorySpecs {
	t.Helper()
	specs := ParseRepository([]string{"NoPkRepository"}, []string{"../testdata"}, []string{})
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

func TestGenUpdateColumnsAutoIncrement(t *testing.T) {
	specs := ParseModel([]string{"AutoIncModel"}, []string{"../testdata"}, []string{})
	if len(specs) != 1 {
		t.Fatalf("expected 1 model spec, got %d", len(specs))
	}
	model := &specs[0]
	columns := GenUpdateColumns(model)

	names := make([]string, len(columns))
	for i, c := range columns {
		names[i] = c.Name
	}

	expected := []string{"name", "birthday"}
	if len(columns) != len(expected) {
		t.Fatalf("expected %d update columns (excluding pk and standalone autoincrement), got %d: %v", len(expected), len(columns), names)
	}
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
			got, err := GenOrderByClause(pt, model)
			if err != nil {
				t.Fatalf("GenOrderByClause error: %v", err)
			}
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
		{"FindByNameOrBirthday", " WHERE (`name` = ?) OR (`birthday` = ?)"},
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

			got, err := GenWhereClausePredicate(pt, m.Signature().Params(), model)
			if err != nil {
				t.Fatalf("GenWhereClausePredicate error: %v", err)
			}
			if got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestGenWhereClausePredicateParamMismatch(t *testing.T) {
	model := loadTestModel(t)
	repo := loadTestRepository(t)

	// Use FindByName's params (has 1 binding param) but parse
	// FindByNameAndBirthday (needs 2 binding params) to trigger mismatch.
	var findByNameParams *types.Tuple
	for _, m := range repo.Methods {
		if m.Name() == "FindByName" {
			findByNameParams = m.Signature().Params()
			break
		}
	}
	if findByNameParams == nil {
		t.Fatal("FindByName method not found")
	}

	pt, err := parser.NewPartTree("FindByNameAndBirthday")
	if err != nil {
		t.Fatalf("NewPartTree error: %v", err)
	}

	_, err = GenWhereClausePredicate(pt, findByNameParams, model)
	if err == nil {
		t.Error("expected error for param count mismatch")
	}
}

func TestGenWhereClausePredicateColumnNotFound(t *testing.T) {
	model := loadTestModel(t)
	repo := loadTestRepository(t)

	// FindByName has 1 param. We'll parse "FindByNonExistent" which
	// references a column that doesn't exist in the model.
	var findByNameParams *types.Tuple
	for _, m := range repo.Methods {
		if m.Name() == "FindByName" {
			findByNameParams = m.Signature().Params()
			break
		}
	}

	pt, err := parser.NewPartTree("FindByNonExistent")
	if err != nil {
		t.Fatalf("NewPartTree error: %v", err)
	}

	_, err = GenWhereClausePredicate(pt, findByNameParams, model)
	if err == nil {
		t.Error("expected error for column not found")
	}
}

func TestGenOrderByClauseColumnNotFound(t *testing.T) {
	model := loadTestModel(t)

	pt, err := parser.NewPartTree("FindByNameOrderByNonExistentAsc")
	if err != nil {
		t.Fatalf("NewPartTree error: %v", err)
	}

	_, err = GenOrderByClause(pt, model)
	if err == nil {
		t.Error("expected error for column not found in order by")
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
		{"FindByNameOrBirthday method", "func (r *userRepositoryImpl) FindByNameOrBirthday("},
		{"DeleteByNameIsNull method", "func (r *userRepositoryImpl) DeleteByNameIsNull("},

		// SQL correctness for IsNull (the original bug)
		{"IsNull SQL has no trailing comma", "WHERE `name` IS NULL\")"}, // closes with ") not ", )"
		{"IsNull+And SQL", "WHERE (`name` IS NULL) AND (`birthday` = ?)"},
		{"And+IsNull SQL", "WHERE (`name` = ?) AND (`birthday` IS NULL)"},

		// OR clause
		{"OR SQL", "WHERE (`name` = ?) OR (`birthday` = ?)"},
		// Delete with IsNull (zero-arg)
		{"DeleteByNameIsNull SQL", "DELETE FROM `user` WHERE `name` IS NULL"},

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

func TestGenerateNoPkRepositoryError(t *testing.T) {
	repo := loadNoPkRepository(t)
	_, err := GenerateRepositoryImplements(repo)
	if err == nil {
		t.Error("expected error for repository with no PK model")
	}
}

func TestGenFuncImplNoPkErrors(t *testing.T) {
	noPkModel := loadNoPkModel(t)
	repo := loadNoPkRepository(t)

	methodsByName := make(map[string]*types.Func)
	for _, m := range repo.Methods {
		methodsByName[m.Name()] = m
	}

	needPk := []string{"Update", "FindById", "ExistsById", "DeleteById"}
	for _, name := range needPk {
		t.Run(name, func(t *testing.T) {
			m, ok := methodsByName[name]
			if !ok {
				t.Skipf("method %s not found", name)
				return
			}
			_, err := genFuncImpl("testImpl", noPkModel, m, repo)
			if err == nil {
				t.Errorf("expected error for %s with no PK", name)
			}
			if !strings.Contains(err.Error(), "pk column not found") {
				t.Errorf("expected pk column error, got: %v", err)
			}
		})
	}
}

func TestLookupPkColumnNil(t *testing.T) {
	m := loadNoPkModel(t)
	col := lookupPkColumn(m)
	if col != nil {
		t.Errorf("expected nil pk column for NoPkModel, got %+v", col)
	}
}

func TestIsPkAutoIncrementNoPk(t *testing.T) {
	m := loadNoPkModel(t)
	if IsPkAutoIncrement(m) {
		t.Error("expected false for model with no PK")
	}
}

func TestPkFieldNameNoPk(t *testing.T) {
	m := loadNoPkModel(t)
	if PkFieldName(m) != "" {
		t.Errorf("expected empty string for model with no PK, got %q", PkFieldName(m))
	}
}
