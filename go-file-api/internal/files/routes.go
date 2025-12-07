package files

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, jwt fiber.Handler) {
	group := app.Group("/files", jwt)

	group.Post("upload/*", UploadFile())
	group.Get("download/*", DownloadFile())
	group.Get("list/*", ListFiles())
	group.Get("search", SearchFiles())
}
