# Test Utilities

This package provides utilities for integration testing with PostgreSQL using testcontainers.

## Overview

The testutil package automatically:
- Starts a PostgreSQL container for each test suite
- Loads the database schema from `schema.sql`
- Provides helper functions for test setup and teardown
- Cleans up containers after tests complete

## Usage

### Basic Test Setup

```go
import (
    "context"
    "testing"
    "go-file-api/internal/testutil"
    "github.com/stretchr/testify/suite"
)

type MyRepositoryTestSuite struct {
    suite.Suite
    pgContainer *testutil.PostgresContainer
    repo        *Repository
    ctx         context.Context
}

func (s *MyRepositoryTestSuite) SetupSuite() {
    s.ctx = context.Background()

    // Start PostgreSQL container with schema loaded
    pgContainer, err := testutil.SetupPostgres(s.ctx)
    require.NoError(s.T(), err)

    s.pgContainer = pgContainer
    s.repo = &Repository{DB: pgContainer.Pool}
}

func (s *MyRepositoryTestSuite) TearDownSuite() {
    // Clean up container after all tests
    if s.pgContainer != nil {
        s.pgContainer.Teardown(s.ctx)
    }
}

func (s *MyRepositoryTestSuite) SetupTest() {
    // Clean all data between tests
    s.pgContainer.CleanupTables(s.ctx)
}

func TestMyRepositoryTestSuite(t *testing.T) {
    suite.Run(t, new(MyRepositoryTestSuite))
}
```

## Requirements

- Docker must be running on your machine
- Go 1.25+
- Dependencies:
  - `github.com/testcontainers/testcontainers-go`
  - `github.com/testcontainers/testcontainers-go/modules/postgres`
  - `github.com/jackc/pgx/v5`
  - `github.com/stretchr/testify`

## Features

### SetupPostgres
Starts a PostgreSQL 15 container and initializes the schema.

```go
pgContainer, err := testutil.SetupPostgres(ctx)
```

### CleanupTables
Removes all data from tables while preserving the schema. Useful for running between tests.

```go
err := pgContainer.CleanupTables(ctx)
```

### Teardown
Closes the connection pool and terminates the container.

```go
err := pgContainer.Teardown(ctx)
```

## Schema Management

The database schema is loaded from `db/migrations/schema.sql` via the `db.InitSchema()` function. This ensures that tests use the exact same schema as production, maintaining consistency across environments.

## Performance Tips

1. **Reuse containers**: Use `SetupSuite` to start the container once for all tests in a suite
2. **Clean between tests**: Use `CleanupTables` in `SetupTest` instead of recreating the container
3. **Parallel tests**: Each test package can run in parallel with its own container

## Example Test Output

```
=== RUN   TestUsersRepositoryTestSuite
2026/02/08 22:01:29 🐳 Creating container for image postgres:15-alpine
2026/02/08 22:01:45 🔔 Container is ready
=== RUN   TestUsersRepositoryTestSuite/TestCreate_Success
=== RUN   TestUsersRepositoryTestSuite/TestFindByEmail_UserExists
--- PASS: TestUsersRepositoryTestSuite (19.15s)
    --- PASS: TestUsersRepositoryTestSuite/TestCreate_Success (0.02s)
    --- PASS: TestUsersRepositoryTestSuite/TestFindByEmail_UserExists (0.01s)
PASS
```
