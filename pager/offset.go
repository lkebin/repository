package pager

import (
	"fmt"
	"math"
	"strings"
)

const (
	Desc = "Desc"
	Asc  = "Asc"
)

type orderByPair struct {
	column    string
	direction string
}

type OffsetPager interface {
	SetOrder(orderByPairs ...string)
	SetPage(page int64)

	Build() (string, int64, error)
}

type offsetPagerImpl struct {
	page         int64
	orderByPairs []orderByPair
}

func NewOffsetPager() OffsetPager {
	return &offsetPagerImpl{
		page: 1,
	}
}

func (p *offsetPagerImpl) SetPage(page int64) {
	if page < 0 {
		page = 1
	}
	p.page = page
}

func (p *offsetPagerImpl) SetOrder(orderByPairs ...string) {
	pairs := orderByPairs
	if math.Mod(float64(len(orderByPairs)), 2) != 0 {
		pairs = append(pairs, Asc)
	}

	var curColumn string
	for k, v := range pairs {
		if math.Mod(float64(k+1), 2) == 0 {
			p.orderByPairs = append(p.orderByPairs, orderByPair{column: curColumn, direction: v})
		} else {
			curColumn = v
		}
	}
}

func (p *offsetPagerImpl) Build() (string, int64, error) {
	orderBy := make([]string, 0)
	for _, v := range p.orderByPairs {
		orderBy = append(orderBy, fmt.Sprintf("`%s` %s", v.column, strings.ToUpper(v.direction)))
	}

	return strings.Join(orderBy, ","), p.page, nil
}
