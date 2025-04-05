package generator

import (
	_ "embed"
	"fmt"
	"go/format"
	"go/types"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/lkebin/repository/parser"
)

//go:embed templates/find.gotpl
var findTpl string

//go:embed templates/exists.gotpl
var existsTpl string

//go:embed templates/count.gotpl
var countTpl string

//go:embed templates/create.gotpl
var createTpl string

//go:embed templates/update.gotpl
var updateTpl string

//go:embed templates/delete.gotpl
var deleteTpl string

func GenerateRepositoryImplements(spec *RepositorySpecs) ([]byte, error) {
	implName := fmt.Sprintf("%sImpl", strings.ToLower(spec.Name[:1])+spec.Name[1:])
	var tpl = &strings.Builder{}
	tpl.WriteString(fmt.Sprintf("// Code generated by \"repository %s\"; DO NOT EDIT.\n", strings.Join(os.Args[1:], " ")))
	tpl.WriteString(fmt.Sprintf("package %s\n", spec.Pkg.Name()))
	tpl.WriteString("\n")
	// Imports
	tpl.WriteString("import (\n")
	imports := make(map[string]any)
	imports["github.com/jmoiron/sqlx"] = nil
	imports["database/sql"] = nil
	for _, v := range spec.Methods {
		for i := 0; i < v.Signature().Params().Len(); i++ {
			importPath := lookupPkgPath(v.Signature().Params().At(i).Type())
			if importPath != "" && importPath != spec.Pkg.Path() {
				imports[importPath] = nil
			}
		}
		for i := 0; i < v.Signature().Results().Len(); i++ {
			importPath := lookupPkgPath(v.Signature().Results().At(i).Type())
			if importPath != "" && importPath != spec.Pkg.Path() {
				imports[importPath] = nil
			}
		}
	}
	for k := range imports {
		tpl.WriteString(fmt.Sprintf("\t\"%s\"\n", k))
	}
	tpl.WriteString(")\n")
	tpl.WriteString("\n")

	tpl.WriteString(fmt.Sprintf("type %s struct {\n", implName))
	tpl.WriteString("\tdb sqlx.ExtContext\n")
	tpl.WriteString("}\n")
	tpl.WriteString("\n")

	// New
	tpl.WriteString(fmt.Sprintf("func New%s(db sqlx.ExtContext) %s {\n", spec.Name, spec.Name))
	tpl.WriteString(fmt.Sprintf("\treturn &%s{db: db}\n", implName))
	tpl.WriteString("}\n")
	tpl.WriteString("\n")

	modelNamed := lookupNamed(spec.EmbeddedTypeArgs[0])
	if modelNamed == nil {
		return nil, fmt.Errorf("parse model error")
	}

	model := getModelFromNamed(modelNamed)
	if model == nil {
		return nil, fmt.Errorf("get model from named error")
	}

	// Implements
	for _, v := range spec.Methods {
		fn, err := genFuncImpl(implName, model, v, spec)
		if err != nil {
			return nil, err
		}
		tpl.WriteString(fn)
		tpl.WriteString("\n")
	}

	return format.Source([]byte(tpl.String()))
}

func parseTuple(params *types.Tuple, spec *RepositorySpecs) string {
	var p = &strings.Builder{}
	for i := 0; i < params.Len(); i++ {
		// Name
		if params.At(i).Name() != "" {
			p.WriteString(params.At(i).Name())
			p.WriteString(" ")
		}

		// Type
		importPath := lookupPkgPath(params.At(i).Type())
		if importPath != "" {
			// If parameter pkg path is equal to repository pkg path, should not include pkg name in type
			if importPath == spec.Pkg.Path() {
				p.WriteString(strings.ReplaceAll(params.At(i).Type().String(), fmt.Sprintf("%s.", importPath), ""))
			} else {
				p.WriteString(strings.ReplaceAll(params.At(i).Type().String(), importPath, strings.Split(importPath, "/")[len(strings.Split(importPath, "/"))-1]))
			}
		} else {
			p.WriteString(params.At(i).Type().String())
		}

		// Comma between params
		if i < params.Len()-1 {
			p.WriteString(", ")
		}
	}
	return p.String()
}

type FuncImpl struct {
	Receiver      string
	Name          string
	TableName     string
	Params        *types.Tuple
	Results       *types.Tuple
	Model         *ModelSpecs
	Repository    *RepositorySpecs
	SelectColumns []*Column
	WhereColumns  []*Column
	PartTree      *parser.PartTree
}

type Column struct {
	Name     string
	Property string
}

func genFuncImpl(implName string, model *ModelSpecs, m *types.Func, spec *RepositorySpecs) (string, error) {
	var tpl = &strings.Builder{}
	fn := FuncImpl{
		Name:       m.Name(),
		Receiver:   fmt.Sprintf("*%s", implName),
		Params:     m.Signature().Params(),
		Results:    m.Signature().Results(),
		TableName:  ToSnakeCase(model.Name),
		Model:      model,
		Repository: spec,
	}

	pt, err := parser.NewPartTree(m.Name())
	if err != nil {
		return "", err
	}
	fn.PartTree = pt

	switch {
	case m.Name() == "Create":
		if err := template.Must(template.New("").Funcs(funcMap()).Parse(createTpl)).Execute(tpl, fn); err != nil {
			return "", fmt.Errorf("create template execute error: %w", err)
		}
	case m.Name() == "Update":
		pkColumn := lookupPkColumn(model)
		if pkColumn == nil {
			return "", fmt.Errorf("pk column not found")
		}
		fn.WhereColumns = []*Column{{Name: pkColumn.Name}}
		if err := template.Must(template.New("").Funcs(funcMap()).Parse(updateTpl)).Execute(tpl, fn); err != nil {
			return "", fmt.Errorf("update template execute error: %w", err)
		}
	case m.Name() == "Count":
		if err := template.Must(template.New("").Funcs(funcMap()).Parse(countTpl)).ExecuteTemplate(tpl, "Count", fn); err != nil {
			return "", fmt.Errorf("count all template execute error: %w", err)
		}
	case m.Name() == "FindById":
		pkColumn := lookupPkColumn(model)
		if pkColumn == nil {
			return "", fmt.Errorf("pk column not found")
		}
		pt, err := parser.NewPartTree("FindBy" + pkColumn.Property)
		if err != nil {
			return "", err
		}
		fn.PartTree = pt
		if err := template.Must(template.New("").Funcs(funcMap()).Parse(findTpl)).ExecuteTemplate(tpl, "FindBy", fn); err != nil {
			return "", fmt.Errorf("find by id template execute error: %w", err)
		}

	case m.Name() == "FindAll":
		if err := template.Must(template.New("").Funcs(funcMap()).Parse(findTpl)).ExecuteTemplate(tpl, "FindAll", fn); err != nil {
			return "", fmt.Errorf("find all template execute error: %w", err)
		}
	case m.Name() == "ExistsById":
		pkColumn := lookupPkColumn(model)
		if pkColumn == nil {
			return "", fmt.Errorf("pk column not found")
		}
		pt, err := parser.NewPartTree("ExistsBy" + pkColumn.Property)
		if err != nil {
			return "", err
		}
		fn.PartTree = pt
		if err := template.Must(template.New("").Funcs(funcMap()).Parse(existsTpl)).Execute(tpl, fn); err != nil {
			return "", fmt.Errorf("exists template execute error: %w", err)
		}

	case m.Name() == "DeleteById":
		pkColumn := lookupPkColumn(model)
		if pkColumn == nil {
			return "", fmt.Errorf("pk column not found")
		}
		pt, err := parser.NewPartTree("DeleteBy" + pkColumn.Property)
		if err != nil {
			return "", err
		}
		fn.PartTree = pt
		if err := template.Must(template.New("").Funcs(funcMap()).Parse(deleteTpl)).Execute(tpl, fn); err != nil {
			return "", fmt.Errorf("delete template execute error: %w", err)
		}

	case pt.Subject != &parser.Subject{}: // A valid subject
		if pt.Subject.IsCount {
			if err := template.Must(template.New("").Funcs(funcMap()).Parse(countTpl)).ExecuteTemplate(tpl, "CountBy", fn); err != nil {
				return "", fmt.Errorf("count template execute error: %w", err)
			}
		} else if pt.Subject.IsExists {
			if err := template.Must(template.New("").Funcs(funcMap()).Parse(existsTpl)).Execute(tpl, fn); err != nil {
				return "", fmt.Errorf("exists template execute error: %w", err)
			}
		} else if pt.Subject.IsDelete {
			if err := template.Must(template.New("").Funcs(funcMap()).Parse(deleteTpl)).Execute(tpl, fn); err != nil {
				return "", fmt.Errorf("delete template execute error: %w", err)
			}
		} else {
			if err := template.Must(template.New("").Funcs(funcMap()).Parse(findTpl)).ExecuteTemplate(tpl, "FindBy", fn); err != nil {
				return "", fmt.Errorf("find template execute error: %w", err)
			}
		}
	}

	return tpl.String(), nil
}

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

func GenOrderByClause(pt *parser.PartTree, m *ModelSpecs) string {
	var s = &strings.Builder{}
	if pt.Predicate.OrderBySource != nil {
		s.WriteString(" ORDER BY ")
		for k, v := range pt.Predicate.OrderBySource.Orders {
			column := lookupColumnByProperty(v.Property, m)
			if column == nil {
				log.Panicf("column not found: %s", v.Property)
			}

			s.WriteString(column.Name)
			s.WriteString(" ")
			s.WriteString(strings.ToUpper(v.Direction))
			if k < len(pt.Predicate.OrderBySource.Orders)-1 {
				s.WriteString(", ")
			}
		}
	}

	return s.String()
}

func GenLimitClause(pt *parser.PartTree) string {
	if pt.Subject.IsLimiting {
		return fmt.Sprintf(" LIMIT %d", pt.Subject.MaxResults)
	}
	return ""
}

func lookupColumnByProperty(property string, m *ModelSpecs) *Column {
	for i := 0; i < m.Struct.NumFields(); i++ {
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

func PkFieldName(model *ModelSpecs) string {
	column := lookupPkColumn(model)
	if column == nil {
		return ""
	}
	return column.Property
}

func ParamName(params *types.Tuple, i int) string {
	return params.At(i).Name()
}

func IsQueryIn(pt *parser.PartTree) bool {
	for _, v := range pt.Predicate.Nodes {
		for _, vv := range v.Children {
			if reflect.DeepEqual(vv.Type, parser.In) {
				return true
			}
		}
	}
	return false
}

func GenParams(params *types.Tuple, spec *RepositorySpecs) string {
	return parseTuple(params, spec)
}

func GenCtxParam(params *types.Tuple) string {
	for i := 0; i < params.Len(); i++ {
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

func GenInsertFieldBinding(params *types.Tuple, model *ModelSpecs) string {
	columns := GenInsertColumns(model)

	var paramName = ""
	for i := 0; i < params.Len(); i++ {
		var p = params.At(i).Type()
		if _, ok := params.At(i).Type().(*types.Pointer); ok {
			p = params.At(i).Type().(*types.Pointer).Elem()
		}

		if p.String() == model.Type.String() {
			paramName = params.At(i).Name()
		}
	}

	var s = &strings.Builder{}
	for i := 0; i < model.Struct.NumFields(); i++ {
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
	for i := 0; i < params.Len(); i++ {
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
	for i := 0; i < model.Struct.NumFields(); i++ {
		tag := reflect.StructTag(model.Struct.Tag(i))
		columnName, opts := ParseTag(tag.Get("db"))
		if opts.Contains("pk") {
			pkFieldName = model.Struct.Field(i).Name()
		}

		for _, v := range columns {
			if v.Name == columnName {
				s.WriteString(paramName)
				s.WriteString(".")
				s.WriteString(model.Struct.Field(i).Name())
				if i < len(columns) {
					s.WriteString(", ")
				}
			}
		}
	}

	// where clause binding
	s.WriteString(", ")
	s.WriteString(paramName)
	s.WriteString(".")
	s.WriteString(pkFieldName)

	return s.String()
}

func GenInsertColumns(model *ModelSpecs) []*Column {
	var columns []*Column
	for i := 0; i < model.Struct.NumFields(); i++ {
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
	for i := 0; i < model.Struct.NumFields(); i++ {
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

func GenInsertClause(tableName string, model *ModelSpecs) string {
	columns := GenInsertColumns(model)

	var s = &strings.Builder{}
	s.WriteString("INSERT INTO ")
	s.WriteString("`")
	s.WriteString(tableName)
	s.WriteString("` (")
	for i, v := range columns {
		s.WriteString("`")
		s.WriteString(v.Name)
		s.WriteString("`")
		if i < len(columns)-1 {
			s.WriteString(", ")
		}
	}
	s.WriteString(") VALUES (")
	for i := range columns {
		s.WriteString("?")
		if i < len(columns)-1 {
			s.WriteString(", ")
		}
	}
	s.WriteString(")")
	return s.String()
}

func GenDeleteClause(tableName string) string {
	return fmt.Sprintf("DELETE FROM `%s`", tableName)
}

func GenUpdateClause(tableName string, model *ModelSpecs) string {
	columns := GenUpdateColumns(model)

	var s = &strings.Builder{}
	s.WriteString("UPDATE ")
	s.WriteString("`")
	s.WriteString(tableName)
	s.WriteString("` ")
	s.WriteString("SET ")
	for i, v := range columns {
		s.WriteString("`")
		s.WriteString(v.Name)
		s.WriteString("` = ?")
		if i < len(columns)-1 {
			s.WriteString(", ")
		}
	}
	return s.String()
}

func GenSelectClause(pt *parser.PartTree, m *ModelSpecs) string {
	columns := lookupColumns(m)
	var s = &strings.Builder{}
	s.WriteString("SELECT ")

	if pt.Subject.IsDistinct {
		s.WriteString("DISTINCT ")
	}

	for i, v := range columns {
		s.WriteString("`")
		s.WriteString(v.Name)
		s.WriteString("`")
		if i < len(columns)-1 {
			s.WriteString(", ")
		}
	}

	return s.String()
}

func GenCountClause(pt *parser.PartTree, m *ModelSpecs) string {
	pkColumn := lookupPkColumn(m)
	var s = &strings.Builder{}
	s.WriteString("SELECT ")
	s.WriteString("COUNT(`")
	s.WriteString(pkColumn.Name)
	s.WriteString("`)")
	return s.String()
}

func GenFromClause(tableName string) string {
	return fmt.Sprintf(" FROM `%s`", tableName)
}

func parseOperator(pt parser.PartType) (string, error) {
	op := ""
	switch {
	case reflect.DeepEqual(pt, parser.Between):
		op = " BETWEEN ? AND ?"
	case reflect.DeepEqual(pt, parser.IsNotNull):
		op = " IS NOT NULL"
	case reflect.DeepEqual(pt, parser.IsNull):
		op = " IS NULL"
	case reflect.DeepEqual(pt, parser.LessThan):
		op = " < ?"
	case reflect.DeepEqual(pt, parser.LessThanEqual):
		op = " <= ?"
	case reflect.DeepEqual(pt, parser.GreaterThan):
		op = " > ?"
	case reflect.DeepEqual(pt, parser.GreaterThanEqual):
		op = " >= ?"
	case reflect.DeepEqual(pt, parser.Before):
		op = " < ?"
	case reflect.DeepEqual(pt, parser.After):
		op = " > ?"
	case reflect.DeepEqual(pt, parser.NotLike):
		op = " NOT LIKE ?"
	case reflect.DeepEqual(pt, parser.Like):
		op = " LIKE ?"
	case reflect.DeepEqual(pt, parser.StartingWith):
		op = " LIKE ?"
	case reflect.DeepEqual(pt, parser.EndingWith):
		op = " LIKE ?"
	case reflect.DeepEqual(pt, parser.IsNotEmpty):
		op = " IS NOT NULL"
	case reflect.DeepEqual(pt, parser.IsEmpty):
		op = " IS NULL"
	case reflect.DeepEqual(pt, parser.NotContaining):
		op = " NOT LIKE ?"
	case reflect.DeepEqual(pt, parser.Containing):
		op = " LIKE ?"
	case reflect.DeepEqual(pt, parser.NotIn):
		op = " NOT IN (?)"
	case reflect.DeepEqual(pt, parser.In):
		op = " IN (?)"
	case reflect.DeepEqual(pt, parser.NegatingSimpleProperty):
		op = " != ?"
	case reflect.DeepEqual(pt, parser.SimpleProperty):
		op = " = ?"
	default:
		return "", fmt.Errorf("operator not implemented")
	}

	return op, nil
}

func GenWhereClause(columns []*Column) string {
	var s = &strings.Builder{}
	s.WriteString(" WHERE ")
	for i, v := range columns {
		s.WriteString("`")
		s.WriteString(v.Name)
		s.WriteString("` = ?")
		if i < len(columns)-1 {
			s.WriteString(" AND ")
		}
	}
	return s.String()
}

func GenWhereClausePredicate(pt *parser.PartTree, params *types.Tuple, m *ModelSpecs) string {
	// validate number of params
	numberOfParams := params.Len() - 1
	for _, v := range pt.Predicate.Nodes {
		for _, vv := range v.Children {
			numberOfParams -= vv.Type.NumberOfArguments
		}
	}
	if numberOfParams != 0 {
		log.Panicf("number of params not match: %d", numberOfParams)
	}
	columns := lookupColumns(m)
	var s = &strings.Builder{}
	s.WriteString(" WHERE ")
	numberOfNodes := len(pt.Predicate.Nodes)
	for k, v := range pt.Predicate.Nodes {
		if numberOfNodes > 1 {
			s.WriteString("(")
		}
		numberOfChildren := len(v.Children)
		for i, n := range v.Children {
			var column *Column
			for _, c := range columns {
				if strings.ToLower(c.Property) == strings.ToLower(n.Property) {
					column = c
				}
			}
			if column == nil {
				log.Panicf("column not found: %s", n.Property)
			}
			if numberOfChildren > 1 {
				s.WriteString("(")
			}
			s.WriteString("`")
			s.WriteString(column.Name)
			s.WriteString("`")
			op, err := parseOperator(n.Type)
			if err != nil {
				log.Panicf("parse operator error: %s", err)
			}
			s.WriteString(op)
			if numberOfChildren > 1 {
				s.WriteString(")")
			}
			if i < len(v.Children)-1 {
				s.WriteString(" AND ")
			}
		}
		if numberOfNodes > 1 {
			s.WriteString(")")
		}

		if k < len(pt.Predicate.Nodes)-1 {
			s.WriteString(" OR ")
		}
	}
	return s.String()
}

func GenVarBinding(params *types.Tuple) string {
	var s = &strings.Builder{}
	// always skip the first params, it's context.Context
	for i := 1; i < params.Len(); i++ {
		s.WriteString(params.At(i).Name())
		if i < params.Len()-1 {
			s.WriteString(", ")
		}
	}
	return s.String()
}

func IsReturnSliceModel(results *types.Tuple) bool {
	return strings.HasPrefix(results.At(0).Type().String(), "[]")
}

func IsPkAutoIncrement(model *ModelSpecs) bool {
	for i := 0; i < model.Struct.NumFields(); i++ {
		tag := reflect.StructTag(model.Struct.Tag(i))
		_, opts := ParseTag(tag.Get("db"))
		if opts.Contains("pk") && opts.Contains("autoincrement") {
			return true
		}
	}
	return false
}

func GenResultModel(results *types.Tuple) string {
	if IsReturnSliceModel(results) {
		if importPath := lookupPkgPath(results.At(0).Type()); importPath != "" {
			return strings.ReplaceAll(results.At(0).Type().String(), importPath, strings.Split(importPath, "/")[len(strings.Split(importPath, "/"))-1])
		}
		return results.At(0).Type().String()
	}

	typ, ok := results.At(0).Type().(*types.Pointer)
	if !ok {
		return results.At(0).Type().String()
	}

	if importPath := lookupPkgPath(typ.Elem()); importPath != "" {
		return strings.ReplaceAll(typ.Elem().String(), importPath, strings.Split(importPath, "/")[len(strings.Split(importPath, "/"))-1])
	}

	return typ.Elem().String()
}

func lookupColumns(model *ModelSpecs) []*Column {
	var columns []*Column
	for i := 0; i < model.Struct.NumFields(); i++ {
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
	for i := 0; i < model.Struct.NumFields(); i++ {
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

type TagOptions string

func ParseTag(tag string) (string, TagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], TagOptions(tag[idx+1:])
	}
	return tag, TagOptions("")
}

func (t TagOptions) Get(optionName string) string {
	if len(t) == 0 {
		return ""
	}

	s := string(t)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}

		j := strings.Index(s, "=")
		if j >= 0 {
			k, v := s[:j], s[j+1:]
			if k == optionName {
				return v
			}
		}

		s = next
	}

	return ""
}

func (t TagOptions) Contains(optionName string) bool {
	if len(t) == 0 {
		return false
	}
	s := string(t)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}

// ToSnakeCase convert string from camel case to snake case
func ToSnakeCase(camelCase string) string {
	re := regexp.MustCompile(`(.)([A-Z][a-z]+)`)
	s1 := re.ReplaceAllString(camelCase, "${1}_${2}")
	re1 := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	s2 := re1.ReplaceAllString(s1, "${1}_${2}")

	return strings.ToLower(s2)
}
