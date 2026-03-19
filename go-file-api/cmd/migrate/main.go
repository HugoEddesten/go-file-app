package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go-file-api/db"
	internaldb "go-file-api/internal/db"

	"github.com/joho/godotenv"
)

func main() {
	// Parse command line flags
	godotenv.Load()

	var (
		host     = flag.String("host", "", "Database host (default: from env or 'localhost')")
		port     = flag.String("port", "", "Database port (default: from env or '5432')")
		user     = flag.String("user", "", "Database user (default: from env or 'admin')")
		password = flag.String("password", "", "Database password (default: from env or 'admin123')")
		dbname   = flag.String("dbname", "", "Database name (default: from env or 'filedb')")
	)

	flag.Parse()

	// Load config (respects env vars and command line flags)
	cfg := internaldb.LoadConfig()

	// Override with command line flags if provided
	if *host != "" {
		cfg.Host = *host
	}
	if *port != "" {
		cfg.Port = *port
	}
	if *user != "" {
		cfg.User = *user
	}
	if *password != "" {
		cfg.Password = *password
	}
	if *dbname != "" {
		cfg.Name = *dbname
	}

	// Display connection info
	fmt.Printf("Connecting to PostgreSQL:\n")
	fmt.Printf("  Host: %s:%s\n", cfg.Host, cfg.Port)
	fmt.Printf("  Database: %s\n", cfg.Name)
	fmt.Printf("  User: %s\n", cfg.User)
	fmt.Println()

	// Connect to database
	database, err := internaldb.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	fmt.Println("Running database migrations...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := db.InitSchema(ctx, database.Pool); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("✅ Migration completed successfully!")
	fmt.Println()
	fmt.Println("Database schema is up to date.")
}
