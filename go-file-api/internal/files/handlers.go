package files

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func UploadFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId").(int)
		fileKey := c.Params("*")

		file, err := c.FormFile("file")
		if err != nil {
			return err
		}

		path := fmt.Sprintf("./uploads/%d/%s", userId, fileKey)
		os.MkdirAll(path, os.ModePerm)

		savePath := filepath.Join(path, file.Filename)
		err = c.SaveFile(file, savePath)
		if err != nil {
			return err
		}

		return c.JSON(fiber.Map{
			"status":   "uploaded",
			"filename": file.Filename,
			"size":     file.Size,
		})
	}
}

func DownloadFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId").(int)
		action := c.Query("action", "send")

		fileKeyEncoded := c.Params("*")
		fileKey, err := url.QueryUnescape(fileKeyEncoded)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid filename encoding")
		}

		path := fmt.Sprintf("./uploads/%d/%s", userId, fileKey)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fiber.NewError(fiber.StatusNotFound, "File not found")
		}

		if action == "download" {
			return c.Download(path)
		}
		return c.SendFile(path)
	}
}

// func GetMetadata() fiber.Handler {
// 	return func(c *fiber.Ctx) error {

// 	}
// }

func ListFiles() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId").(int)

		fileKeyEncoded := c.Params("*")
		fileKey, err := url.QueryUnescape(fileKeyEncoded)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid filename encoding")
		}

		path := fmt.Sprintf("./uploads/%d/%s", userId, fileKey)

		files := make([]FileResponse, 0)

		entries, err := os.ReadDir(path)
		if err != nil {
			return c.JSON(files)
		}

		for _, entry := range entries {
			files = append(files, FileResponse{
				Name: entry.Name(),
				Key:  fmt.Sprintf("%s/%s", fileKey, entry.Name()),
			})
		}

		return c.JSON(files)
	}
}

func SearchFiles() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId").(int)
		search := c.Query("q")
		if search == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing query parameter 'q'")
		}
		var matches []string

		root := fmt.Sprintf("./uploads/%d", userId)
		filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				filename := filepath.Base(path)
				if strings.Contains(strings.ToLower(filename), search) {
					matches = append(matches, filename)
				}
			}

			return nil
		})

		return c.JSON(matches)
	}
}
