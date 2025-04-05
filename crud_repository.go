package repository

import (
	"context"
)

type CrudRepository[M any, K comparable] interface {
	Create(ctx context.Context, model *M) (*M, error)
	Update(ctx context.Context, model *M) error
	FindAll(ctx context.Context) ([]*M, error)
	FindById(ctx context.Context, id K) (*M, error)
	DeleteById(ctx context.Context, id K) error
	ExistsById(ctx context.Context, id K) (bool, error)
	Count(ctx context.Context) (int64, error)
}
