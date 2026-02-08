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
	assert.False(t, shouldValidate)
}

func TestResolveVaultId_FromURLParam(t *testing.T) {
	app := fiber.New()
	var resolvedId int
	var resolvedErr error

	app.Get("/vault/:vaultId/info", func(c *fiber.Ctx) error {
		resolvedId, resolvedErr = ResolveVaultId(c)
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/vault/42/info", nil)
	app.Test(req)

	assert.NoError(t, resolvedErr)
	assert.Equal(t, 42, resolvedId)
}

func TestResolveVaultId_FromQueryParam(t *testing.T) {
	app := fiber.New()
	var resolvedId int
	var resolvedErr error

	// Route without :vaultId param - should fall through to query
	app.Get("/vault/info", func(c *fiber.Ctx) error {
		resolvedId, resolvedErr = ResolveVaultId(c)
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/vault/info?vaultId=99", nil)
	app.Test(req)

	assert.NoError(t, resolvedErr)
	assert.Equal(t, 99, resolvedId)
}

func TestResolveVaultId_FromLocals(t *testing.T) {
	app := fiber.New()
	var resolvedId int
	var resolvedErr error

	// Route without :vaultId param or query - should fall through to locals
	app.Get("/vault/info", func(c *fiber.Ctx) error {
		// Simulate middleware setting vaultId in locals
		c.Locals("vaultId", 123)
		resolvedId, resolvedErr = ResolveVaultId(c)
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/vault/info", nil)
	app.Test(req)

	assert.NoError(t, resolvedErr)
	assert.Equal(t, 123, resolvedId)
}

func TestResolveVaultId_NotFound(t *testing.T) {
	app := fiber.New()
	var resolvedId int
	var resolvedErr error

	// No vaultId anywhere
	app.Get("/vault/info", func(c *fiber.Ctx) error {
		resolvedId, resolvedErr = ResolveVaultId(c)
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/vault/info", nil)
	app.Test(req)

	assert.Error(t, resolvedErr)
	assert.Equal(t, "vaultId not found", resolvedErr.Error())
	assert.Equal(t, 0, resolvedId)
}

func TestResolveVaultId_InvalidURLParam(t *testing.T) {
	app := fiber.New()
	var resolvedId int
	var resolvedErr error

	app.Get("/vault/:vaultId/info", func(c *fiber.Ctx) error {
		resolvedId, resolvedErr = ResolveVaultId(c)
		return c.SendString("ok")
	})

	// Non-numeric vaultId - should return error
	req := httptest.NewRequest("GET", "/vault/invalid/info", nil)
	app.Test(req)

	assert.Error(t, resolvedErr)
	assert.Equal(t, 0, resolvedId)
}

func TestResolveVaultId_ZeroValue(t *testing.T) {
	app := fiber.New()
	var resolvedId int
	var resolvedErr error

	app.Get("/vault/info", func(c *fiber.Ctx) error {
		resolvedId, resolvedErr = ResolveVaultId(c)
		return c.SendString("ok")
	})

	// vaultId=0 should be rejected (not > 0)
	req := httptest.NewRequest("GET", "/vault/info?vaultId=0", nil)
	app.Test(req)

	assert.Error(t, resolvedErr)
	assert.Equal(t, 0, resolvedId)
}

func TestResolveVaultId_NegativeValue(t *testing.T) {
	app := fiber.New()
	var resolvedId int
	var resolvedErr error

	app.Get("/vault/info", func(c *fiber.Ctx) error {
		resolvedId, resolvedErr = ResolveVaultId(c)
		return c.SendString("ok")
	})

	// vaultId=-5 should be rejected (not > 0)
	req := httptest.NewRequest("GET", "/vault/info?vaultId=-5", nil)
	app.Test(req)

	assert.Error(t, resolvedErr)
	assert.Equal(t, 0, resolvedId)
}

func TestResolveVaultId_PriorityOrder(t *testing.T) {
	app := fiber.New()
	var resolvedId int
	var resolvedErr error

	// All three sources present - should prioritize URL param
	app.Get("/vault/:vaultId/info", func(c *fiber.Ctx) error {
		c.Locals("vaultId", 999)
		resolvedId, resolvedErr = ResolveVaultId(c)
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/vault/10/info?vaultId=20", nil)
	app.Test(req)

	assert.NoError(t, resolvedErr)
	assert.Equal(t, 10, resolvedId) // URL param takes priority
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
	// Normal paths
	assert.Equal(t, "/documents/file.txt", cleanPath("/documents/file.txt"))
	assert.Equal(t, "/documents/file.txt", cleanPath("documents/file.txt"))
	assert.Equal(t, "/", cleanPath("/"))

	// Path traversal attempts - path.Clean normalizes these
	// Note: Path traversal is prevented by pathAllowed() checking permissions,
	// not by cleanPath itself. cleanPath just normalizes the path.
	assert.Equal(t, "/admin", cleanPath("/documents/../admin"))
	assert.Equal(t, "/", cleanPath(".."))
	assert.Equal(t, "/admin", cleanPath("../admin"))
}
