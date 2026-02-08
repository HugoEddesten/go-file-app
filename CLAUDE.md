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
```

### Backend
```bash
cd go-file-api
go run ./cmd/api/main.go
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

No migration tool — tables must be created manually. Implied schema:
- `users` (id, email, password_hash, created_at, updated_at)
- `vaults` (id, name, created_at, updated_at)
- `vault_users` (id, vault_id, user_id, path, role, created_at, updated_at)

## Notes

- CORS is configured for `http://localhost:5173` only
- Frontend env: `go-file-ui-react/.env` sets `VITE_API_URL=http://127.0.0.1:3000/`
- Path aliases: `@/` maps to `src/` in the React app (configured in vite.config.ts and tsconfig)
- No test infrastructure exists yet in either backend or frontend
