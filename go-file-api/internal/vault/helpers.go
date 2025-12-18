package vault

import (
	"errors"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ResolveVaultPath(c *fiber.Ctx) (string, bool) {
	// 1. Wildcard route
	if p := c.Params("*"); p != "" {
		return cleanPath(p), true
	}

	// 2. URL param
	if p := c.Params("path"); p != "" {
		return cleanPath(p), true
	}

	// 3. Query
	if p := c.Query("path"); p != "" {
		return cleanPath(p), true
	}

	// 4. JSON body
	if v := c.Locals("path"); v != nil {
		return cleanPath(v.(string)), true
	}

	return "", false
}

func ResolveVaultId(c *fiber.Ctx) (int, error) {
	// 2. URL param
	if v, err := c.ParamsInt("vaultId"); err == nil && v != 0 {
		return v, nil
	}

	// 3. Query
	if v := c.QueryInt("vaultId"); v != 0 {
		return v, nil
	}

	// 4. JSON body
	if v := c.Locals("vaultId").(int); v != 0 {
		return v, nil
	}

	return 0, errors.New("NO VAULT_ID FOUND")
}

func cleanPath(p string) string {
	p = path.Clean("/" + p)
	if strings.Contains(p, "..") {
		return ""
	}
	return p
}

func pathAllowed(allowed, requested string) bool {
	if allowed == "/" {
		return true
	}

	if !strings.HasPrefix(requested, allowed) {
		return false
	}

	if len(requested) == len(allowed) {
		return true
	}

	return requested[len(allowed)] == '/'
}
