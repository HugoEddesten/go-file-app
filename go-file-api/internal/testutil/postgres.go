package testutil

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

//go:embed schema.sql
var schemaSQL string

// PostgresContainer holds the container instance and connection pool
type PostgresContainer struct {
	Container *postgres.PostgresContainer
	Pool      *pgxpool.Pool
}

// SetupPostgres starts a PostgreSQL container and initializes the schema
func SetupPostgres(ctx context.Context) (*PostgresContainer, error) {
	// Start PostgreSQL container
	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	// Get connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	// Create connection pool
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		pgContainer.Terminate(ctx)
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Initialize schema
	if _, err := pool.Exec(ctx, schemaSQL); err != nil {
		pool.Close()
		pgContainer.Terminate(ctx)
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return &PostgresContainer{
		Container: pgContainer,
		Pool:      pool,
	}, nil
}

// Teardown closes the connection pool and terminates the container
func (pc *PostgresContainer) Teardown(ctx context.Context) error {
	if pc.Pool != nil {
		pc.Pool.Close()
	}
	if pc.Container != nil {
		return pc.Container.Terminate(ctx)
	}
	return nil
}

// CleanupTables removes all data from tables (useful between tests)
func (pc *PostgresContainer) CleanupTables(ctx context.Context) error {
	_, err := pc.Pool.Exec(ctx, `
		TRUNCATE TABLE vault_users, vaults, users RESTART IDENTITY CASCADE;
	`)
	return err
}
