package parser

import (
	"regexp"
	"strings"
)

var (
	queryPattern   = `Find|Read|Get|Query|Search|Stream`
	countPattern   = `Count`
	existsPattern  = `Exists`
	deletePattern  = `Delete|Remove`
	prefixTemplate = regexp.MustCompile(`^(` + queryPattern + `|` + countPattern + `|` + existsPattern + `|` + deletePattern + `)(\p{Lu}.*?)??By`)
	andPattern     = regexp.MustCompile(`And(\p{Lu})`)
)

type PartTree struct {
	Subject   *Subject
	Predicate *Predicate
}

func NewPartTree(source string) (*PartTree, error) {
	pt := &PartTree{}

	matches := prefixTemplate.FindAllStringSubmatch(source, -1)

	if matches == nil {
		pt.Subject = NewSubject("")
		predicate, err := NewPredicate(source)
		if err != nil {
			return nil, err
		}
		pt.Predicate = predicate
	} else {
		pt.Subject = NewSubject(matches[0][0])
		predicate, err := NewPredicate(source[len(matches[0][0]):])
		if err != nil {
			return nil, err
		}
		pt.Predicate = predicate
	}

	return pt, nil
}

type OrPart struct {
	Children []*Part
}

func NewOrPart(source string, isAlwaysIgnoreCase bool) *OrPart {
	orPart := &OrPart{}

	withoutAnd := andPattern.ReplaceAllString(source, " $1")
	for _, part := range strings.Split(withoutAnd, " ") {
		orPart.Children = append(orPart.Children, NewPart(part, isAlwaysIgnoreCase))
	}

	return orPart
}
