package testdata

type User struct {
	ID        int64  `db:"id,pk,autoincrement"`
	Name      string `db:"name"`
	Birthday  string `db:"birthday"`
	CreatedAt string `db:"created_at,unsafe"`
	UpdatedAt string `db:"updated_at"`
}

type UserUUID struct {
	ID        string `db:"id,pk,uuid"`
	Name      string `db:"name"`
	Birthday  string `db:"birthday"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}
