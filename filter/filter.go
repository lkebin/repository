package filter

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cast"
)

type options struct {
}

type Option func(*options)

type optionFunc func(*options)

type Rule struct {
	Column   string
	Type     string
	Operator string
	Inverse  bool
	Query    any

	// Rules is used for `or` and `and` type
	Rules []*Rule
}

type Filter interface {
	Rules() []*Rule
	SetRules(rules []*Rule)

	Search() string
	SetSearch(query string)

	In(column string, terms ...any)
	NotIn(column string, terms ...any)

	InJSON(column string, terms ...any)
	NotInJSON(column string, terms ...any)

	Is(column string, term any)
	Not(column string, term any)

	Exist(column string)
	NotExist(column string)

	Eq(column string, value any)
	Ne(column string, value any)
	Gt(column string, value any)
	Gte(column string, value any)
	Lt(column string, value any)
	Lte(column string, value any)
	Between(column string, start, end any)

	And(...*Rule)
	Or(...*Rule)

	Build() (string, []any, error)
}

type filterImpl struct {
	rules  []*Rule
	search string
}

func NewFilter(opts ...Option) Filter {
	o := &options{}
	for _, v := range opts {
		v(o)
	}

	f := &filterImpl{}

	return f
}

func (f *filterImpl) SetSearch(kw string) {
	f.search = kw
}

func (f *filterImpl) Search() string {
	return f.search
}

func (f *filterImpl) SetRules(rules []*Rule) {
	f.rules = rules
}

func (f *filterImpl) Rules() []*Rule {
	return f.rules
}

func And(rules ...*Rule) *Rule {
	return &Rule{
		Type:  "and",
		Rules: rules,
	}
}

func (f *filterImpl) And(rules ...*Rule) {
	f.rules = append(f.rules, And(rules...))
}

func Or(rules ...*Rule) *Rule {
	return &Rule{
		Type:  "or",
		Rules: rules,
	}
}

func (f *filterImpl) Or(rules ...*Rule) {
	f.rules = append(f.rules, Or(rules...))
}

func In(column string, terms ...any) *Rule {
	return &Rule{
		Column:   column,
		Type:     "terms",
		Operator: "in",
		Query:    terms,
	}
}

func (f *filterImpl) In(column string, terms ...any) {
	f.rules = append(f.rules, In(column, terms...))
}

func NotIn(column string, terms ...any) *Rule {
	return &Rule{
		Column:   column,
		Type:     "terms",
		Operator: "notin",
		Query:    terms,
	}
}

func (f *filterImpl) NotIn(column string, terms ...any) {
	f.rules = append(f.rules, NotIn(column, terms...))
}

func InJSON(column string, terms ...any) *Rule {
	return &Rule{
		Column:   column,
		Type:     "terms",
		Operator: "injson",
		Query:    terms,
	}
}

func (f *filterImpl) InJSON(column string, terms ...any) {
	f.rules = append(f.rules, InJSON(column, terms...))
}

func NotInJSON(column string, terms ...any) *Rule {
	return &Rule{
		Column:   column,
		Type:     "terms",
		Operator: "notinjson",
		Query:    terms,
	}
}

func (f *filterImpl) NotInJSON(column string, terms ...any) {
	f.rules = append(f.rules, NotInJSON(column, terms...))
}

func Is(column string, term any) *Rule {
	return &Rule{
		Column:   column,
		Type:     "term",
		Operator: "is",
		Query:    term,
	}
}

func (f *filterImpl) Is(column string, term any) {
	f.rules = append(f.rules, Is(column, term))
}

func Not(column string, term any) *Rule {
	return &Rule{
		Column:   column,
		Type:     "term",
		Operator: "not",
		Query:    term,
	}
}

func (f *filterImpl) Not(column string, term any) {
	f.rules = append(f.rules, Not(column, term))
}

func Exist(column string) *Rule {
	return &Rule{
		Column:  column,
		Type:    "exists",
		Inverse: false,
	}
}

func (f *filterImpl) Exist(column string) {
	f.rules = append(f.rules, Exist(column))
}

func NotExist(column string) *Rule {
	return &Rule{
		Column:  column,
		Type:    "exists",
		Inverse: true,
	}
}

func (f *filterImpl) NotExist(column string) {
	f.rules = append(f.rules, NotExist(column))
}

func Eq(column string, value any) *Rule {
	return &Rule{
		Column:   column,
		Type:     "term",
		Operator: "eq",
		Query:    value,
	}
}

func (f *filterImpl) Eq(column string, value any) {
	f.rules = append(f.rules, Eq(column, value))
}

func Ne(column string, value any) *Rule {
	return &Rule{
		Column:   column,
		Type:     "term",
		Operator: "ne",
		Query:    value,
	}
}

func (f *filterImpl) Ne(column string, value any) {
	f.rules = append(f.rules, Ne(column, value))
}

func Gt(column string, value any) *Rule {
	return &Rule{
		Column:   column,
		Type:     "range",
		Operator: "gt",
		Query:    value,
	}
}

func (f *filterImpl) Gt(column string, value any) {
	f.rules = append(f.rules, Gt(column, value))
}

func Gte(column string, value any) *Rule {
	return &Rule{
		Column:   column,
		Type:     "range",
		Operator: "gte",
		Query:    value,
	}
}

func (f *filterImpl) Gte(column string, value any) {
	f.rules = append(f.rules, Gte(column, value))
}

func Lt(column string, value any) *Rule {
	return &Rule{
		Column:   column,
		Type:     "range",
		Operator: "lt",
		Query:    value,
	}
}

func (f *filterImpl) Lt(column string, value any) {
	f.rules = append(f.rules, Lt(column, value))
}

func Lte(column string, value any) *Rule {
	return &Rule{
		Column:   column,
		Type:     "range",
		Operator: "lte",
		Query:    value,
	}
}

func (f *filterImpl) Lte(column string, value any) {
	f.rules = append(f.rules, Lte(column, value))
}

func Between(column string, start, end any) *Rule {
	return &Rule{
		Column:   column,
		Type:     "range",
		Operator: "between",
		Query: map[string]any{
			"gte": start,
			"lte": end,
		},
	}
}

func (f *filterImpl) Between(column string, start, end any) {
	f.rules = append(f.rules, Between(column, start, end))
}

func (f *filterImpl) parseRule(rule *Rule) (string, []any, error) {
	switch rule.Type {
	case "terms", "term", "range":
		if rule.Operator == "" {
			return "", nil, errors.New("invalid filter data, operator not a string")
		}

		if rule.Query == nil {
			return "", nil, errors.New("invalid filter data, missing query")
		}

		return f.whereBuild(rule.Column, rule.Operator, rule.Query)
	case "exists":
		if rule.Inverse {
			return f.whereBuild(rule.Column, "none", nil)
		} else {
			return f.whereBuild(rule.Column, "notnone", nil)
		}
	case "and":
		where := ""
		values := make([]any, 0)
		for _, v := range rule.Rules {
			ww, vv, err := f.parseRule(v)
			if err != nil {
				return ww, vv, err
			}

			where += (" AND " + ww)
			values = append(values, vv...)
		}
		return "(" + strings.TrimLeft(where, " AND") + ")", values, nil
	case "or":
		where := ""
		values := make([]any, 0)
		for _, v := range rule.Rules {
			ww, vv, err := f.parseRule(v)
			if err != nil {
				return ww, vv, err
			}

			where += (" OR " + ww)
			values = append(values, vv...)
		}
		return "(" + strings.TrimLeft(where, " OR ") + ")", values, nil
	}

	return "", nil, fmt.Errorf("invalid rule type")
}

func (f *filterImpl) Build() (where string, values []any, err error) {
	where = ""
	values = make([]any, 0)

	for _, v := range f.rules {
		ww, vv, err := f.parseRule(v)
		if err != nil {
			return where, values, fmt.Errorf("build filter error: %w", err)
		}

		if where != "" {
			where += " AND "
		}
		where += ww
		values = append(values, vv...)
	}

	if where == "" {
		where = "true"
	}

	return where, values, nil
}

func (f *filterImpl) whereBuild(column string, operator string, query any) (string, []any, error) {
	w := ""
	v := make([]any, 0)
	switch operator {
	case "injson":
		w = fmt.Sprintf(`(JSON_CONTAINS(%s, ?))`, column)
		b, err := json.Marshal(cast.ToSlice(query))
		if err != nil {
			return w, v, fmt.Errorf("parse injson error: %w", err)
		}
		v = append(v, string(b))
	case "notinjson":
		w = fmt.Sprintf(`(!JSON_CONTAINS(%s, ?))`, column)
		b, err := json.Marshal(cast.ToSlice(query))
		if err != nil {
			return w, v, fmt.Errorf("parse notinjson error: %w", err)
		}
		v = append(v, string(b))
	case "in":
		wx, values, err := sqlx.In(fmt.Sprintf(`(%s IN(?))`, column), cast.ToSlice(query))
		if err != nil {
			return w, v, fmt.Errorf("filter parse in error: %w", err)
		}
		w = wx
		v = append(v, values...)
	case "notin":
		wx, values, err := sqlx.In(fmt.Sprintf(`(%s NOT IN(?))`, column), cast.ToSlice(query))
		if err != nil {
			return w, v, fmt.Errorf("filter parse notin error: %w", err)
		}
		w = wx
		v = append(v, values...)
	case "none":
		w = fmt.Sprintf(`(%s IS NULL)`, column)
	case "notnone":
		w = fmt.Sprintf(`(%s IS NOT NULL)`, column)
	case "eq", "is":
		w = fmt.Sprintf(`(%s = ?)`, column)
		v = append(v, query)
	case "ne", "not":
		w = fmt.Sprintf(`(%s != ?)`, column)
		v = append(v, query)
	case "gt":
		w = fmt.Sprintf(`(%s > ?)`, column)
		v = append(v, query)
	case "gte":
		w = fmt.Sprintf(`(%s >= ?)`, column)
		v = append(v, column)
	case "lt":
		w = fmt.Sprintf(`(%s < ?)`, column)
		v = append(v, query)
	case "lte":
		w = fmt.Sprintf(`(%s <= ?)`, column)
		v = append(v, query)
	case "between":
		w = fmt.Sprintf(`(%s BETWEEN ? AND ?)`, column)
		qMap := query.(map[string]any)
		if q, ok := qMap["gte"]; ok {
			v = append(v, q)
		} else {
			return w, v, fmt.Errorf("filter parse between error, missing [%s] key", operator)
		}

		if q, ok := qMap["lte"]; ok {
			v = append(v, q)
		} else {
			return w, v, fmt.Errorf("filter parse between error, missing [%s] key", operator)
		}
	}
	return w, v, nil
}
