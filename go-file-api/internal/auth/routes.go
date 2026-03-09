package auth

import (
	"go-file-api/internal/email"
	"go-file-api/internal/invites"
	"go-file-api/internal/jwt"
	"go-file-api/internal/users"
	"go-file-api/internal/vault"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(
	app *fiber.App,
	userRepo *users.Repository,
	vaultRepo *vault.Repository,
	inviteRepo *invites.Repository,
	emailSvc email.EmailService,
	jwtService *jwt.JWTService,
) {
	app.Post("/auth/register", Register(userRepo, vaultRepo, inviteRepo, emailSvc, jwtService))
	app.Post("/auth/reset-password", SendResetPasswordEmail(userRepo, emailSvc))
	app.Post("/auth/reset-password/:token", ResetPassword(userRepo))
	app.Post("/auth/login", Login(userRepo, jwtService))
	app.Post("/auth/logout", Logout())
	app.Get("/auth/me", Me(userRepo, jwtService))
}
