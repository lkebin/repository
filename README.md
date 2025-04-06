The SQL Query Methods Generator for Go
-------------------------------------

This is a simple SQL query method generator for Go. It is inspired by the Spring.

## Usage

Install the repository command

```bash
go install github.com/lkebin/repository/cmd/repository@latest
```

Declare a struct which works with [sqlx](https://github.com/jmoiron/sqlx), add additonal options for `db` tag. The following struct added `pk` and `autoincrement` options.

```go
package example

type User struct {
    Id        int64  `db:"id,pk,autoincrement"`
    Name      string `db:"name"`
    Birthday  string `db:"birthday"`
    CreatedAt string `db:"created_at"`
    UpdatedAt string `db:"updated_at"`
}
```

Declare an interface type within another file, named `user_repository.go`.

```go
//go:generate repository -type=UserRepository
package example

import (
    "context"
    "github.com/lkebin/repository"
)

type UserRepository interface {
    repository.CrudRepository[User, int64]

    FindByBirthday(ctx context.Context, birthday string) ([]*User, error)
}
```

Run `go generate`, a new file `user_repository_impl.go` generated, check the content of generated file.
