package generator

import (
	"fmt"
	"go/types"
	"reflect"
	"strings"
	"text/template"

	"github.com/lkebin/repository/parser"
)

func funcMap() template.FuncMap {
	fm := make(template.FuncMap)
	fm["Add"] = func(a, b int) int {
		return a + b
	}

	fm["Params"] = GenParams
	fm["CtxParam"] = GenCtxParam
	fm["Results"] = GenResults
	fm["SelectClause"] = GenSelectClause
	fm["ExistsClause"] = GenCountClause
	fm["CountClause"] = GenCountClause
	fm["InsertClause"] = GenInsertClause
	fm["UpdateClause"] = GenUpdateClause
	fm["DeleteClause"] = GenDeleteClause
	fm["FromClause"] = GenFromClause
	fm["WhereClausePredicate"] = GenWhereClausePredicate
	fm["WhereClause"] = GenWhereClause
	fm["OrderByClause"] = GenOrderByClause
	fm["LimitClause"] = GenLimitClause
	fm["VarBinding"] = GenVarBinding
	fm["ResultModel"] = GenResultModel
	fm["IsReturnSliceModel"] = IsReturnSliceModel
	fm["IsPkAutoIncrement"] = IsPkAutoIncrement
	fm["IsQueryIn"] = IsQueryIn
	fm["InsertFieldBinding"] = GenInsertFieldBinding
	fm["UpdateFieldBinding"] = GenUpdateFieldBinding
	fm["ParamName"] = ParamName
	fm["PkFieldName"] = PkFieldName
	return fm
}

func parseTuple(params *types.Tuple, spec *RepositorySpecs) string {
	var p = &strings.Builder{}
	for i := range params.Len() {
		if params.At(i).Name() != "" {
			p.WriteString(params.At(i).Name())
			p.WriteString(" ")
		}

		importPath := lookupPkgPath(params.At(i).Type())
		if importPath != "" {
			if importPath == spec.Pkg.Path() {
				p.WriteString(strings.ReplaceAll(params.At(i).Type().String(), fmt.Sprintf("%s.", importPath), ""))
			} else {
				p.WriteString(strings.ReplaceAll(params.At(i).Type().String(), importPath, strings.Split(importPath, "/")[len(strings.Split(importPath, "/"))-1]))
			}
		} else {
			p.WriteString(params.At(i).Type().String())
		}

		if i < params.Len()-1 {
			p.WriteString(", ")
		}
	}
	return p.String()
}

func GenParams(params *types.Tuple, spec *RepositorySpecs) string {
	return parseTuple(params, spec)
}

func GenCtxParam(params *types.Tuple) string {
	for i := range params.Len() {
		if (params.At(i).Type().String()) == "context.Context" {
			return params.At(i).Name()
		}
	}
	return "ctx"
}

func GenResults(results *types.Tuple, spec *RepositorySpecs) string {
	if results.Len() > 1 {
		return fmt.Sprintf("(%s)", parseTuple(results, spec))
	}
	return parseTuple(results, spec)
}

func GenVarBinding(params *types.Tuple) string {
	if params.Len() <= 1 {
		return ""
	}
	var s = &strings.Builder{}
	s.WriteString(", ")
	for i := 1; i < params.Len(); i++ {
		s.WriteString(params.At(i).Name())
		if i < params.Len()-1 {
			s.WriteString(", ")
		}
	}
	return s.String()
}

func GenResultModel(results *types.Tuple, spec *RepositorySpecs) string {
	if IsReturnSliceModel(results) {
		importPath := lookupPkgPath(results.At(0).Type())
		if importPath != "" {
			if importPath == spec.Pkg.Path() {
				return strings.ReplaceAll(results.At(0).Type().String(), fmt.Sprintf("%s.", importPath), "")
			} else {
				return strings.ReplaceAll(results.At(0).Type().String(), importPath, strings.Split(importPath, "/")[len(strings.Split(importPath, "/"))-1])
			}
		}
		return results.At(0).Type().String()
	}

	typ, ok := results.At(0).Type().(*types.Pointer)
	if !ok {
		return results.At(0).Type().String()
	}

	importPath := lookupPkgPath(typ.Elem())
	if importPath != "" {
		if importPath == spec.Pkg.Path() {
			return strings.ReplaceAll(typ.Elem().String(), fmt.Sprintf("%s.", importPath), "")
		} else {
			return strings.ReplaceAll(typ.Elem().String(), importPath, strings.Split(importPath, "/")[len(strings.Split(importPath, "/"))-1])
		}
	}

	return typ.Elem().String()
}

func IsReturnSliceModel(results *types.Tuple) bool {
	return strings.HasPrefix(results.At(0).Type().String(), "[]")
}

func IsQueryIn(pt *parser.PartTree) bool {
	for _, v := range pt.Predicate.Nodes {
		for _, vv := range v.Children {
			if vv.Type.Name == "In" || vv.Type.Name == "NotIn" {
				return true
			}
		}
	}
	return false
}

func ParamName(params *types.Tuple, i int) string {
	return params.At(i).Name()
}

func GenInsertFieldBinding(params *types.Tuple, model *ModelSpecs) string {
	columns := GenInsertColumns(model)

	var paramName = ""
	for i := range params.Len() {
		var p = params.At(i).Type()
		if _, ok := params.At(i).Type().(*types.Pointer); ok {
			p = params.At(i).Type().(*types.Pointer).Elem()
		}

		if p.String() == model.Type.String() {
			paramName = params.At(i).Name()
		}
	}

	var s = &strings.Builder{}
	for i := range model.Struct.NumFields() {
		tag := reflect.StructTag(model.Struct.Tag(i))
		cn, _ := ParseTag(tag.Get("db"))
		for _, v := range columns {
			if v.Name == cn {
				s.WriteString(paramName)
				s.WriteString(".")
				s.WriteString(model.Struct.Field(i).Name())
				if i < model.Struct.NumFields()-1 {
					s.WriteString(", ")
				}
			}
		}
	}
	return s.String()
}

func GenUpdateFieldBinding(params *types.Tuple, model *ModelSpecs) string {
	columns := GenUpdateColumns(model)

	var paramName = ""
	for i := range params.Len() {
		var p = params.At(i).Type()
		if _, ok := params.At(i).Type().(*types.Pointer); ok {
			p = params.At(i).Type().(*types.Pointer).Elem()
		}

		if p.String() == model.Type.String() {
			paramName = params.At(i).Name()
		}
	}

	var pkFieldName = ""
	var s = &strings.Builder{}
	var numOfSet = 0
	for i := range model.Struct.NumFields() {
		tag := reflect.StructTag(model.Struct.Tag(i))
		columnName, opts := ParseTag(tag.Get("db"))
		if opts.Contains("pk") {
			pkFieldName = model.Struct.Field(i).Name()
		}

		flag := false
		for _, v := range columns {
			if v.Name == columnName {
				flag = true
				break
			}
		}

		if flag {
			s.WriteString(paramName)
			s.WriteString(".")
			s.WriteString(model.Struct.Field(i).Name())

			numOfSet++

			if numOfSet < len(columns) {
				s.WriteString(", ")
			}
		}
	}

	s.WriteString(", ")
	s.WriteString(paramName)
	s.WriteString(".")
	s.WriteString(pkFieldName)

	return s.String()
}
