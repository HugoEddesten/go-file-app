# Database Migrations

This directory contains the database schema that serves as the **single source of truth** for both development and production databases.

## Schema File

- **`schema.sql`** - Complete database schema with all tables, constraints, and indexes

The schema uses `CREATE TABLE IF NOT EXISTS`, making it safe to run multiple times without destroying existing data.

## Usage

### Option 1: Manual Migration (Recommended for Production)

Run the migration command to initialize or update the database schema:

```bash
# Using default configuration (from env vars or defaults)
go run ./cmd/migrate

# With custom connection parameters
go run ./cmd/migrate -host localhost -port 5432 -dbname filedb -user admin -password admin123
```

### Option 2: Auto-migrate on Application Startup (Development)

Start the application with the `--auto-migrate` flag:

```bash
go run ./cmd/api/main.go --auto-migrate
```

This will automatically initialize the schema before starting the server.

### Option 3: Docker Compose (Optional)

You can also add an init script to docker-compose.yml to run migrations automatically when the container starts.

## Configuration

The migration tools respect the same environment variables as the main application:

| Variable    | Default   |
|-------------|-----------|
| DB_HOST     | localhost |
| DB_PORT     | 5432      |
| DB_USER     | admin     |
| DB_PASSWORD | admin123  |
| DB_NAME     | filedb    |

## Safety

✅ **Safe to run multiple times** - Uses `IF NOT EXISTS` clauses
✅ **Won't destroy data** - Only creates tables if they don't exist
✅ **Won't modify existing tables** - Existing structure is preserved

⚠️ **Note**: This approach is suitable for initial schema creation and simple projects. For complex schema changes (adding columns, modifying constraints, etc.), consider using a proper migration tool like golang-migrate.

## Testing

The test suite (using testcontainers) uses the same `schema.sql` file, ensuring that tests run against the exact same schema as production.

## Current Schema

### Tables

1. **users**
   - id (SERIAL PRIMARY KEY)
   - email (VARCHAR, UNIQUE)
   - password_hash (VARCHAR)
   - created_at, updated_at (TIMESTAMP)

2. **vaults**
   - id (SERIAL PRIMARY KEY)
   - name (VARCHAR)
   - created_at, updated_at (TIMESTAMP)

3. **vault_users**
   - id (SERIAL PRIMARY KEY)
   - vault_id (INTEGER, FK to vaults)
   - user_id (INTEGER, FK to users)
   - path (VARCHAR)
   - role (INTEGER, 1-4: Owner/Admin/Editor/Viewer)
   - created_at, updated_at (TIMESTAMP)
   - UNIQUE constraint on (vault_id, user_id, path)

### Indexes

- `idx_vault_users_vault_id` - Fast lookups by vault
- `idx_vault_users_user_id` - Fast lookups by user
- `idx_vault_users_vault_user` - Fast lookups by vault+user combination
