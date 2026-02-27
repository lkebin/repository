package testdata

import (
	"context"

	"github.com/lkebin/repository"
)

type UserRepository interface {
	repository.CrudRepository[User, int64]

	FindByName(ctx context.Context, name string) (*User, error)
	FindByNameIsNull(ctx context.Context) ([]*User, error)
	FindByNameAndBirthday(ctx context.Context, name string, birthday string) ([]*User, error)
	FindByNameIsNullAndBirthday(ctx context.Context, birthday string) ([]*User, error)
	FindByNameIn(ctx context.Context, name []string) ([]*User, error)
	FindByBirthdayBetween(ctx context.Context, start string, end string) ([]*User, error)
	FindByNameOrderByNameAsc(ctx context.Context, name string) ([]*User, error)
	FindTop10ByName(ctx context.Context, name string) ([]*User, error)
	CountByName(ctx context.Context, name string) (int64, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	DeleteByName(ctx context.Context, name string) error
}
