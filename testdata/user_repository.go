package testdata

import (
	"context"
	repo "repository/repository"
)

type UserRepository interface {
	repo.CrudRepository[User, int64]

	FindByName(ctx context.Context, name string) (*User, error)
	FindByNameAndBirthdayIn(ctx context.Context, name string, birthday []string) ([]*User, error)
	FindByNameAndBirthday(ctx context.Context, name string, birthday string) ([]*User, error)
	FindByBirthdayBetween(ctx context.Context, start string, end string) ([]*User, error)
	FindByBirthdayIsBefore(ctx context.Context, date string) ([]*User, error)
}
