package repository

import (
	"context"
)

type CrudRepository[M any, K comparable] interface {
	Create(ctx context.Context, model *M) (*M, error)
	Update(ctx context.Context, model *M) error
	Find(ctx context.Context, id K) (*M, error)
	Delete(ctx context.Context, id K) error
}
