package storage

import (
	"context"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOService handles object storage operations
type MinIOService struct {
	client *minio.Client
}

// NewMinIOService creates a new MinIO service client
// endpoint: typically "localhost:9000" for local development
// accessKey: MINIO_ROOT_USER from docker-compose
// secretKey: MINIO_ROOT_PASSWORD from docker-compose
// useSSL: false for local development
func NewMinIOService(endpoint, accessKey, secretKey string, useSSL bool) (*MinIOService, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinIOService{
		client: client,
	}, nil
}

// EnsureBucket creates a bucket if it doesn't exist
// In MinIO, buckets are like top-level folders for your objects
func (s *MinIOService) EnsureBucket(ctx context.Context, bucketName string) error {
	exists, err := s.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exists {
		err = s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		log.Printf("Bucket '%s' created successfully", bucketName)
	}

	return nil
}

// UploadObject uploads a file to MinIO
// bucketName: the bucket to upload to (e.g., "file-vault")
// objectName: the path/name for the object (e.g., "vault-123/documents/file.pdf")
// reader: the file content
// objectSize: size in bytes (-1 if unknown, MinIO will buffer it)
// contentType: MIME type (e.g., "application/pdf", "image/png")
func (s *MinIOService) UploadObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, contentType string) error {
	_, err := s.client.PutObject(ctx, bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// DownloadObject retrieves a file from MinIO
// Returns a reader - don't forget to close it when done!
func (s *MinIOService) DownloadObject(ctx context.Context, bucketName, objectName string) (*minio.Object, error) {
	return s.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
}

// DeleteObject removes a file from MinIO
func (s *MinIOService) DeleteObject(ctx context.Context, bucketName, objectName string) error {
	return s.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

// ListObjects lists all objects with a given prefix
// prefix: filter objects by prefix (e.g., "vault-123/" to list all files in vault 123)
// Set recursive to true to list all objects in subdirectories
func (s *MinIOService) ListObjects(ctx context.Context, bucketName, prefix string, recursive bool) <-chan minio.ObjectInfo {
	return s.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: recursive,
	})
}

// TODO: Experiment with these features:
// - Presigned URLs (temporary download links): s.client.PresignedGetObject()
// - Object metadata: Add custom metadata in PutObjectOptions.UserMetadata
// - Versioning: Enable bucket versioning for file history
// - Server-side encryption: Add encryption in PutObjectOptions
// - Multipart uploads: For large files, use s.client.PutObjectWithContext with parts
// - Copy objects: s.client.CopyObject() for duplicating files
// - Stat object: s.client.StatObject() to get metadata without downloading
