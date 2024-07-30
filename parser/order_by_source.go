package parser

import (
	"errors"
	"regexp"
)

var (
	blockSplit            = "(?<=Asc|Desc)(?=\\p{Lu})"
	directionSplit        = regexp.MustCompile("(.+?)(Asc|Desc)?$")
	directionKeywords     = []string{"Asc", "Desc"}
	ErrInvalidOrderSyntax = errors.New("Invalid order-by clause syntax")
)

type Order struct {
	property  string
	direction string
}

type OrderBySource struct {
	orders []*Order
}

func NewOrderBySource(clause string) (*OrderBySource, error) {
	if clause == "" {
		return nil, nil
	}

	orderBySource := &OrderBySource{}

	split := regexp.MustCompile(blockSplit).Split(clause, -1)

	for _, part := range split {
		matcher := directionSplit.FindAllString(part, -1)
		if matcher != nil {
			return nil, ErrInvalidOrderSyntax
		}

		propertyString := matcher[1]
		directionString := matcher[2]

		if directionString == "" {
			return nil, ErrInvalidOrderSyntax
		}

		for _, v := range directionKeywords {
			if v == propertyString {
				return nil, ErrInvalidOrderSyntax
			}
		}

		orderBySource.orders = append(orderBySource.orders, &Order{property: propertyString, direction: directionString})
	}

	return orderBySource, nil
}
