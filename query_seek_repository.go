package repository

import (
	"context"

	"github.com/lkebin/repository/filter"
	"github.com/lkebin/repository/pager"
)

// QuerySeekRepository provides seek (cursor-based) query capability.
// It is more efficient than offset-based pagination for large datasets.
type QuerySeekRepository[M any] interface {
	// QuerySeek queries records using cursor-based pagination.
	//
	// Parameters:
	//   - ctx: context for the query
	//   - size: number of records to return (page size)
	//   - p: SeekPager for ordering and cursor positioning
	//   - f: Filter for filtering records
	//
	// Returns:
	//   - []*M: list of records
	//   - int64: total count of records matching the filter
	//   - error: any error encountered
	QuerySeek(ctx context.Context, size int64, p pager.SeekPager, f filter.Filter) ([]*M, int64, error)
}
