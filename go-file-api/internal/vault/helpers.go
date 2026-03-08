package vault

import (
	"errors"
	"net/url"
	"path"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ResolveVaultPath(c *fiber.Ctx) (string, bool, error) {
	extractors := []func() (string, bool){
		func() (string, bool) {
			if c.Route().Path != "" && strings.Contains(c.Route().Path, "*") {
				return c.Params("*"), true
			}
			return "", false
		},
		func() (string, bool) {
			p := c.Params("path")
			return p, p != ""
		},
		func() (string, bool) {
			p := c.Query("path")
			return p, p != ""
		},
		func() (string, bool) {
			body := new(PathBodyValidation)
			if err := c.BodyParser(body); err == nil && body.Path != "" {
				return body.Path, true
			}
			return "", false
		},
	}

	for _, extract := range extractors {
		p, found := extract()
		if found {
			path, err := cleanPath(p)
			if err != nil {
				return "", false, err
			}
			return path, true, nil
		}
	}

	return "", false, nil
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

func cleanPath(p string) (string, error) {
	// path.Clean normalizes the path and resolves ".." sequences
	// This prevents path traversal by converting paths like "/documents/../admin"
	// into "/admin", which will then be checked against user permissions
	fileKey, err := url.QueryUnescape(p)
	if err != nil {
		return "", errors.New("invalid filename encoding")
	}

	return path.Clean("/" + fileKey), nil
}

// editableByAdmin returns the subset of targets that adminEntries have permission to edit.
// A target is editable if it is not an owner and at least one admin entry covers its path with
// Admin role or higher (Owner or Admin).
func editableByAdmin(adminEntries []VaultUser, targets []VaultUser) []VaultUser {
	var editable []VaultUser
	for _, target := range targets {
		if target.Role == VaultRoleOwner {
			continue
		}
		for _, admin := range adminEntries {
			if admin.Role <= VaultRoleAdmin && pathAllowed(admin.Path, target.Path) {
				editable = append(editable, target)
				break
			}
		}
	}
	return editable
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
