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

### 2. Create the Database Schema

Connect to the database and create the required tables. There is no automated migration tool — the schema must be set up manually.

**Tables:**
- `users` — id, email, password_hash, created_at, updated_at
- `vaults` — id, name, created_at, updated_at
- `vault_users` — id, vault_id, user_id, path, role, created_at, updated_at

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
- `POST /auth/register` — Create account (also creates a default vault)
- `POST /auth/login` — Authenticate and receive JWT cookie
- `GET /auth/me` — Get current user info

### Vaults
- `GET /vault/get-user-vaults` — List vaults for authenticated user
- `POST /vault/create` — Create a new vault
- `GET /vault/get-vault/:vaultId` — Get vault details with users
- `POST /vault/assign-user/:vaultId` — Add a user to a vault
- `PUT /vault/update-vault-user/:vaultId` — Update user role/path

### Files
- `POST /files/upload/:vaultId` — Upload a file
- `POST /files/create/:vaultId` — Create a file or directory
- `GET /files/download/:vaultId` — Download a file
- `GET /files/metadata/:vaultId` — Get file metadata
- `GET /files/list/:vaultId` — List directory contents
- `GET /files/search/:vaultId` — Search files by name

## Environment Variables

**Backend** (defaults used if not set):

| Variable    | Default   |
|-------------|-----------|
| DB_HOST     | localhost |
| DB_PORT     | 5432      |
| DB_USER     | admin     |
| DB_PASSWORD | admin123  |
| DB_NAME     | filedb    |

**Frontend** (`go-file-ui-react/.env`):

| Variable     | Default                    |
|--------------|----------------------------|
| VITE_API_URL | http://127.0.0.1:3000/     |
