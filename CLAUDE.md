# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A file management application with vault-based access control. Users create vaults (file containers), upload/manage files, and share vaults with other users via role-based permissions.

## Architecture

- **Backend:** Go 1.25 with Fiber v2 web framework, PostgreSQL 15 via pgx driver (raw SQL, no ORM or migrations)
- **Frontend (active):** React 19 + TypeScript + Vite in `go-file-ui-react/`
- **Frontend (legacy):** Svelte in `go-file-ui/` (not actively used)
- **Auth:** JWT tokens stored in HTTP-only cookies, bcrypt password hashing

### Backend Structure (`go-file-api/`)

Entry point: `cmd/api/main.go` — Fiber app on port 3000.

Modules in `internal/`:
- `auth/` — Register, login, /auth/me handlers. Registration creates a default vault.
- `jwt/` — JWT service (HS256, 24h expiry), protected route middleware via cookie validation
- `users/` — User repository (FindByEmail, Create)
- `vault/` — Vault CRUD, vault-user role assignment, VaultAccessMiddleware for path-based permissions. Roles: Owner(1), Admin(2), Editor(3), Viewer(4)
- `files/` — Upload, download, create, list, search, metadata. Files stored at `./uploads/{vaultId}/{path}`
- `db/` — pgxpool connection config (env vars: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
- `testutil/` — Test utilities (PostgresContainer with testcontainers)

Top-level `db/` package:
- `migrations/schema.sql` — Single source of truth for database schema
- `migrate.go` — InitSchema function (embeds and executes schema.sql)

Commands in `cmd/`:
- `api/` — Main application entry point (Fiber server on port 3000)
- `migrate/` — Standalone migration tool

Pattern: Handler functions (controllers) + Repository structs (data access) + Fiber middleware chain.

### Frontend Structure (`go-file-ui-react/`)

Feature-folder organization in `src/features/`:
- `auth/` — Login/Register with React Hook Form + Zod validation
- `vaults/` — Vault list with tabs (My Vaults / Shared)
- `home/` — File browser, preview, upload, drag-drop

Key libraries: Zustand (vaultId state), TanStack React Query (server state), Axios (API client with credentials), shadcn/ui + Radix UI (components), Tailwind CSS v4.

API client configured in `src/lib/api.ts` with base URL from `VITE_API_URL` env var.

Routing in `src/root/router/` with ProtectedRoute/PublicRoute wrappers.

## Build & Run Commands

### Database
```bash
# Start PostgreSQL + pgAdmin (from go-file-api/)
docker compose up -d
# pgAdmin at localhost:8080 (admin@local.com / admin123)
# Postgres at localhost:5432 (admin / admin123 / filedb)

# Initialize schema (run once or when schema changes)
go run ./cmd/migrate
```

### Backend
```bash
cd go-file-api
go run ./cmd/api/main.go              # Normal start
go run ./cmd/api/main.go --auto-migrate  # Start with auto-migration (dev)
```

### Frontend
```bash
cd go-file-ui-react
npm install
npm run dev        # Dev server (port 5173)
npm run build      # TypeScript check + Vite production build
npm run lint       # ESLint
```

## Database

### Schema Management

Schema defined in `db/migrations/schema.sql` (single source of truth). Uses embedded Go SQL for migrations.

**Migration Methods:**
1. **Standalone tool**: `go run ./cmd/migrate` (recommended for production)
2. **Auto-migrate**: Start API with `--auto-migrate` flag for development
3. **Programmatic**: `db.InitSchema(ctx, pool)` function used by tests

Safe to run multiple times (uses `CREATE TABLE IF NOT EXISTS`). Won't destroy existing data.

**Schema:**
- `users` (id, email UNIQUE, password_hash, created_at, updated_at)
- `vaults` (id, name, created_at, updated_at)
- `vault_users` (id, vault_id FK, user_id FK, path, role INTEGER 1-4, created_at, updated_at)
  - UNIQUE constraint on (vault_id, user_id, path)
  - Indexes: idx_vault_users_vault_id, idx_vault_users_user_id, idx_vault_users_vault_user

See `db/migrations/README.md` for detailed migration documentation.

## Testing

### Backend Tests

Located in `go-file-api/internal/*/` with `*_test.go` files.

**Test Infrastructure:**
- Framework: testify/suite for test organization, testify/assert and testify/require for assertions
- Database: Testcontainers with PostgreSQL 15-alpine for isolated test databases
- Utilities: `internal/testutil` package provides PostgresContainer setup, teardown, and table cleanup

**Test Coverage:**
- `users/repository_test.go` — Create, FindByEmail, duplicate handling, validation
- `vault/repository_test.go` — Create, AddUserToVault, GetVaultUsers, UpdateVaultUser, GetVault, GetVaultsForUser
- `vault/helpers_test.go` — Vault helper functions
- `vault/middleware_test.go` — VaultAccessMiddleware for role-based permissions
- `jwt/middleware_test.go` — JWT authentication middleware

**Running Tests:**
```bash
cd go-file-api
go test ./...                    # Run all tests
go test -v ./...                 # Verbose output
go test -race -coverprofile=coverage.out -covermode=atomic ./...  # With race detection and coverage
go tool cover -func=coverage.out  # Display coverage report
```

**CI/CD:**
GitHub Actions workflow (`.github/workflows/test.yml`) runs tests automatically on pull requests and pushes to main with race detection and coverage reporting.

### Frontend Tests

No test infrastructure exists yet in the frontend.

## Notes

- CORS is configured for `http://localhost:5173` only
- Frontend env: `go-file-ui-react/.env` sets `VITE_API_URL=http://127.0.0.1:3000/`
- Path aliases: `@/` maps to `src/` in the React app (configured in vite.config.ts and tsconfig)
