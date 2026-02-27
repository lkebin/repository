package generator

import (
	"testing"
)

func TestParseRepository(t *testing.T) {
	specs := ParseRepository([]string{"UserRepository"}, []string{"../testdata"}, []string{})
	if len(specs) != 1 {
		t.Fatalf("expect 1 spec, got %d", len(specs))
	}

	spec := specs[0]
	if spec.Name != "UserRepository" {
		t.Errorf("expected name UserRepository, got %s", spec.Name)
	}

	if spec.Pkg == nil {
		t.Fatal("expected non-nil Pkg")
	}

	if len(spec.Methods) == 0 {
		t.Fatal("expected at least one method")
	}

	methodNames := make(map[string]bool)
	for _, m := range spec.Methods {
		methodNames[m.Name()] = true
	}

	for _, expected := range []string{"FindById", "FindAll", "Create", "Update", "DeleteById", "ExistsById", "Count", "FindByName"} {
		if !methodNames[expected] {
			t.Errorf("expected method %q not found", expected)
		}
	}
}

func TestParseModel(t *testing.T) {
	specs := ParseModel([]string{"User"}, []string{"../testdata"}, []string{})
	if len(specs) != 1 {
		t.Fatalf("expect 1 spec, got %d", len(specs))
	}

	spec := specs[0]
	if spec.Name != "User" {
		t.Errorf("expected name User, got %s", spec.Name)
	}

	if spec.Struct == nil {
		t.Fatal("expected non-nil Struct")
	}

	expectedFields := []string{"Id", "Name", "Birthday", "CreatedAt", "UpdatedAt"}
	if spec.Struct.NumFields() != len(expectedFields) {
		t.Fatalf("expected %d fields, got %d", len(expectedFields), spec.Struct.NumFields())
	}

	for i, name := range expectedFields {
		if spec.Struct.Field(i).Name() != name {
			t.Errorf("field[%d]: expected %q, got %q", i, name, spec.Struct.Field(i).Name())
		}
	}
}
