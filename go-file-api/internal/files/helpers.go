package files

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// getClientKeyFromBucketPath converts a MinIO bucket path like "vault-1/documents/file.txt"
// to a client key like "/documents/file.txt".
func getClientKeyFromBucketPath(bucketPath string) string {
	idx := strings.Index(bucketPath, "/")
	if idx == -1 {
		return "/"
	}
	return bucketPath[idx:]
}

func getClientKeyFromFilePath(fullPath string) string {
	clean := filepath.Clean(fullPath)

	parts := strings.Split(clean, string(os.PathSeparator))

	clientKey := strings.Join(parts[1:], "/")

	return "/" + clientKey
}

func isPreviewable(mime string) bool {
	switch {
	case strings.HasPrefix(mime, "image/"):
		return true
	case mime == "application/pdf":
		return true
	case strings.HasPrefix(mime, "text/"):
		return true
	default:
		return false
	}
}

func getContentType(reader io.Reader) string {
	buffer := make([]byte, 512)
	_, _ = reader.Read(buffer)
	contentType := http.DetectContentType(buffer)

	// Reset reader position if it supports seeking
	if seeker, ok := reader.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}

	return contentType
}

func getBucketPath(vaultId int, clientKey string) string {
	return fmt.Sprintf("vault-%d%s", vaultId, clientKey)
}
