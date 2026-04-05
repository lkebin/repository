package pager

import (
	"fmt"
	"math"
	"strings"
)

const (
	OrderAsc            = "ASC"
	OrderDesc           = "DESC"
	DefaultLastIDColumn = "id"
)

// Next represents the cursor values for seek pagination.
// It contains the values of the order columns and the last ID column.
type Next []any

// NewNext creates a new Next cursor from the provided values.
func NewNext(values ...any) Next {
	n := make(Next, 0, len(values))
	for _, v := range values {
		n = append(n, v)
	}
	return n
}

// orderPair represents a single column ordering
type orderPair struct {
	column string
	order  string
}

// SeekPager provides cursor-based pagination (seek method).
// It is more efficient than offset-based pagination for large datasets.
type SeekPager interface {
	// Build returns the ORDER BY clause, WHERE clause for seek condition, and the values for the WHERE clause.
	Build() (orderBy string, seekWhere string, seekValues []any, err error)
	// SetOrder sets the order columns and directions (column1, order1, column2, order2, ...).
	// If order is missing for a column, ASC is used as default.
	SetOrder(orderPairs ...string)
	// SetLastIDColumn sets the column used as the last ID (default is "id").
	SetLastIDColumn(column string)
	// SetNext sets the cursor values for the next page.
	SetNext(next Next)
}

type seekPagerImpl struct {
	orderPairs   []orderPair
	lastIDColumn string
	next         Next
}

// NewSeekPager creates a new SeekPager with default settings.
func NewSeekPager() SeekPager {
	return &seekPagerImpl{
		orderPairs:   make([]orderPair, 0),
		lastIDColumn: DefaultLastIDColumn,
	}
}

func (p *seekPagerImpl) SetOrder(orderPairs ...string) {
	pairs := orderPairs
	// If odd number of arguments, append default ASC
	if math.Mod(float64(len(orderPairs)), 2) != 0 {
		pairs = append(pairs, OrderAsc)
	}

	p.orderPairs = make([]orderPair, 0)
	var curColumn string
	for k, v := range pairs {
		if math.Mod(float64(k+1), 2) == 0 {
			p.orderPairs = append(p.orderPairs, orderPair{column: curColumn, order: v})
		} else {
			curColumn = v
		}
	}
}

func (p *seekPagerImpl) SetLastIDColumn(column string) {
	p.lastIDColumn = column
}

func (p *seekPagerImpl) SetNext(next Next) {
	p.next = nil
	for _, v := range next {
		p.next = append(p.next, v)
	}
}

func (p *seekPagerImpl) Build() (string, string, []any, error) {
	var (
		orderByGroup []string
		seekWhere    = "true"
		seekValues   []any
		lastIDColumn = p.lastIDColumn
	)

	// Build ORDER BY clause from orderPairs
	for _, v := range p.orderPairs {
		if v.order == OrderAsc || v.order == OrderDesc {
			orderByGroup = append(orderByGroup, fmt.Sprintf("`%s` %s", v.column, v.order))
		} else {
			// For raw order expression
			orderByGroup = append(orderByGroup, fmt.Sprintf("`%s`", v.order))
		}
	}

	// Always append lastID column as tie-breaker
	orderByGroup = append(orderByGroup, fmt.Sprintf("`%s` %s", lastIDColumn, OrderAsc))
	orderBy := strings.Join(orderByGroup, ",")

	// Build seek WHERE clause if we have cursor values
	if p.next != nil && len(p.next) > 0 && len(p.next)-1 == len(p.orderPairs) {
		orPairs := make([]string, 0)

		for k, v := range p.orderPairs {
			operator := ">"
			if v.order == OrderDesc {
				operator = "<"
			}

			andPairs := make([]string, 0)
			if k > 0 {
				for i := 0; i < k; i++ {
					andPairs = append(andPairs, fmt.Sprintf("`%s` = ?", p.orderPairs[i].column))
					seekValues = append(seekValues, p.next[i])
				}
			}
			andPairs = append(andPairs, fmt.Sprintf("`%s` %s ?", v.column, operator))
			seekValues = append(seekValues, p.next[k])

			orPairs = append(orPairs, fmt.Sprintf("(%s)", strings.Join(andPairs, " AND ")))
		}

		// Final condition: all order columns equal and ID > lastID
		andPairs := make([]string, 0)
		for i := 0; i < len(p.orderPairs); i++ {
			andPairs = append(andPairs, fmt.Sprintf("`%s` = ?", p.orderPairs[i].column))
			seekValues = append(seekValues, p.next[i])
		}
		andPairs = append(andPairs, fmt.Sprintf("`%s` > ?", lastIDColumn))
		seekValues = append(seekValues, p.next[len(p.next)-1])

		orPairs = append(orPairs, fmt.Sprintf("(%s)", strings.Join(andPairs, " AND ")))

		seekWhere = fmt.Sprintf("(%s)", strings.Join(orPairs, " OR "))
	}

	return orderBy, seekWhere, seekValues, nil
}
