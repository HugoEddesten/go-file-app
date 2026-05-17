package auth

import (
	"go-file-api/internal/invites"
	"go-file-api/internal/jwt"
	"go-file-api/internal/users"
	"go-file-api/internal/vault"

	"eddesten-mail/client"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(
	app *fiber.App,
	userRepo *users.Repository,
	vaultRepo *vault.Repository,
	inviteRepo *invites.Repository,
	emailClient *client.Client,
	jwtService *jwt.JWTService,
) {
	app.Post("/auth/register", Register(userRepo, vaultRepo, inviteRepo, emailClient, jwtService))
	app.Post("/auth/reset-password", SendResetPasswordEmail(userRepo, emailClient))
	app.Post("/auth/reset-password/:token", ResetPassword(userRepo))
	app.Post("/auth/login", Login(userRepo, jwtService))
	app.Post("/auth/logout", Logout())
	app.Get("/auth/me", Me(userRepo, jwtService))
}
