# Database Migration Guide

This project uses a simple, safe migration system with `schema.sql` as the single source of truth for database structure.

## 🎯 Quick Start

### For Development (First Time Setup)

1. **Start the database:**
   ```bash
   docker compose up -d
   ```

2. **Run migrations:**
   ```bash
   go run ./cmd/migrate
   ```

3. **Start the application:**
   ```bash
   go run ./cmd/api/main.go
   ```

### For Development (Auto-migrate on Startup)

Start the application with auto-migration enabled:

```bash
go run ./cmd/api/main.go --auto-migrate
```

This will automatically initialize the schema before starting the server.

## 📋 Migration Commands

### Run Migrations Manually

```bash
# Using default configuration (from environment variables or defaults)
go run ./cmd/migrate

# With custom connection parameters
go run ./cmd/migrate \
  -host localhost \
  -port 5432 \
  -dbname filedb \
  -user admin \
  -password admin123
```

### Build and Run the Migration Tool

```bash
# Build the migration binary
go build -o migrate ./cmd/migrate

# Run it
./migrate
```

## 🔧 Configuration

The migration system respects these environment variables:

| Variable    | Default   | Description          |
|-------------|-----------|----------------------|
| DB_HOST     | localhost | PostgreSQL host      |
| DB_PORT     | 5432      | PostgreSQL port      |
| DB_USER     | admin     | Database user        |
| DB_PASSWORD | admin123  | Database password    |
| DB_NAME     | filedb    | Database name        |

## ✅ Safety Features

- **Idempotent**: Uses `CREATE TABLE IF NOT EXISTS` - safe to run multiple times
- **Non-destructive**: Only creates tables if they don't exist
- **Preserves data**: Won't modify or delete existing tables or data
- **No dependencies**: Uses embedded SQL file, no external migration tool required

## 📁 File Structure

```
go-file-api/
├── db/
│   ├── migrate.go              # Migration logic
│   └── migrations/
│       ├── schema.sql          # Single source of truth for database schema
│       └── README.md           # Migration documentation
├── internal/
│   ├── db/                     # Database connection package
│   └── testutil/               # Test utilities (uses same schema)
└── cmd/
    ├── migrate/main.go         # Migration CLI tool
    └── api/main.go            # Main application (supports --auto-migrate flag)
```

## 🧪 Testing

Tests automatically use the same `schema.sql` file via testcontainers:

```bash
# Run all tests
go test ./internal/users ./internal/vault

# Run with verbose output
go test -v ./internal/users
```

Tests create temporary PostgreSQL containers, initialize the schema, and clean up automatically.

## 🚀 Production Deployment

### Option 1: Run Migration Command (Recommended)

```bash
# SSH into your server or run in your deployment pipeline
export DB_HOST=production-db-host
export DB_PASSWORD=secure-password
go run ./cmd/migrate
```

### Option 2: Docker Compose Init Script

Add to your `docker-compose.yml`:

```yaml
services:
  postgres:
    image: postgres:15
    volumes:
      - ./db/migrations/schema.sql:/docker-entrypoint-initdb.d/schema.sql
      - postgres_data:/var/lib/postgresql/data
```

This will run the schema on first container startup.

### Option 3: Auto-migrate in Startup Script

Your application startup script can include:

```bash
# Start with auto-migration
./api --auto-migrate
```

## ⚠️ Important Notes

1. **This is for initial schema setup** - The current system is perfect for creating tables initially and for small projects.

2. **For schema changes** (adding columns, modifying constraints, etc.), you have two options:
   - **Simple approach**: Manually ALTER tables in production, then update `schema.sql`
   - **Robust approach**: Migrate to a proper migration tool like [golang-migrate](https://github.com/golang-migrate/migrate) when your schema starts changing frequently

3. **Always backup** before running any schema changes in production

4. **Test migrations** in a staging environment first

## 🔄 Updating the Schema

When you need to add new tables or modify the schema:

1. Update `db/migrations/schema.sql`
2. Test locally: `go run ./cmd/migrate`
3. Run tests: `go test ./...`
4. Commit changes
5. Deploy and run migrations in production

## 📖 Examples

### First Time Setup

```bash
# Clone the repository
git clone <repo-url>
cd go-file-api

# Start PostgreSQL
docker compose up -d

# Initialize the database
go run ./cmd/migrate

# Start the application
go run ./cmd/api/main.go
```

### CI/CD Pipeline

```yaml
# Example GitHub Actions workflow
- name: Setup database
  run: docker compose up -d

- name: Wait for PostgreSQL
  run: sleep 5

- name: Run migrations
  run: go run ./cmd/migrate

- name: Run tests
  run: go test ./...
```

## 🆘 Troubleshooting

### "Connection refused"
- Ensure PostgreSQL is running: `docker compose ps`
- Check connection parameters match docker-compose.yml

### "Table already exists"
- This is normal! The migration is idempotent and can be run multiple times

### Tests failing with schema errors
- Ensure Docker is running (testcontainers requires Docker)
- Check that `db/migrations/schema.sql` is valid SQL

## 📚 Further Reading

- [PostgreSQL CREATE TABLE Documentation](https://www.postgresql.org/docs/current/sql-createtable.html)
- [Testcontainers for Go](https://golang.testcontainers.org/)
- [golang-migrate](https://github.com/golang-migrate/migrate) (for more complex migration needs)
