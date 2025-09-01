package repository

import (
	"context"

	"github.com/lkebin/repository/filter"
	"github.com/lkebin/repository/pager"
)

type QueryRepository[M any] interface {
	Query(ctx context.Context, size int64, p pager.OffsetPager, f filter.Filter) ([]*M, int64, error)
}
