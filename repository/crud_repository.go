package repository

import (
	"context"
)

type CrudRepository[M any, K comparable] interface {
	Save(ctx context.Context, model *M) (*M, error)
	FindByID(ctx context.Context, id K) (*M, error)
	ExistsByID(ctx context.Context, id K) (bool, error)
	Delete(ctx context.Context, model *M) error
	FindAll(ctx context.Context) ([]*M, int64, error)
}
