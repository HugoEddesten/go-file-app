package files

import (
	"go-file-api/internal/vault"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, vaultRepo *vault.Repository, jwtMiddleware fiber.Handler) {
	group := app.Group("/files/:vaultId", jwtMiddleware)

	group.Post("upload/*", vault.VaultAccessMiddleware(vaultRepo, vault.VaultRoleEditor), UploadFile())
	group.Post("create/*", vault.VaultAccessMiddleware(vaultRepo, vault.VaultRoleEditor), CreateFile())
	group.Get("download/*", vault.VaultAccessMiddleware(vaultRepo, vault.VaultRoleViewer), DownloadFile())
	group.Get("metadata/*", vault.VaultAccessMiddleware(vaultRepo, vault.VaultRoleViewer), GetMetadata())
	group.Get("list/*", vault.VaultAccessMiddleware(vaultRepo, vault.VaultRoleViewer), ListFiles())
	group.Get("search", vault.VaultAccessMiddleware(vaultRepo, vault.VaultRoleViewer), SearchFiles())
}
