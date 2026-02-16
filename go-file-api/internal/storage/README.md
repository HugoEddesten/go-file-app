# Storage Package

MinIO-based object storage service for file management.

## Setup

1. **Install the MinIO Go SDK:**
   ```bash
   go get github.com/minio/minio-go/v7
   ```

2. **Start MinIO via Docker Compose:**
   ```bash
   docker compose up -d minio
   ```

3. **Access MinIO Console:**
   - URL: http://localhost:9001
   - Username: `minioadmin`
   - Password: `minioadmin`

## Quick Start

```go
import "your-module/internal/storage"

// Initialize service
service, err := storage.NewMinIOService(
    "localhost:9000",
    "minioadmin",
    "minioadmin",
    false, // useSSL
)

// Create bucket
ctx := context.Background()
err = service.EnsureBucket(ctx, "file-vault")

// Upload a file
file, _ := os.Open("document.pdf")
defer file.Close()
info, _ := file.Stat()

err = service.UploadObject(
    ctx,
    "file-vault",
    "vault-1/documents/document.pdf",
    file,
    info.Size(),
    "application/pdf",
)

// Download a file
object, err := service.DownloadObject(ctx, "file-vault", "vault-1/documents/document.pdf")
defer object.Close()

// Copy to local file or stream to HTTP response
io.Copy(destination, object)
```

## Object Naming Convention

Recommended pattern: `vault-{vaultID}/{path}/{filename}`

Example: `vault-123/documents/reports/2024/annual.pdf`

## Environment Variables

Consider adding these to your `.env`:
```
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_USE_SSL=false
MINIO_BUCKET=file-vault
```

## Features to Explore

- **Presigned URLs**: Temporary download links without exposing credentials
- **Object Metadata**: Attach custom key-value pairs to files
- **Versioning**: Keep file history and restore previous versions
- **Server-side Encryption**: Encrypt files at rest
- **Multipart Uploads**: Efficient large file uploads with resume capability
- **Lifecycle Policies**: Auto-delete or transition old files

## Migration from Filesystem

To migrate existing files from `./uploads/`:
1. Read file from filesystem
2. Upload to MinIO with same vault structure
3. Update database metadata if needed
4. Verify upload success
5. Remove old file (or keep as backup initially)
