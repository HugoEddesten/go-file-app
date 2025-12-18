package jwt

import "github.com/gofiber/fiber/v2"

func Protected(jwtService *JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Cookies("auth")
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing auth cookie",
			})
		}

		claims, err := jwtService.ValidateToken(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		c.Locals("userId", claims.UserId)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}
