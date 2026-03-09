package files

import (
	"go-file-api/internal/storage"
	"go-file-api/internal/vault"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, vaultRepo *vault.Repository, minIOService *storage.MinIOService, jwtMiddleware fiber.Handler) {
	group := app.Group("/files/:vaultId", jwtMiddleware)

	group.Post("upload/*", vault.VaultAccessMiddleware(vaultRepo, vault.VaultRoleEditor), UploadFile(minIOService))
	group.Post("create/*", vault.VaultAccessMiddleware(vaultRepo, vault.VaultRoleEditor), CreateFile(minIOService))
	group.Get("download/*", vault.VaultAccessMiddleware(vaultRepo, vault.VaultRoleViewer), DownloadFile(minIOService))
	group.Get("metadata/*", vault.VaultAccessMiddleware(vaultRepo, vault.VaultRoleViewer), GetMetadata(minIOService))
	group.Get("list/*", vault.VaultAccessMiddleware(vaultRepo, vault.VaultRoleViewer), ListFiles(minIOService))
	group.Get("search", vault.VaultAccessMiddleware(vaultRepo, vault.VaultRoleViewer), SearchFiles(minIOService))
	group.Put("rename/*", vault.VaultAccessMiddleware(vaultRepo, vault.VaultRoleEditor), RenameFile(minIOService))
	group.Put("move/*", vault.VaultAccessMiddleware(vaultRepo, vault.VaultRoleEditor), MoveFile(minIOService))
	group.Delete("delete/*", vault.VaultAccessMiddleware(vaultRepo, vault.VaultRoleEditor), DeleteFile(minIOService))
}
