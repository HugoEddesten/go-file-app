package vault

import (
	"go-file-api/internal/users"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(
	app *fiber.App,
	userRepo *users.Repository,
	vaultRepo *Repository,
	jwtMiddleware fiber.Handler,
) {
	group := app.Group("/vault", jwtMiddleware)

	group.Get("get-user-vaults", GetUserVaults(vaultRepo))
	group.Post("create", CreateVault(vaultRepo))

	group.Get("get-vault/:vaultId", VaultAccessMiddleware(vaultRepo, VaultRoleViewer), GetVault(vaultRepo))
	group.Post("assign-user/:vaultId", VaultAccessMiddleware(vaultRepo, VaultRoleAdmin), AssignUserToVault(vaultRepo))
	group.Put("update-vault-user/:vaultId", VaultAccessMiddleware(vaultRepo, VaultRoleAdmin), UpdateVaultUser(vaultRepo))
}
