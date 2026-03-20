package files

import (
	"fmt"
	"go-file-api/internal/locals"
	"go-file-api/internal/storage"
	"go-file-api/internal/vault"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func UploadFile(minIOService *storage.MinIOService, vaultRepo *vault.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		vaultId := locals.VaultId(c)
		fileKey := locals.RequestedVaultPath(c)

		file, err := c.FormFile("file")
		if err != nil {
			return err
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		prefix := fmt.Sprintf("vault-%d%s", vaultId, fileKey)
		savePath := path.Join(prefix, file.Filename)
		contentType := getContentType(fileReader)

		var oldSize int64
		if existing, err := minIOService.StatObject(c.Context(), storage.VaultBucket, savePath); err == nil {
			oldSize = existing.Size
		}

		err = minIOService.UploadObject(c.Context(), storage.VaultBucket, savePath, fileReader, -1, contentType)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to upload file")
		}

		_ = vaultRepo.UpdateStorageUsed(c.Context(), vaultId, file.Size-oldSize)

		return c.JSON(fiber.Map{
			"status":   "uploaded",
			"filename": file.Filename,
			"size":     file.Size,
		})
	}
}

func CreateFile(minIOService *storage.MinIOService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		vaultId := locals.VaultId(c)
		parentDirKey := locals.RequestedVaultPath(c)
		ext := c.Query("ext")

		bucketDir := getBucketPath(vaultId, parentDirKey)
		const baseName = "new"

		bucketKey, name, err := minIOService.NextAvailablePath(c.Context(), storage.VaultBucket, bucketDir, baseName, ext)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		if ext == "" {
			// Represent the folder with a trailing-slash placeholder object
			folderKey := bucketKey + "/"
			if err := minIOService.UploadObject(c.Context(), storage.VaultBucket, folderKey, strings.NewReader(""), 0, "application/x-directory"); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError)
			}
			return c.SendStatus(fiber.StatusCreated)
		}

		if err := minIOService.UploadObject(c.Context(), storage.VaultBucket, bucketKey, strings.NewReader(""), 0, "application/octet-stream"); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusCreated).JSON(FileResponse{
			Name: name,
			Key:  getClientKeyFromBucketPath(bucketKey),
		})
	}
}

func DownloadFile(minIOService *storage.MinIOService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		vaultId := locals.VaultId(c)
		fileKey := locals.RequestedVaultPath(c)
		action := c.Query("action", "send")

		objectPath := fmt.Sprintf("vault-%d%s", vaultId, fileKey)

		stat, err := minIOService.StatObject(c.Context(), storage.VaultBucket, objectPath)
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("File not found")
		}

		object, err := minIOService.DownloadObject(c.Context(), storage.VaultBucket, objectPath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Unable to download file")
		}

		filename := path.Base(fileKey)

		c.Set(fiber.HeaderContentType, stat.ContentType)

		// Set Content-Disposition based on action
		if action == "download" {
			c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"", filename))
		} else {
			c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("inline; filename=\"%s\"", filename))
		}

		return c.SendStream(object)
	}
}

func GetMetadata(minIOService *storage.MinIOService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		vaultId := locals.VaultId(c)

		clientKey := locals.RequestedVaultPath(c)
		fileKey := getBucketPath(vaultId, clientKey)

		objectInfo, err := minIOService.StatObject(c.Context(), storage.VaultBucket, fileKey)
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("File not found")
		}

		meta := FileMetadata{
			Name:        path.Base(objectInfo.Key),
			MimeType:    objectInfo.ContentType,
			Size:        objectInfo.Size,
			Editable:    strings.HasPrefix(objectInfo.ContentType, "text/"),
			Previewable: isPreviewable(objectInfo.ContentType),
		}

		return c.JSON(meta)
	}
}

func ListFiles(minIOService *storage.MinIOService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		vaultId := locals.VaultId(c)
		requestedKey := locals.RequestedVaultPath(c)

		fileKey := getBucketPath(vaultId, requestedKey)
		if !strings.HasSuffix(fileKey, "/") {
			fileKey += "/"
		}
		files := make([]FileResponse, 0)
		entries := minIOService.ListObjects(c.Context(), storage.VaultBucket, fileKey, false)

		for entry := range entries {
			if entry.Err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": entry.Err.Error()})
			}

			// Skip the folder placeholder for the listed directory itself
			if entry.Key == fileKey {
				continue
			}

			clientKey := getClientKeyFromFilePath(entry.Key)

			files = append(files, FileResponse{
				Name: path.Base(entry.Key),
				Key:  clientKey,
			})
		}

		return c.JSON(files)
	}
}

func SearchFiles(minIOService *storage.MinIOService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		vaultId := locals.VaultId(c)
		search := c.Query("q")
		if search == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing query parameter 'q'")
		}

		prefix := getBucketPath(vaultId, "/")
		entries := minIOService.ListObjects(c.Context(), storage.VaultBucket, prefix, true)

		var matches []string
		for entry := range entries {
			if entry.Err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": entry.Err.Error()})
			}
			filename := path.Base(entry.Key)
			if strings.Contains(strings.ToLower(filename), strings.ToLower(search)) {
				matches = append(matches, filename)
			}
		}

		return c.JSON(matches)
	}
}

func DeleteFile(minIOService *storage.MinIOService, vaultRepo *vault.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		vaultId := locals.VaultId(c)
		fileKey := locals.RequestedVaultPath(c)

		bucketPath := getBucketPath(vaultId, fileKey)

		// Try as a single file first
		stat, err := minIOService.StatObject(c.Context(), storage.VaultBucket, bucketPath)
		if err == nil {
			if err := minIOService.DeleteObject(c.Context(), storage.VaultBucket, bucketPath); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "failed to delete file")
			}
			_ = vaultRepo.UpdateStorageUsed(c.Context(), vaultId, -stat.Size)
			return c.SendStatus(fiber.StatusNoContent)
		}

		// Treat as a folder — collect all objects with their sizes, then delete
		type objEntry struct {
			key  string
			size int64
		}
		prefix := bucketPath + "/"
		var objects []objEntry
		for obj := range minIOService.ListObjects(c.Context(), storage.VaultBucket, prefix, true) {
			if obj.Err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "failed to list objects")
			}
			objects = append(objects, objEntry{obj.Key, obj.Size})
		}

		if len(objects) == 0 {
			return fiber.NewError(fiber.StatusNotFound, "file or folder not found")
		}

		var freed int64
		for _, obj := range objects {
			if err := minIOService.DeleteObject(c.Context(), storage.VaultBucket, obj.key); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "failed to delete object")
			}
			freed += obj.size
		}
		_ = vaultRepo.UpdateStorageUsed(c.Context(), vaultId, -freed)

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func MoveFile(minIOService *storage.MinIOService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		vaultId := locals.VaultId(c)
		fileKey := locals.RequestedVaultPath(c)

		var req MoveRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
		}

		if req.DestinationKey == "" {
			return fiber.NewError(fiber.StatusBadRequest, "destinationKey is required")
		}

		oldBucketPath := getBucketPath(vaultId, fileKey)
		filename := path.Base(fileKey)
		newClientKey := path.Join(req.DestinationKey, filename)
		newBucketPath := getBucketPath(vaultId, newClientKey)

		_, err := minIOService.StatObject(c.Context(), storage.VaultBucket, oldBucketPath)
		if err == nil {
			if exists, _ := minIOService.ObjectExists(c.Context(), storage.VaultBucket, newBucketPath); exists {
				return fiber.NewError(fiber.StatusConflict, "a file with that name already exists at the destination")
			}
			if err := minIOService.CopyObject(c.Context(), storage.VaultBucket, oldBucketPath, newBucketPath); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "failed to move file")
			}
			if err := minIOService.DeleteObject(c.Context(), storage.VaultBucket, oldBucketPath); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "failed to clean up original file")
			}
		} else {
			oldPrefix := oldBucketPath + "/"
			newPrefix := newBucketPath + "/"

			var keys []string
			for obj := range minIOService.ListObjects(c.Context(), storage.VaultBucket, oldPrefix, true) {
				if obj.Err != nil {
					return fiber.NewError(fiber.StatusInternalServerError, "failed to list objects")
				}
				keys = append(keys, obj.Key)
			}

			if len(keys) == 0 {
				return fiber.NewError(fiber.StatusNotFound, "file or folder not found")
			}

			for obj := range minIOService.ListObjects(c.Context(), storage.VaultBucket, newPrefix, false) {
				if obj.Err == nil {
					return fiber.NewError(fiber.StatusConflict, "a folder with that name already exists at the destination")
				} else {
					break
				}
			}

			for _, key := range keys {
				suffix := strings.TrimPrefix(key, oldPrefix)
				newKey := newPrefix + suffix
				if err := minIOService.CopyObject(c.Context(), storage.VaultBucket, key, newKey); err != nil {
					return fiber.NewError(fiber.StatusInternalServerError, "failed to move folder")
				}
				if err := minIOService.DeleteObject(c.Context(), storage.VaultBucket, key); err != nil {
					return fiber.NewError(fiber.StatusInternalServerError, "failed to clean up original objects")
				}
			}
		}

		return c.JSON(FileResponse{
			Name: filename,
			Key:  newClientKey,
		})
	}
}

func RenameFile(minIOService *storage.MinIOService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		vaultId := locals.VaultId(c)
		fileKey := locals.RequestedVaultPath(c)

		var req RenameRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
		}

		if req.NewName == "" {
			return fiber.NewError(fiber.StatusBadRequest, "newName is required")
		}

		if strings.Contains(req.NewName, "/") || strings.Contains(req.NewName, "\\") {
			return fiber.NewError(fiber.StatusBadRequest, "newName cannot contain path separators")
		}

		oldBucketPath := getBucketPath(vaultId, fileKey)
		newClientKey := path.Join(path.Dir(fileKey), req.NewName)
		newBucketPath := getBucketPath(vaultId, newClientKey)

		// Check if source is a single file
		_, err := minIOService.StatObject(c.Context(), storage.VaultBucket, oldBucketPath)
		if err == nil {
			// It's a file — check for conflict then copy+delete
			if exists, _ := minIOService.ObjectExists(c.Context(), storage.VaultBucket, newBucketPath); exists {
				return fiber.NewError(fiber.StatusConflict, "a file or folder with that name already exists")
			}
			if err := minIOService.CopyObject(c.Context(), storage.VaultBucket, oldBucketPath, newBucketPath); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "failed to rename file")
			}
			if err := minIOService.DeleteObject(c.Context(), storage.VaultBucket, oldBucketPath); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "failed to clean up original file")
			}
		} else {
			// Treat as a folder — collect all objects under the old prefix
			oldPrefix := oldBucketPath + "/"
			newPrefix := newBucketPath + "/"

			var keys []string
			for obj := range minIOService.ListObjects(c.Context(), storage.VaultBucket, oldPrefix, true) {
				if obj.Err != nil {
					return fiber.NewError(fiber.StatusInternalServerError, "failed to list objects")
				}
				keys = append(keys, obj.Key)
			}

			if len(keys) == 0 {
				return fiber.NewError(fiber.StatusNotFound, "file or folder not found")
			}

			// Check for conflict at the destination prefix
			for obj := range minIOService.ListObjects(c.Context(), storage.VaultBucket, newPrefix, false) {
				if obj.Err == nil {
					return fiber.NewError(fiber.StatusConflict, "a file or folder with that name already exists")
				} else {
					break
				}
			}

			// Copy each object to the new prefix, then delete the original
			for _, key := range keys {
				suffix := strings.TrimPrefix(key, oldPrefix)
				newKey := newPrefix + suffix
				if err := minIOService.CopyObject(c.Context(), storage.VaultBucket, key, newKey); err != nil {
					return fiber.NewError(fiber.StatusInternalServerError, "failed to rename folder")
				}
				if err := minIOService.DeleteObject(c.Context(), storage.VaultBucket, key); err != nil {
					return fiber.NewError(fiber.StatusInternalServerError, "failed to clean up original objects")
				}
			}
		}

		return c.JSON(FileResponse{
			Name: req.NewName,
			Key:  newClientKey,
		})
	}
}
