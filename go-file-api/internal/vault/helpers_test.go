package vault

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestResolveVaultPath_Wildcard(t *testing.T) {
	app := fiber.New()
	var resolvedPath string
	var shouldValidate bool

	app.Get("/vault/:vaultId/*", func(c *fiber.Ctx) error {
		resolvedPath, shouldValidate = ResolveVaultPath(c)
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/vault/1/documents/file.txt", nil)
	app.Test(req)

	assert.Equal(t, "/documents/file.txt", resolvedPath)
	assert.True(t, shouldValidate)
}

func TestResolveVaultPath_QueryParam(t *testing.T) {
	app := fiber.New()
	var resolvedPath string
	var shouldValidate bool

	// Route with NO wildcard
	app.Get("/vault/:vaultId/list", func(c *fiber.Ctx) error {
		resolvedPath, shouldValidate = ResolveVaultPath(c)
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/vault/1/list?path=/uploads/image.png", nil)
	app.Test(req)

	// Should resolve from query param
	assert.Equal(t, "/uploads/image.png", resolvedPath)
	assert.True(t, shouldValidate)
}

func TestResolveVaultPath_JSONBody(t *testing.T) {
	app := fiber.New()
	var resolvedPath string
	var shouldValidate bool

	app.Post("/vault/:vaultId/create", func(c *fiber.Ctx) error {
		resolvedPath, shouldValidate = ResolveVaultPath(c)
		return c.SendString("ok")
	})

	body := map[string]string{"path": "/uploads/newfile.txt"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/vault/1/create", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	app.Test(req)

	// Should resolve from JSON body
	assert.Equal(t, "/uploads/newfile.txt", resolvedPath)
	assert.True(t, shouldValidate)
}

func TestResolveVaultPath_NoPath(t *testing.T) {
	app := fiber.New()
	var resolvedPath string
	var shouldValidate bool

	app.Get("/vault/:vaultId/info", func(c *fiber.Ctx) error {
		resolvedPath, shouldValidate = ResolveVaultPath(c)
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/vault/1/info", nil)
	app.Test(req)

	assert.Equal(t, "/", resolvedPath)
	// No path provided anywhere, so shouldValidate is false
	assert.False(t, shouldValidate)
}

func TestPathAllowed_RootAccess(t *testing.T) {
	assert.True(t, pathAllowed("/", "/any/path/here"))
	assert.True(t, pathAllowed("/", "/"))
}

func TestPathAllowed_SubdirectoryAccess(t *testing.T) {
	assert.True(t, pathAllowed("/documents", "/documents"))
	assert.True(t, pathAllowed("/documents", "/documents/file.txt"))
	assert.True(t, pathAllowed("/documents", "/documents/sub/file.txt"))
}

func TestPathAllowed_DeniedPaths(t *testing.T) {
	assert.False(t, pathAllowed("/documents", "/admin"))
	assert.False(t, pathAllowed("/documents", "/doc")) // Not a subdirectory
	assert.False(t, pathAllowed("/doc", "/documents")) // Similar prefix but not allowed
}

func TestCleanPath(t *testing.T) {
	assert.Equal(t, "/documents/file.txt", cleanPath("/documents/file.txt"))
	assert.Equal(t, "/documents/file.txt", cleanPath("documents/file.txt"))
	// BUG: path.Clean resolves ".." so "/documents/../admin" becomes "/admin"
	// The Contains check happens AFTER Clean, so it doesn't catch this traversal
	assert.Equal(t, "/admin", cleanPath("/documents/../admin"))
	assert.Equal(t, "/", cleanPath("/"))
	// path.Clean("/" + "..") becomes "/", so Contains doesn't find ".."
	assert.Equal(t, "/", cleanPath(".."))
	assert.Equal(t, "/admin", cleanPath("../admin"))
}
