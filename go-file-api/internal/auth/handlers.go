package auth

import (
	"go-file-api/internal/jwt"
	"go-file-api/internal/users"
	"go-file-api/internal/vault"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(userRepo *users.Repository, vaultRepo *vault.Repository, jwtService *jwt.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		body := new(AuthRequest)
		if err := c.BodyParser(body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
		}

		if len(body.Email) < 5 || len(body.Password) < 5 {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Could not hash password")
		}

		userId, err := userRepo.Create(body.Email, string(hashed))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Email already exists")
		}

		_, err = vaultRepo.Create(ctx, "my_vault", userId)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Could not create vault")
		}

		token, _ := jwtService.GenerateToken(userId, body.Email)

		c.Cookie(&fiber.Cookie{
			Name:     "auth",
			Value:    token,
			HTTPOnly: true,
			Secure:   true,
			SameSite: fiber.CookieSameSiteNoneMode,
			Path:     "/",
		})

		return c.SendStatus(200)
	}
}

func Login(repo *users.Repository, jwtService *jwt.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := new(AuthRequest)

		if err := c.BodyParser(body); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
		}

		user, err := repo.FindByEmail(body.Email)

		if err != nil || user == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password))
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid email or password")
		}

		token, err := jwtService.GenerateToken(user.Id, user.Email)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "could not generate token"})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "auth",
			Value:    token,
			HTTPOnly: true,
			Secure:   true,
			SameSite: fiber.CookieSameSiteNoneMode,
			Path:     "/",
		})

		return c.SendStatus(200)
	}
}

func Me(repo *users.Repository, jwtService *jwt.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authCookie := c.Cookies("auth")

		claims, err := jwtService.ValidateToken(authCookie)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid cookie")
		}

		return c.JSON(fiber.Map{
			"userId": claims.UserId,
			"email":  claims.Email,
		})
	}
}
