package generator

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/lkebin/repository/parser"
)

func GenSelectClause(m *ModelSpecs, isDistinct bool) string {
	columns := lookupColumns(m)
	var s = &strings.Builder{}
	s.WriteString("SELECT ")

	if isDistinct {
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

func GenCountClause(m *ModelSpecs) string {
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

func GenDeleteClause(tableName string) string {
	return fmt.Sprintf("DELETE FROM `%s`", tableName)
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

func GenWhereClausePredicate(pt *parser.PartTree, params *types.Tuple, m *ModelSpecs) (string, error) {
	numberOfParams := params.Len() - 1
	for _, v := range pt.Predicate.Nodes {
		for _, vv := range v.Children {
			numberOfParams -= vv.Type.NumberOfArguments
		}
	}
	if numberOfParams != 0 {
		return "", fmt.Errorf("number of params not match: %d", numberOfParams)
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
				return "", fmt.Errorf("column not found: %s", n.Property)
			}
			if numberOfChildren > 1 {
				s.WriteString("(")
			}
			s.WriteString("`")
			s.WriteString(column.Name)
			s.WriteString("`")
			op, err := parseOperator(n.Type)
			if err != nil {
				return "", fmt.Errorf("parse operator error: %w", err)
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
	return s.String(), nil
}

func GenOrderByClause(pt *parser.PartTree, m *ModelSpecs) (string, error) {
	var s = &strings.Builder{}
	if pt.Predicate.OrderBySource != nil {
		s.WriteString(" ORDER BY ")
		for k, v := range pt.Predicate.OrderBySource.Orders {
			column := lookupColumnByProperty(v.Property, m)
			if column == nil {
				return "", fmt.Errorf("column not found: %s", v.Property)
			}

			s.WriteString(column.Name)
			s.WriteString(" ")
			s.WriteString(strings.ToUpper(v.Direction))
			if k < len(pt.Predicate.OrderBySource.Orders)-1 {
				s.WriteString(", ")
			}
		}
	}

	return s.String(), nil
}

func GenLimitClause(pt *parser.PartTree) string {
	if pt.Subject.IsLimiting {
		return fmt.Sprintf(" LIMIT %d", pt.Subject.MaxResults)
	}
	return ""
}

func parseOperator(pt parser.PartType) (string, error) {
	switch pt.Name {
	case "Between":
		return " BETWEEN ? AND ?", nil
	case "IsNotNull":
		return " IS NOT NULL", nil
	case "IsNull":
		return " IS NULL", nil
	case "LessThan":
		return " < ?", nil
	case "LessThanEqual":
		return " <= ?", nil
	case "GreaterThan":
		return " > ?", nil
	case "GreaterThanEqual":
		return " >= ?", nil
	case "Before":
		return " < ?", nil
	case "After":
		return " > ?", nil
	case "NotLike":
		return " NOT LIKE ?", nil
	case "Like":
		return " LIKE ?", nil
	case "StartingWith":
		return " LIKE ?", nil
	case "EndingWith":
		return " LIKE ?", nil
	case "IsNotEmpty":
		return " IS NOT NULL", nil
	case "IsEmpty":
		return " IS NULL", nil
	case "NotContaining":
		return " NOT LIKE ?", nil
	case "Containing":
		return " LIKE ?", nil
	case "NotIn":
		return " NOT IN (?)", nil
	case "In":
		return " IN (?)", nil
	case "True":
		return " = TRUE", nil
	case "False":
		return " = FALSE", nil
	case "NegatingSimpleProperty":
		return " != ?", nil
	case "SimpleProperty":
		return " = ?", nil
	default:
		return "", fmt.Errorf("operator not implemented: %s", pt.Name)
	}
}
