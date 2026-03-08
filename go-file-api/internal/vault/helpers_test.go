package vault

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestEditableByAdmin(t *testing.T) {
	owner := VaultUser{Id: 1, Path: "/", Role: VaultRoleOwner}
	adminRoot := VaultUser{Id: 2, Path: "/", Role: VaultRoleAdmin}
	adminDocs := VaultUser{Id: 3, Path: "/documents", Role: VaultRoleAdmin}
	editorRoot := VaultUser{Id: 4, Path: "/", Role: VaultRoleEditor}

	targetRoot := VaultUser{Id: 10, Path: "/", Role: VaultRoleEditor}
	targetDocs := VaultUser{Id: 11, Path: "/documents/file.txt", Role: VaultRoleViewer}
	targetOther := VaultUser{Id: 12, Path: "/other", Role: VaultRoleViewer}

	t.Run("owner is never editable", func(t *testing.T) {
		result := editableByAdmin([]VaultUser{adminRoot}, []VaultUser{owner})
		assert.Empty(t, result)
	})

	t.Run("admin with root access can edit any target", func(t *testing.T) {
		result := editableByAdmin([]VaultUser{adminRoot}, []VaultUser{targetRoot, targetDocs, targetOther})
		assert.Len(t, result, 3)
	})

	t.Run("admin scoped to /documents can edit targets under /documents", func(t *testing.T) {
		result := editableByAdmin([]VaultUser{adminDocs}, []VaultUser{targetDocs, targetOther})
		assert.Len(t, result, 1)
		assert.Equal(t, targetDocs.Id, result[0].Id)
	})

	t.Run("editor cannot edit anyone even with path coverage", func(t *testing.T) {
		result := editableByAdmin([]VaultUser{editorRoot}, []VaultUser{targetDocs})
		assert.Empty(t, result)
	})

	t.Run("empty admin entries returns nothing", func(t *testing.T) {
		result := editableByAdmin([]VaultUser{}, []VaultUser{targetDocs})
		assert.Empty(t, result)
	})

	t.Run("empty targets returns nothing", func(t *testing.T) {
		result := editableByAdmin([]VaultUser{adminRoot}, []VaultUser{})
		assert.Empty(t, result)
	})
}

func TestResolveVaultPath_Wildcard(t *testing.T) {
	app := fiber.New()
	var resolvedPath string
	var shouldValidate bool
	var resolvedErr error

	app.Get("/vault/:vaultId/*", func(c *fiber.Ctx) error {
		resolvedPath, shouldValidate, resolvedErr = ResolveVaultPath(c)
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/vault/1/documents/file.txt", nil)
	app.Test(req)

	assert.NoError(t, resolvedErr)
	assert.Equal(t, "/documents/file.txt", resolvedPath)
	assert.True(t, shouldValidate)
}

func TestResolveVaultPath_WildcardRootDir(t *testing.T) {
	app := fiber.New()
	var resolvedPath string
	var shouldValidate bool
	var resolvedErr error

	app.Get("/vault/:vaultId/*", func(c *fiber.Ctx) error {
		resolvedPath, shouldValidate, resolvedErr = ResolveVaultPath(c)
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/vault/1/", nil)
	app.Test(req)

	assert.NoError(t, resolvedErr)
	assert.Equal(t, "/", resolvedPath)
	assert.True(t, shouldValidate)
}

func TestResolveVaultPath_QueryParam(t *testing.T) {
	app := fiber.New()
	var resolvedPath string
	var shouldValidate bool
	var resolvedErr error

	// Route with NO wildcard
	app.Get("/vault/:vaultId/list", func(c *fiber.Ctx) error {
		resolvedPath, shouldValidate, resolvedErr = ResolveVaultPath(c)
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/vault/1/list?path=/uploads/image.png", nil)
	app.Test(req)

	// Should resolve from query param
	assert.NoError(t, resolvedErr)
	assert.Equal(t, "/uploads/image.png", resolvedPath)
	assert.True(t, shouldValidate)
}

func TestResolveVaultPath_JSONBody(t *testing.T) {
	app := fiber.New()
	var resolvedPath string
	var shouldValidate bool
	var resolvedErr error

	app.Post("/vault/:vaultId/create", func(c *fiber.Ctx) error {
		resolvedPath, shouldValidate, resolvedErr = ResolveVaultPath(c)
		return c.SendString("ok")
	})

	body := map[string]string{"path": "/uploads/newfile.txt"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/vault/1/create", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	app.Test(req)

	// Should resolve from JSON body
	assert.NoError(t, resolvedErr)
	assert.Equal(t, "/uploads/newfile.txt", resolvedPath)
	assert.True(t, shouldValidate)
}

func TestResolveVaultPath_NoPath(t *testing.T) {
	app := fiber.New()
	var resolvedPath string
	var shouldValidate bool
	var resolvedErr error

	app.Get("/vault/:vaultId/info", func(c *fiber.Ctx) error {
		resolvedPath, shouldValidate, resolvedErr = ResolveVaultPath(c)
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/vault/1/info", nil)
	app.Test(req)

	assert.NoError(t, resolvedErr)
	assert.Equal(t, "", resolvedPath)
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
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		// Normal paths
		{"absolute path", "/documents/file.txt", "/documents/file.txt", false},
		{"relative path", "documents/file.txt", "/documents/file.txt", false},
		{"root path", "/", "/", false},

		// Path traversal attempts - path.Clean normalizes these
		{"traversal to admin", "/documents/../admin", "/admin", false},
		{"parent directory", "..", "/", false},
		{"relative traversal", "../admin", "/admin", false},

		// Error cases
		{"invalid encoding", "%ZZ", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := cleanPath(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
