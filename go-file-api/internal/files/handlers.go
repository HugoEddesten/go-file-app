package files

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func UploadFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		vaultId := c.Locals("vaultId").(int)
		fileKey := c.Params("*")

		file, err := c.FormFile("file")
		if err != nil {
			return err
		}

		path := fmt.Sprintf("./uploads/%d/%s", vaultId, fileKey)
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

func CreateFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		vaultId := c.Locals("vaultId").(int)
		parentDirKey := c.Params("*")
		ext := c.Query("ext")

		baseDir := fmt.Sprintf("./uploads/%d/%s", vaultId, parentDirKey)

		if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		const baseName = "new"

		fullPath, name, err := nextAvailablePath(baseDir, baseName, ext)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		if ext == "" {
			if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError)
			}
			return c.SendStatus(fiber.StatusCreated)
		} else {
			file, err := os.Create(fullPath)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError)
			}
			defer file.Close()
		}

		clientFileKey := getClientKeyFromFilePath(fullPath)

		return c.Status(fiber.StatusCreated).JSON(FileResponse{
			Name: name,
			Key:  clientFileKey,
		})
	}
}

func DownloadFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		vaultId := c.Locals("vaultId").(int)
		action := c.Query("action", "send")

		fileKeyEncoded := c.Params("*")
		fileKey, err := url.QueryUnescape(fileKeyEncoded)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid filename encoding")
		}

		path := fmt.Sprintf("./uploads/%d/%s", vaultId, fileKey)

		file, err := os.Open(path)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "file not found")
		}
		defer file.Close()

		buffer := make([]byte, 512)
		_, _ = file.Read(buffer)
		contentType := http.DetectContentType(buffer)

		c.Set("Content-Type", contentType)

		if action == "download" {
			return c.Download(path)
		}
		return c.SendFile(path)
	}
}

func GetMetadata() fiber.Handler {
	return func(c *fiber.Ctx) error {
		vaultId := c.Locals("vaultId").(int)

		fileKeyEncoded := c.Params("*")
		fileKey, _ := url.QueryUnescape(fileKeyEncoded)
		path := fmt.Sprintf("./uploads/%d/%s", vaultId, fileKey)

		file, err := os.Open(path)
		if err != nil {
			return fiber.ErrNotFound
		}
		defer file.Close()

		stat, _ := file.Stat()

		buffer := make([]byte, 512)
		file.Read(buffer)
		mimeType := http.DetectContentType(buffer)

		meta := FileMetadata{
			Name:        stat.Name(),
			MimeType:    mimeType,
			Size:        stat.Size(),
			Editable:    strings.HasPrefix(mimeType, "text/"),
			Previewable: isPreviewable(mimeType),
		}

		return c.JSON(meta)
	}
}

func ListFiles() fiber.Handler {
	return func(c *fiber.Ctx) error {
		vaultId := c.Locals("vaultId").(int)
		fileKeyEncoded := c.Params("*")

		fileKey, err := url.QueryUnescape(fileKeyEncoded)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid filename encoding")
		}

		path := fmt.Sprintf("./uploads/%d/%s", vaultId, fileKey)

		files := make([]FileResponse, 0)

		entries, err := os.ReadDir(path)
		if err != nil {
			return c.JSON(files)
		}

		for _, entry := range entries {
			var key string
			if fileKey == "" {
				key = fmt.Sprintf("/%s", entry.Name())

			} else {
				key = fmt.Sprintf("/%s/%s", fileKey, entry.Name())
			}

			files = append(files, FileResponse{
				Name: entry.Name(),
				Key:  key,
			})
		}

		return c.JSON(files)
	}
}

func SearchFiles() fiber.Handler {
	return func(c *fiber.Ctx) error {
		vaultId := c.Locals("vaultId").(int)
		search := c.Query("q")
		if search == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing query parameter 'q'")
		}
		var matches []string

		root := fmt.Sprintf("./uploads/%d", vaultId)
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
