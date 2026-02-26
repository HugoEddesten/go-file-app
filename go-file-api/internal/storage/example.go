package storage

import (
	"context"
	"log"
)

// Example usage - you can delete this file once you're comfortable with the API

// InitializeStorage sets up MinIO and creates the default bucket
func InitializeStorage() (*MinIOService, error) {
	// Create MinIO service
	// These should come from environment variables in production
	service, err := NewMinIOService(
		"localhost:9000",
		"minioadmin",
		"minioadmin",
		false, // useSSL - false for local development
	)
	if err != nil {
		return nil, err
	}

	// Create default bucket for file storage
	ctx := context.Background()
	err = service.EnsureBucket(ctx, "file-vault")
	if err != nil {
		return nil, err
	}

	log.Println("MinIO storage initialized successfully")
	return service, nil
}

// Example: How you might structure object names
// func GetObjectPath(vaultID int, filepath string) string {
// 	return fmt.Sprintf("vault-%d/%s", vaultID, filepath)
// }

// Example: Upload a file
// func (s *MinIOService) UploadVaultFile(ctx context.Context, vaultID int, filename string, reader io.Reader, size int64, contentType string) error {
// 	objectName := GetObjectPath(vaultID, filename)
// 	return s.UploadObject(ctx, "file-vault", objectName, reader, size, contentType)
// }

// Example: Download a file
// func (s *MinIOService) DownloadVaultFile(ctx context.Context, vaultID int, filename string) (*minio.Object, error) {
// 	objectName := GetObjectPath(vaultID, filename)
// 	return s.DownloadObject(ctx, "file-vault", objectName)
// }
