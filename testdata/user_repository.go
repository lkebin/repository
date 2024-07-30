package testdata

import (
	"context"
	repo "repository/repository"
)

type UserRepository interface {
	repo.CrudRepository[User, string]

	FindByName(ctx context.Context, name string) (*User, error)
}

// type UserRepository[X User, Y string] interface {
// 	repo.CrudRepository[X, Y]

// 	FindByName(ctx context.Context, name string) (*X, error)
// }
