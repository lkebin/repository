package generator

import (
	"go/types"
	"reflect"
)

func lookupColumns(model *ModelSpecs) []*Column {
	var columns []*Column
	for i := range model.Struct.NumFields() {
		tag := reflect.StructTag(model.Struct.Tag(i))
		columnName, _ := ParseTag(tag.Get("db"))
		if columnName != "" {
			columns = append(columns, &Column{
				Name:     columnName,
				Property: model.Struct.Field(i).Name(),
			})
		}
	}
	return columns
}

func lookupPkColumn(model *ModelSpecs) *Column {
	for i := range model.Struct.NumFields() {
		tag := reflect.StructTag(model.Struct.Tag(i))
		columnName, opts := ParseTag(tag.Get("db"))
		if columnName != "" && opts.Contains("pk") {
			return &Column{
				Name:     columnName,
				Property: model.Struct.Field(i).Name(),
			}
		}
	}
	return nil
}

func lookupColumnByProperty(property string, m *ModelSpecs) *Column {
	for i := range m.Struct.NumFields() {
		if m.Struct.Field(i).Name() == property {
			tag := reflect.StructTag(m.Struct.Tag(i))
			columnName, _ := ParseTag(tag.Get("db"))
			return &Column{
				Name:     columnName,
				Property: property,
			}
		}
	}
	return nil
}

func GenInsertColumns(model *ModelSpecs) []*Column {
	var columns []*Column
	for i := range model.Struct.NumFields() {
		tag := reflect.StructTag(model.Struct.Tag(i))
		columnName, opts := ParseTag(tag.Get("db"))
		if opts.Contains("autoincrement") {
			continue
		}
		if opts.Contains("unsafe") {
			continue
		}
		columns = append(columns, &Column{
			Name: columnName,
		})
	}

	return columns
}

func GenUpdateColumns(model *ModelSpecs) []*Column {
	var columns []*Column
	for i := range model.Struct.NumFields() {
		tag := reflect.StructTag(model.Struct.Tag(i))
		columnName, opts := ParseTag(tag.Get("db"))
		if opts.Contains("pk") {
			continue
		}
		if opts.Contains("autoincrement") {
			continue
		}
		if opts.Contains("unsafe") {
			continue
		}
		columns = append(columns, &Column{
			Name: columnName,
		})
	}

	return columns
}

func IsPkAutoIncrement(model *ModelSpecs) bool {
	for i := range model.Struct.NumFields() {
		tag := reflect.StructTag(model.Struct.Tag(i))
		_, opts := ParseTag(tag.Get("db"))
		if opts.Contains("pk") && opts.Contains("autoincrement") {
			return true
		}
	}
	return false
}

func PkFieldName(model *ModelSpecs) string {
	column := lookupPkColumn(model)
	if column == nil {
		return ""
	}
	return column.Property
}

func getModelFromNamed(n *types.Named) *ModelSpecs {
	specs := ParseModel([]string{n.Obj().Name()}, []string{n.Obj().Pkg().Path()}, []string{})
	if len(specs) > 0 {
		return &specs[0]
	}

	return nil
}

func lookupNamed(typ types.Type) *types.Named {
	switch t := typ.(type) {
	case *types.Pointer:
		return lookupNamed(t.Elem())
	case *types.Slice:
		return lookupNamed(t.Elem())
	case *types.Named:
		return t
	}

	return nil
}

func lookupPkgPath(typ types.Type) string {
	switch t := typ.(type) {
	case *types.Pointer:
		return lookupPkgPath(t.Elem())
	case *types.Slice:
		return lookupPkgPath(t.Elem())
	case *types.Named:
		if pkg := t.Obj().Pkg(); pkg != nil {
			return pkg.Path()
		}
	}
	return ""
}
