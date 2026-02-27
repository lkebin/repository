# Agent Instructions for Repository Project

This document provides context and guidelines for AI agents working on this Go codebase.

## Project Overview
This project (`github.com/lkebin/repository`) is a Go library and code generator that implements the Repository pattern, likely inspired by Spring Data JPA. It parses interface definitions and generates implementation code for CRUD and query operations using `sqlx`.

## Build & Test Commands

### Build
To build the CLI tool:
```bash
go build -o repository ./cmd/repository
```

To run the example generation (requires the binary to be built or run via `go run`):
```bash
# Using the Makefile shortcut
make example

# Or manually
cd example && go generate .
```

### Testing
Run all tests:
```bash
go test ./...
```

Run tests for a specific package:
```bash
go test ./parser/...
```

Run a single specific test case:
```bash
# Format: go test -v <package_path> -run <TestNameRegex>
go test -v ./parser -run TestParsePredicate
```

### Linting & Formatting
Ensure code is formatted using standard Go tools:
```bash
go fmt ./...
go vet ./...
```

## Code Style & Conventions

### Formatting
- **Standard**: Strictly adhere to `gofmt` standards.
- **Indentation**: Use tabs for indentation, not spaces.

### Imports
Group imports into three blocks separated by a newline:
1. Standard library (e.g., `"fmt"`, `"os"`)
2. Third-party libraries (e.g., `"github.com/jmoiron/sqlx"`)
3. Local/Internal project imports (e.g., `"github.com/lkebin/repository/generator"`)

```go
import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/lkebin/repository/parser"
)
```

### Naming Conventions
- **Go Types/Functions**: Use `PascalCase` for exported identifiers and `camelCase` for unexported ones.
- **Generated Files**: The project conventions use snake_case for generated files (e.g., `user_repository_impl.go`), usually suffixed with `_impl.go`.
- **Database Tags**: Struct tags for database mapping use `db:"column_name"`.
  - Options include: `pk` (primary key), `autoincrement`, `unsafe` (ignore in inserts).

### Error Handling
- **Pattern**: Use the standard `if err != nil` pattern.
- **Return**: Errors should almost always be the last return value.
- **Context**: Wrap errors when propagating up call stacks to add context if helpful, but keep it simple.
- **Panic**: Avoid panics in library code. Only panic in the CLI `main` package if a fatal error occurs during startup or configuration. Internal generator logic occasionally uses `log.Panicf` for unreachable states, but prefer returning errors.

### Type Safety
- The project uses Go generics (Go 1.18+) extensively (e.g., `CrudRepository[Entity, ID]`). Ensure strictly typed interfaces are maintained.

## Project Structure
- `cmd/repository/`: Entry point for the CLI tool. Parses flags and orchestrates generation.
- `generator/`: Core logic for generating Go source code.
  - `generator.go`: Main generation logic, template execution, and helper functions (`funcMap`).
  - `templates/`: Embedded `.gotpl` files for different method types (Find, Create, etc.).
- `parser/`: Logic for parsing Go interfaces and method names (DSL parsing).
  - `part_tree.go`: Breaks down method names into Subject and Predicate.
  - `predicate.go`: Parses the `By...` part of the method name.
- `example/`: Example usage and integration tests.
- `repository.go` / `crud_repository.go`: Core interfaces defined for the library consumers.

## Architecture & Logic Details

### Template System
The generator uses `text/template`. Templates are embedded as string variables in `generator.go` (e.g., `findTpl`).
- **FuncMap**: A custom `funcMap` in `generator.go` provides SQL generation helpers:
  - `SelectClause`, `WhereClause`, `OrderByClause`: Generate SQL fragments based on the parsed model.
  - `VarBinding`, `Params`, `CtxParam`: Helper functions for generating Go method signatures.

### DSL Parsing (`parser` package)
The project heavily relies on parsing method names to infer SQL queries.
- **Regex Patterns**: `parser/part_tree.go` defines patterns like `queryPattern`, `countPattern` to identify the operation type.
- **PartTree**: Represents the structure of a parsed method name.
- **Subject**: The action (Find, Count, Delete, Exists).
- **Predicate**: The conditions (e.g., `ByNameAndAge`).
- **Operators**: Handled in `generator.go`'s `parseOperator` function (maps DSL operators like `Between`, `LessThan` to SQL).

## Workflow Rules
1. **Generation**: When modifying logic that affects code generation, always run the `example` generation to verify the output is valid Go code.
2. **Backward Compatibility**: This is a library; avoid breaking public interfaces (`CrudRepository`, `QueryRepository`) unless necessary.
3. **No Magic**: Prefer explicit code generation over runtime reflection where possible, aligning with the project's design philosophy.
4. **Testing**: Add unit tests in `parser/` when adding new DSL keywords. Add integration tests in `example/` when adding new generation features.
