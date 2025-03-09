package parser

import (
	"regexp"
	"strings"
)

var (
	queryPattern    = "Find|Read|Get|Query|Search|Stream"
	countPattern    = "Count"
	existsPattern   = "Exists"
	deletePattern   = "Delete|Remove"
	prefixTemplate  = regexp.MustCompile("^(" + queryPattern + "|" + countPattern + "|" + existsPattern + "|" + deletePattern + ")((\\p{Lu}.*?))??By")
	keywordTemplate = "(%s)(?=(\\p{Lu}|\\P{InBASIC_LATIN}))"
)

type PartTree struct {
	Subject   *Subject
	Predicate *Predicate
}

func NewPartTree(source string) (*PartTree, error) {
	pt := &PartTree{}

	matches := prefixTemplate.FindAllString(source, -1)

	if matches == nil {
		pt.Subject = NewSubject("")
		predicate, err := NewPredicate(source)
		if err != nil {
			return nil, err
		}
		pt.Predicate = predicate
	} else {
		pt.Subject = NewSubject(matches[0])
		predicate, err := NewPredicate(source[len(matches[0]):])
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

	split := strings.Split(source, "And")

	for _, part := range split {
		orPart.Children = append(orPart.Children, NewPart(part, isAlwaysIgnoreCase))
	}

	return orPart
}
