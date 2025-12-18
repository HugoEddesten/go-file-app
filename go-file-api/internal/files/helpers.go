package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func nextAvailablePath(dir, base, ext string) (string, string, error) {
	name := base + ext
	fullPath := filepath.Join(dir, name)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fullPath, name, nil
	}

	for i := 2; ; i++ {
		name = fmt.Sprintf("%s(%d)%s", base, i, ext)
		fullPath = filepath.Join(dir, name)

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fullPath, name, nil
		}
	}
}

func getClientKeyFromFilePath(fullPath string) string {
	clean := filepath.Clean(fullPath)

	parts := strings.Split(clean, string(os.PathSeparator))

	for i := 0; i < len(parts)-2; i++ {
		if parts[i] == "uploads" {
			return "\\" + filepath.Join(parts[i+2:]...)
		}
	}

	return "\\" + fullPath
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
