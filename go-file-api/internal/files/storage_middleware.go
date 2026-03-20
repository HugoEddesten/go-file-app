package files

import (
	"go-file-api/internal/locals"
	"go-file-api/internal/vault"

	"github.com/gofiber/fiber/v2"
)

// StorageLimitMiddleware rejects upload requests that would exceed the vault's storage limit.
// It must run after VaultAccessMiddleware (which sets vaultId in locals).
// Uses Content-Length as an upper-bound pre-check; precise accounting happens in the handler.
func StorageLimitMiddleware(vaultRepo *vault.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		contentLength := c.Request().Header.ContentLength()
		if contentLength <= 0 {
			return c.Next()
		}

		vaultId := locals.VaultId(c)
		storage, err := vaultRepo.GetVaultStorage(c.Context(), vaultId)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "could not check storage limits")
		}

		available := storage.LimitBytes - storage.UsedBytes
		if int64(contentLength) > available {
			return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
				"error":     "Storage limit exceeded",
				"limit":     storage.LimitBytes,
				"used":      storage.UsedBytes,
				"available": available,
			})
		}

		return c.Next()
	}
}
