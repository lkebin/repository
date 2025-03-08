package testdata

import (
	"context"
	repo "repository/repository"
)

type UserRepository interface {
	repo.CrudRepository[User, int64]

	FindByName(ctx context.Context, name string) (*User, error)
}
