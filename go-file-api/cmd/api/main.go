package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"go-file-api/db"
	"go-file-api/internal/auth"
	internaldb "go-file-api/internal/db"
	"go-file-api/internal/email"
	"go-file-api/internal/files"
	"go-file-api/internal/invites"
	"go-file-api/internal/jwt"
	"go-file-api/internal/storage"
	"go-file-api/internal/users"
	"go-file-api/internal/vault"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Parse command line flags
	autoMigrate := flag.Bool("auto-migrate", false, "Automatically run database migrations on startup")
	flag.Parse()

	godotenv.Load()

	database, err := internaldb.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Run migrations if requested
	if *autoMigrate {
		log.Println("Running database migrations...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := db.InitSchema(ctx, database.Pool); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("✅ Database schema initialized successfully")
	}

	minIOService, err := storage.InitializeStorage()

	jwtService := jwt.New(os.Getenv("JWT_SECRET"), "go-file-api", time.Hour*24)
	jwtMiddleware := jwt.Protected(jwtService)

	vaultRepo := vault.Repository{DB: database.Pool}
	userRepo := users.Repository{DB: database.Pool}
	inviteRepo := invites.Repository{DB: database.Pool}
	emailSvc := email.New()

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, https://files.eddesten.dev",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	auth.RegisterRoutes(app, &userRepo, &vaultRepo, &inviteRepo, emailSvc, jwtService)
	files.RegisterRoutes(app, &vaultRepo, minIOService, jwtMiddleware)
	vault.RegisterRoutes(app, &userRepo, &vaultRepo, &inviteRepo, emailSvc, jwtMiddleware)

	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
