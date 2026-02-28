package testdata

import (
	"context"

	"github.com/lkebin/repository"
)

type NoPkRepository interface {
	repository.CrudRepository[NoPkModel, int64]

	FindByName(ctx context.Context, name string) (*NoPkModel, error)
}
