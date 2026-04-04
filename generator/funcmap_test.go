package generator

import (
	"go/types"
	"strings"
	"testing"
)

func TestParamName(t *testing.T) {
	repo := loadTestRepository(t)
	for _, m := range repo.Methods {
		if m.Name() == "FindByName" {
			got := ParamName(m.Signature().Params(), 0)
			if got != "ctx" {
				t.Errorf("expected first param name %q, got %q", "ctx", got)
			}
			got = ParamName(m.Signature().Params(), 1)
			if got != "name" {
				t.Errorf("expected second param name %q, got %q", "name", got)
			}
			return
		}
	}
	t.Fatal("FindByName method not found")
}

func TestGenCtxParam(t *testing.T) {
	repo := loadTestRepository(t)

	// Method with context.Context param
	for _, m := range repo.Methods {
		if m.Name() == "FindByName" {
			got := GenCtxParam(m.Signature().Params())
			if got != "ctx" {
				t.Errorf("expected %q, got %q", "ctx", got)
			}
			return
		}
	}
	t.Fatal("FindByName method not found")
}

func TestGenResultModelSliceSamePkg(t *testing.T) {
	repo := loadTestRepository(t)
	// FindByNameIsNull returns []*User -- slice, same package
	for _, m := range repo.Methods {
		if m.Name() == "FindByNameIsNull" {
			got := GenResultModel(m.Signature().Results(), repo)
			if !strings.Contains(got, "User") {
				t.Errorf("expected result model containing User, got %q", got)
			}
			if strings.Contains(got, "testdata.") {
				t.Errorf("same-pkg slice should not contain package prefix, got %q", got)
			}
			return
		}
	}
	t.Fatal("FindByNameIsNull method not found")
}

func TestGenResultModelPointerSamePkg(t *testing.T) {
	repo := loadTestRepository(t)
	// FindByName returns *User -- pointer, same package
	for _, m := range repo.Methods {
		if m.Name() == "FindByName" {
			got := GenResultModel(m.Signature().Results(), repo)
			if !strings.Contains(got, "User") {
				t.Errorf("expected result model containing User, got %q", got)
			}
			if strings.Contains(got, "testdata.") {
				t.Errorf("same-pkg pointer should not contain package prefix, got %q", got)
			}
			return
		}
	}
	t.Fatal("FindByName method not found")
}

func TestGenResultModelSingleResult(t *testing.T) {
	repo := loadTestRepository(t)
	// DeleteByName returns error -- single non-tuple result
	for _, m := range repo.Methods {
		if m.Name() == "DeleteByName" {
			got := GenResults(m.Signature().Results(), repo)
			if got != "error" {
				t.Errorf("expected %q, got %q", "error", got)
			}
			return
		}
	}
	t.Fatal("DeleteByName method not found")
}

func TestGenResultsMultiple(t *testing.T) {
	repo := loadTestRepository(t)
	// FindByName returns (*User, error)
	for _, m := range repo.Methods {
		if m.Name() == "FindByName" {
			got := GenResults(m.Signature().Results(), repo)
			if !strings.HasPrefix(got, "(") || !strings.HasSuffix(got, ")") {
				t.Errorf("multi-result should be parenthesized, got %q", got)
			}
			return
		}
	}
	t.Fatal("FindByName method not found")
}

func TestGenParamsOutput(t *testing.T) {
	repo := loadTestRepository(t)
	for _, m := range repo.Methods {
		if m.Name() == "FindByName" {
			got := GenParams(m.Signature().Params(), repo)
			if !strings.Contains(got, "ctx") || !strings.Contains(got, "name") {
				t.Errorf("expected params to contain ctx and name, got %q", got)
			}
			return
		}
	}
	t.Fatal("FindByName method not found")
}

func TestGenResultModelNonPointerNonSlice(t *testing.T) {
	// Return type is plain int64 (not pointer, not slice).
	param := types.NewVar(0, nil, "", types.Typ[types.Int64])
	tuple := types.NewTuple(param)
	repo := loadTestRepository(t)
	got := GenResultModel(tuple, repo)
	if got != "int64" {
		t.Errorf("expected %q, got %q", "int64", got)
	}
}

func TestGenResultModelSliceNoImportPath(t *testing.T) {
	// Return type is []string (slice, no import path).
	sliceType := types.NewSlice(types.Typ[types.String])
	param := types.NewVar(0, nil, "", sliceType)
	tuple := types.NewTuple(param)
	repo := loadTestRepository(t)
	got := GenResultModel(tuple, repo)
	if got != "[]string" {
		t.Errorf("expected %q, got %q", "[]string", got)
	}
}

func TestGenResultModelSliceDifferentPkg(t *testing.T) {
	// Return type is []*otherPkg.SomeType (slice, different pkg).
	otherPkg := types.NewPackage("github.com/other/pkg", "pkg")
	namedType := types.NewNamed(types.NewTypeName(0, otherPkg, "SomeType", nil), types.NewStruct(nil, nil), nil)
	sliceType := types.NewSlice(types.NewPointer(namedType))
	param := types.NewVar(0, nil, "", sliceType)
	tuple := types.NewTuple(param)
	repo := loadTestRepository(t)
	got := GenResultModel(tuple, repo)
	if !strings.Contains(got, "pkg.SomeType") {
		t.Errorf("expected result to contain %q, got %q", "pkg.SomeType", got)
	}
}

func TestGenResultModelPointerDifferentPkg(t *testing.T) {
	// Return type is *otherPkg.SomeType (pointer, different pkg).
	otherPkg := types.NewPackage("github.com/other/pkg", "pkg")
	namedType := types.NewNamed(types.NewTypeName(0, otherPkg, "SomeType", nil), types.NewStruct(nil, nil), nil)
	ptrType := types.NewPointer(namedType)
	param := types.NewVar(0, nil, "", ptrType)
	tuple := types.NewTuple(param)
	repo := loadTestRepository(t)
	got := GenResultModel(tuple, repo)
	if !strings.Contains(got, "pkg.SomeType") {
		t.Errorf("expected result to contain %q, got %q", "pkg.SomeType", got)
	}
}

func TestGenResultModelPointerNoImportPath(t *testing.T) {
	// Return type is *int64 (pointer, no import path).
	ptrType := types.NewPointer(types.Typ[types.Int64])
	param := types.NewVar(0, nil, "", ptrType)
	tuple := types.NewTuple(param)
	repo := loadTestRepository(t)
	got := GenResultModel(tuple, repo)
	if got != "int64" {
		t.Errorf("expected %q, got %q", "int64", got)
	}
}

func TestGenCtxParamFallback(t *testing.T) {
	// Create a Tuple with no context.Context param to trigger the "ctx" fallback.
	param := types.NewVar(0, nil, "name", types.Typ[types.String])
	tuple := types.NewTuple(param)
	got := GenCtxParam(tuple)
	if got != "ctx" {
		t.Errorf("expected fallback %q, got %q", "ctx", got)
	}
}

func TestGenVarBindingEmpty(t *testing.T) {
	// Single param (context only) should return empty string.
	param := types.NewVar(0, nil, "ctx", types.Typ[types.String])
	tuple := types.NewTuple(param)
	got := GenVarBinding(tuple)
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestLookupNamed(t *testing.T) {
	pkg := types.NewPackage("example.com/test", "test")
	named := types.NewNamed(types.NewTypeName(0, pkg, "Foo", nil), types.NewStruct(nil, nil), nil)

	// Direct named type
	if got := lookupNamed(named); got != named {
		t.Error("expected named type returned directly")
	}

	// Pointer wrapping named
	ptr := types.NewPointer(named)
	if got := lookupNamed(ptr); got != named {
		t.Error("expected named from pointer unwrap")
	}

	// Slice wrapping pointer wrapping named
	sl := types.NewSlice(types.NewPointer(named))
	if got := lookupNamed(sl); got != named {
		t.Error("expected named from slice->pointer unwrap")
	}

	// Basic type returns nil
	if got := lookupNamed(types.Typ[types.Int]); got != nil {
		t.Error("expected nil for basic type")
	}
}

func TestIsReturnSliceModel(t *testing.T) {
	repo := loadTestRepository(t)
	for _, m := range repo.Methods {
		switch m.Name() {
		case "FindByNameIsNull":
			if !IsReturnSliceModel(m.Signature().Results()) {
				t.Error("FindByNameIsNull should return slice model")
			}
		case "FindByName":
			if IsReturnSliceModel(m.Signature().Results()) {
				t.Error("FindByName should not return slice model")
			}
		}
	}
}
