package parser

import (
	"errors"
	"regexp"
	"strings"
)

type Predicate struct {
	IsAlwaysIgnoreCase bool
	Nodes              []*OrPart
	OrderBySource      *OrderBySource
}

var (
	allIgnoreCase = regexp.MustCompile(`AllIgnor(ing|e)Case`)
	orderBy       = `OrderBy`
	orPattern     = regexp.MustCompile(`Or(\p{Lu})`)
)

func NewPredicate(predicate string) (*Predicate, error) {
	p := &Predicate{}

	parts := strings.Split(p.detectAndSetAllIgnoreCase(predicate), orderBy)
	if len(parts) > 2 {
		return nil, errors.New("OrderBy must not be used more than once in a query method name")
	}

	withoutOr := orPattern.ReplaceAllString(parts[0], " $1")
	for _, v := range strings.Split(withoutOr, " ") {
		p.Nodes = append(p.Nodes, NewOrPart(v, p.IsAlwaysIgnoreCase))
	}

	if len(parts) == 2 {
		orderBySource, err := NewOrderBySource(parts[1])
		if err != nil {
			return nil, err
		}

		p.OrderBySource = orderBySource
	}

	return p, nil
}

func (p *Predicate) detectAndSetAllIgnoreCase(predicate string) string {
	indexes := allIgnoreCase.FindAllStringIndex(predicate, -1)
	if indexes != nil {
		p.IsAlwaysIgnoreCase = true
		predicate = predicate[0:indexes[0][0]] + predicate[indexes[0][1]:]
	}

	return predicate
}
