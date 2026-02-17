package files

import (
	"fmt"
	"go-file-api/internal/storage"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func UploadFile(minIOService *storage.MinIOService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		vaultId := c.Locals("vaultId").(int)
		fileKey := c.Params("*")

		file, err := c.FormFile("file")
		if err != nil {
			return err
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		path := fmt.Sprintf("vault-%d/%s", vaultId, fileKey)
		savePath := filepath.Join(path, file.Filename)
		contentType := getContentType(fileReader)

		err = minIOService.UploadObject(c.Context(), "file-vault", savePath, fileReader, -1, contentType)
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

		contentType := getContentType(file)

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

func ListFiles(minIOService *storage.MinIOService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		vaultId := c.Locals("vaultId").(int)
		fileKeyEncoded := c.Params("*")

		fileKey, err := url.QueryUnescape(fileKeyEncoded)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid filename encoding")
		}

		filePath := fmt.Sprintf("vault-%d\\%s", vaultId, fileKey)
		fmt.Print(filePath)
		files := make([]FileResponse, 0)

		entries := minIOService.ListObjects(c.Context(), "file-vault", filePath, false)

		for entry := range entries {
			if entry.Err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": entry.Err.Error()})
			}
			var key string
			key = fmt.Sprintf("/%s", entry.Key)

			files = append(files, FileResponse{
				Name: path.Base(entry.Key),
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
