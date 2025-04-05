//go:generate ../repository -type=UserRepository,UserUuidRepository
package testdata

import (
	"context"

	"github.com/lkebin/repository"
)

type UserRepository interface {
	repository.CrudRepository[User, int64]

	FindByName(ctx context.Context, name string) (*User, error)
	FindByNameAndBirthdayIn(ctx context.Context, name string, birthday []string) ([]*User, error)
	FindByNameAndBirthday(ctx context.Context, name string, birthday string) ([]*User, error)
	FindByBirthdayBetween(ctx context.Context, start string, end string) ([]*User, error)
	FindByBirthdayIsBefore(ctx context.Context, date string) ([]*User, error)
	FindByBirthdayOrderByNameAsc(ctx context.Context, date string) ([]*User, error)
	FindByBirthdayOrderByNameAscBirthdayDesc(ctx context.Context, date string) ([]*User, error)
	FindTop20ByBirthday(ctx context.Context, date string) ([]*User, error)
}

type UserUuidRepository interface {
	repository.CrudRepository[UserUuid, repository.ID]

	FindByName(ctx context.Context, name string) (*UserUuid, error)
}
