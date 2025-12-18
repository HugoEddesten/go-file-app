package auth

import (
	"go-file-api/internal/jwt"
	"go-file-api/internal/users"
	"go-file-api/internal/vault"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(
	app *fiber.App,
	userRepo *users.Repository,
	vaultRepo *vault.Repository,
	jwtService *jwt.JWTService,
) {
	app.Post("/auth/register", Register(userRepo, vaultRepo, jwtService))
	app.Post("/auth/login", Login(userRepo, jwtService))
	app.Get("/auth/me", Me(userRepo, jwtService))
}
