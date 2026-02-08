package db

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/schema.sql
var schemaSQL string

// InitSchema runs the database schema initialization.
// It uses CREATE TABLE IF NOT EXISTS, so it's safe to run multiple times.
// It will create tables if they don't exist, but won't modify existing tables.
func InitSchema(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, schemaSQL); err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}
	return nil
}

// MustInitSchema is like InitSchema but panics on error.
// Useful for application startup where schema must be initialized.
func MustInitSchema(ctx context.Context, pool *pgxpool.Pool) {
	if err := InitSchema(ctx, pool); err != nil {
		panic(fmt.Sprintf("schema initialization failed: %v", err))
	}
}
