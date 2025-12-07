package auth

import (
	"go-file-api/internal/jwt"
	"go-file-api/internal/users"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(app *fiber.App, pool *pgxpool.Pool, jwtService *jwt.JWTService) {
	repo := &users.Repository{DB: pool}

	app.Post("/auth/register", Register(repo, jwtService))
	app.Post("/auth/login", Login(repo, jwtService))
	app.Get("/auth/me", Me(repo, jwtService))
}
