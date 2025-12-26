package vault

import (
	"errors"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ResolveVaultPath(c *fiber.Ctx) (string, bool) {
	// 1. Wildcard route
	if p := c.Params("*", "/"); p != "" {
		return cleanPath(p), true
	}

	// 2. URL param
	if p := c.Params("path", "/"); p != "" {
		return cleanPath(p), true
	}

	// 3. Query
	if p := c.Query("path", "/"); p != "" {
		return cleanPath(p), true
	}

	// 4. JSON body
	body := new(PathBodyValidation)
	if err := c.BodyParser(body); err == nil {
		return cleanPath(body.Path), true
	}

	return "/", false
}

func ResolveVaultId(c *fiber.Ctx) (int, error) {
	// 1. URL param (/vaults/:vaultId)
	if v, err := c.ParamsInt("vaultId"); err == nil && v > 0 {
		return v, nil
	}

	// 2. Query (?vaultId=)
	if v := c.QueryInt("vaultId"); v > 0 {
		return v, nil
	}

	// 3. Locals (set by middleware or handler)
	if v, ok := c.Locals("vaultId").(int); ok && v > 0 {
		return v, nil
	}

	return 0, errors.New("vaultId not found")
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
