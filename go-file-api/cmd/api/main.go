package main

import (
	"go-file-api/internal/auth"
	"go-file-api/internal/db"
	"go-file-api/internal/files"
	"go-file-api/internal/jwt"
	"go-file-api/internal/users"
	"go-file-api/internal/vault"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	jwtService := jwt.New("my_secret_key_123", "go-file-api", time.Hour*24)
	jwtMiddleware := jwt.Protected(jwtService)

	vaultRepo := vault.Repository{DB: database.Pool}
	userRepo := users.Repository{DB: database.Pool}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	auth.RegisterRoutes(app, &userRepo, &vaultRepo, jwtService)
	files.RegisterRoutes(app, &vaultRepo, jwtMiddleware)
	vault.RegisterRoutes(app, &userRepo, &vaultRepo, jwtMiddleware)

	os.MkdirAll("uploads", os.ModePerm)

	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
