package vault

import "github.com/gofiber/fiber/v2"

func VaultAccessMiddleware(
	vaultRepo *Repository,
	requiredRole VaultRole,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		userId := c.Locals("userId").(int) // set by auth middleware

		requestedPath, shouldValidatePath := ResolveVaultPath(c)

		vaultId, err := ResolveVaultId(c)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid vault id")
		}

		vaultUsers, err := vaultRepo.GetVaultUsers(ctx, vaultId, userId)
		if err != nil || len(vaultUsers) == 0 {
			return fiber.NewError(fiber.StatusForbidden)
		}

		allowed := false
		for _, vu := range vaultUsers {
			if vu.Role <= requiredRole && (!shouldValidatePath || pathAllowed(vu.Path, requestedPath)) {
				allowed = true
				c.Locals("vaultRole", vu.Role)
				if shouldValidatePath {
					c.Locals("requestedVaultPath", requestedPath)
				}
				break
			}
		}

		if !allowed {
			return fiber.NewError(fiber.StatusForbidden, "Insufficient access level")
		}

		c.Locals("vaultId", vaultId)

		return c.Next()
	}
}
