package testdata

import "github.com/lkebin/repository"

type User struct {
	Id        int64  `db:"id,pk,autoincrement"`
	Name      string `db:"name"`
	Birthday  string `db:"birthday"`
	CreatedAt string `db:"created_at,unsafe"`
	UpdatedAt string `db:"updated_at"`
}

type UserUuid struct {
	Id        repository.ID `db:"id,pk,uuid"`
	Name      string        `db:"name"`
	Birthday  string        `db:"birthday"`
	CreatedAt string        `db:"created_at"`
	UpdatedAt string        `db:"updated_at"`
}
