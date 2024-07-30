package parser

import (
	"errors"
	"regexp"
	"strings"
)

type Predicate struct {
	isAlwaysIgnoreCase bool
	nodes              []*OrPart
	orderBySource      *OrderBySource
}

var (
	allIgnoreCase = regexp.MustCompile("AllIgnor(ing|e)Case")
	orderBy       = "OrderBy"
)

func NewPredicate(predicate string) (*Predicate, error) {
	p := &Predicate{}

	parts := strings.Split(p.detectAndSetAllIgnoreCase(predicate), orderBy)
	if len(parts) > 2 {
		return nil, errors.New("OrderBy must not be used more than once in a query method name")
	}

	for _, v := range strings.Split(parts[0], "Or") {
		p.nodes = append(p.nodes, NewOrPart(v, p.isAlwaysIgnoreCase))
	}

	if len(parts) == 2 {
		orderBySource, err := NewOrderBySource(parts[1])
		if err != nil {
			return nil, err
		}

		p.orderBySource = orderBySource
	}

	return p, nil
}

func (p *Predicate) detectAndSetAllIgnoreCase(predicate string) string {
	indexes := allIgnoreCase.FindAllStringIndex(predicate, -1)
	if indexes != nil {
		p.isAlwaysIgnoreCase = true
		predicate = predicate[0:indexes[0][0]] + predicate[indexes[0][1]:]
	}

	return predicate
}
