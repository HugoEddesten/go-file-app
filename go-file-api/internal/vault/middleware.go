package vault

import (
	"context"
	"go-file-api/internal/locals"

	"github.com/gofiber/fiber/v2"
)

// VaultUserGetter defines the interface needed by VaultAccessMiddleware
type VaultUserGetter interface {
	GetVaultUsers(ctx context.Context, vaultId, userId int) ([]VaultUser, error)
}

func VaultAccessMiddleware(
	vaultRepo VaultUserGetter,
	requiredRole VaultRole,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		userId := locals.UserId(c) // set by auth middleware

		requestedPath, shouldValidatePath, err := ResolveVaultPath(c)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Unable to resolve path")
		}

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
				locals.SetVaultRole(c, int(vu.Role))
				if shouldValidatePath {
					locals.SetRequestedVaultPath(c, requestedPath)
				}
				break
			}
		}

		if !allowed {
			return fiber.NewError(fiber.StatusForbidden, "Insufficient access level")
		}

		locals.SetVaultId(c, vaultId)

		return c.Next()
	}
}
