package files

import (
	"fmt"
	"go-file-api/internal/storage"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func UploadFile(minIOService *storage.MinIOService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		vaultId := c.Locals("vaultId").(int)
		fileKey := c.Locals("requestedVaultPath")

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
		parentDirKey := c.Locals("requestedVaultPath")
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
		fileKey := c.Locals("requestedVaultPath")
		action := c.Query("action", "send")

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

		fileKey := c.Locals("requestedVaultPath")
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
		fileKey := c.Locals("requestedVaultPath")

		path := fmt.Sprintf("./uploads/%d%s", vaultId, fileKey)

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

func RenameFile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		vaultId := c.Locals("vaultId").(int)

		// Get the original file/folder path from middleware
		// The VaultAccessMiddleware already resolved and cleaned the path
		fileKey := c.Locals("requestedVaultPath").(string)

		// Parse request body for new name
		var req RenameRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
		}

		// Validate new name
		if req.NewName == "" {
			return fiber.NewError(fiber.StatusBadRequest, "newName is required")
		}

		// Prevent path traversal attacks
		if strings.Contains(req.NewName, "/") || strings.Contains(req.NewName, "\\") {
			return fiber.NewError(fiber.StatusBadRequest, "newName cannot contain path separators")
		}

		// Construct old and new paths
		oldPath := fmt.Sprintf("./uploads/%d%s", vaultId, fileKey)

		// Get parent directory
		parentDir := filepath.Dir(oldPath)
		newPath := filepath.Join(parentDir, req.NewName)

		// Check if old path exists
		stat, err := os.Stat(oldPath)
		if os.IsNotExist(err) {
			return fiber.NewError(fiber.StatusNotFound, "file or folder not found")
		}

		// Check if new path already exists
		if _, err := os.Stat(newPath); err == nil {
			return fiber.NewError(fiber.StatusConflict, "a file or folder with that name already exists")
		}

		// Try to rename with retry logic for Windows file locking issues
		if err := rename(oldPath, newPath, stat.IsDir()); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to rename file or folder")
		}

		// Generate new client key
		clientKey := getClientKeyFromFilePath(newPath)

		return c.JSON(FileResponse{
			Name: req.NewName,
			Key:  clientKey,
		})
	}
}

// renameWithFallback attempts to rename a file or directory using os.Rename with retries.
// If that fails (common on Windows due to file locking), it falls back to copy-then-delete.
func rename(oldPath, newPath string, isDir bool) error {
	if isDir {
		// For directories, use recursive copy
		if err := copyDir(oldPath, newPath); err != nil {
			return err
		}
	} else {
		// For files, use file copy
		if err := copyFile(oldPath, newPath); err != nil {
			fmt.Print(err.Error())
			return err
		}
	}

	// If copy succeeded, remove the old path
	// Use RemoveAll to handle both files and directories
	return os.RemoveAll(oldPath)
}

// copyFile copies a single file from src to dst, preserving permissions
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Get source file info for permissions
	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	// Create destination file with same permissions
	destFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, sourceInfo.Mode())
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy file contents
	_, err = sourceFile.WriteTo(destFile)
	return err
}

// copyDir recursively copies a directory from src to dst
func copyDir(src, dst string) error {
	// Get source directory info
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create destination directory with same permissions
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// Read directory contents
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectory
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}
