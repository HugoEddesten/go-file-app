# Go File App

A file management application with vault-based access control. Users create vaults (file containers), upload and manage files, and share vaults with other users through role-based permissions.

## Tech Stack

**Backend:** Go with Fiber web framework, PostgreSQL via pgx
**Frontend:** React 19, TypeScript, Vite, Tailwind CSS, shadcn/ui
**Auth:** JWT (HTTP-only cookies) + bcrypt password hashing

## Getting Started

### Prerequisites

- Go 1.25+
- Node.js (for the React frontend)
- Docker & Docker Compose (for PostgreSQL)

### 1. Start the Database

```bash
cd go-file-api
docker compose up -d
```

This starts PostgreSQL on port 5432 and pgAdmin on port 8080.

| Service  | URL / Port        | Credentials                |
|----------|-------------------|----------------------------|
| Postgres | localhost:5432    | admin / admin123 / filedb  |
| pgAdmin  | localhost:8080    | admin@local.com / admin123 |

### 2. Initialize the Database Schema

The project includes a migration system with the schema defined in `go-file-api/db/migrations/schema.sql`.

**Option 1: Standalone Migration Tool (Recommended)**

```bash
cd go-file-api
go run ./cmd/migrate
```

**Option 2: Auto-migrate on Startup**

```bash
cd go-file-api
go run ./cmd/api/main.go --auto-migrate
```

The migration is **safe to run multiple times** (uses `IF NOT EXISTS`) and won't destroy existing data.

**Schema:**
- `users` — id, email (unique), password_hash, created_at, updated_at
- `vaults` — id, name, created_at, updated_at
- `vault_users` — id, vault_id, user_id, path, role (1-4), created_at, updated_at
  - Unique constraint on (vault_id, user_id, path)
  - Indexes for performance optimization
- `vault_invites` — id, vault_id, invited_by, email, role, path, token (unique), expires_at, accepted_at, created_at
  - Indexes on token and email for fast lookups

### 3. Run the Backend

```bash
cd go-file-api
go run ./cmd/api/main.go
```

The API starts on port 3000.

### 4. Run the Frontend

```bash
cd go-file-ui-react
npm install
npm run dev
```

The dev server starts on port 5173.

## Project Structure

```
go-file-api/          # Go backend
  cmd/api/            # Application entry point
  internal/
    auth/             # Register, login, session handlers
    jwt/              # JWT token service and middleware
    users/            # User repository
    vault/            # Vault CRUD, sharing, role-based access
    files/            # File upload, download, listing, search
    email/            # Email service (SMTP + Resend) and HTML templates
    invites/          # Vault invite repository and types
    db/               # Database connection config

go-file-ui-react/     # React frontend
  src/
    features/
      auth/           # Login and registration
      vaults/         # Vault listing and management
      home/           # File browser, preview, upload
    components/       # Shared UI components (shadcn/ui)
    contexts/         # Auth, drag-drop, and vault state
    lib/              # API client (Axios)
    root/router/      # Routing with protected/public routes
```

## Vault Roles

| Role   | Value | Description              |
|--------|-------|--------------------------|
| Owner  | 1     | Full control             |
| Admin  | 2     | Manage vault and users   |
| Editor | 3     | Upload and modify files  |
| Viewer | 4     | Read-only access         |

Users can be restricted to specific paths within a vault for fine-grained access control.

## API Overview

### Auth
- `POST /auth/register` — Create account (also creates a default vault and redeems any pending invites)
- `POST /auth/login` — Authenticate and receive JWT cookie
- `GET /auth/me` — Get current user info

### Vaults
- `GET /vault/get-user-vaults` — List vaults for authenticated user
- `POST /vault/create` — Create a new vault
- `GET /vault/get-vault/:vaultId` — Get vault details with users
- `POST /vault/assign-user/:vaultId` — Add a user to a vault (sends invite email if user doesn't exist)
- `PUT /vault/update-vault-user/:vaultId` — Update user role/path
- `GET /vault/invites/:vaultId` — List pending invites for a vault (admin only)

### Invites
- `GET /invites/:token` — Get invite info by token (public, used to pre-fill the register form)

### Files
- `POST /files/upload/:vaultId` — Upload a file
- `POST /files/create/:vaultId` — Create a file or directory
- `GET /files/download/:vaultId` — Download a file
- `GET /files/metadata/:vaultId` — Get file metadata
- `GET /files/list/:vaultId` — List directory contents
- `GET /files/search/:vaultId` — Search files by name

## Environment Variables

**Backend** (defaults used if not set):

| Variable       | Default                    | Description                              |
|----------------|----------------------------|------------------------------------------|
| DB_HOST        | localhost                  |                                          |
| DB_PORT        | 5432                       |                                          |
| DB_USER        | admin                      |                                          |
| DB_PASSWORD    | admin123                   |                                          |
| DB_NAME        | filedb                     |                                          |
| EMAIL_PROVIDER | smtp                       | `smtp` (dev/Mailhog) or `resend`         |
| EMAIL_FROM     | noreply@go-file-app.local  | Sender address                           |
| SMTP_HOST      | localhost                  | SMTP server host                         |
| SMTP_PORT      | 1025                       | SMTP server port (Mailhog default)       |
| RESEND_API_KEY | —                          | Required when `EMAIL_PROVIDER=resend`    |
| APP_URL        | http://localhost:5173       | Base URL used in invite email links      |

**Frontend** (`go-file-ui-react/.env`):

| Variable     | Default                    |
|--------------|----------------------------|
| VITE_API_URL | http://127.0.0.1:3000/     |

## Email

The `internal/email` package provides a provider-agnostic `EmailService` interface with two implementations selected via the `EMAIL_PROVIDER` env var:

- **`smtp`** (default) — connects to any SMTP server; use [Mailhog](https://github.com/mailhog/MailHog) locally on port 1025
- **`resend`** — sends via the [Resend](https://resend.com) API; requires `RESEND_API_KEY`

### Emails sent

| Trigger | Email |
|---|---|
| User registers | Welcome email to the new user |
| Vault shared with existing user | "You now have access to X" notification |
| Vault shared with unknown email | Invite email with a registration link |

### Vault invite flow

1. Admin shares a vault with an email that has no account
2. A `vault_invites` record is created with a unique token (expires in 7 days)
3. An invite email is sent containing a link to `/register/:token`
4. The frontend fetches `GET /invites/:token` to pre-fill and lock the email field
5. On registration, all pending invites for that email are automatically redeemed — `vault_users` rows are created and invites marked accepted

### Local development with Mailhog

```bash
docker run -d -p 1025:1025 -p 8025:8025 mailhog/mailhog
```

View captured emails at `http://localhost:8025`.

## Testing

### Backend Tests

The Go API has comprehensive repository and middleware tests using:
- **Framework:** testify/suite, testify/assert, testify/require
- **Database:** Testcontainers with PostgreSQL 15-alpine (isolated test databases)
- **Utilities:** `internal/testutil` package for container setup and cleanup

**Test Coverage:**
- User repository (Create, FindByEmail)
- Vault repository (Create, AddUserToVault, GetVaultUsers, UpdateVaultUser, GetVault, GetVaultsForUser)
- Vault helpers and middleware (role-based access control)
- JWT authentication middleware

**Running Tests:**

```bash
cd go-file-api
go test ./...                    # Run all tests
go test -v ./...                 # Verbose output
go test -race -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out  # Coverage report
```

**CI/CD:**

Tests run automatically via GitHub Actions on pull requests and pushes to main (`.github/workflows/test.yml`).

### Frontend Tests

No test infrastructure yet.
