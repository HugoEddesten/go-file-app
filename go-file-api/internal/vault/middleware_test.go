package vault

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// MockVaultRepository implements VaultUserGetter for testing
type MockVaultRepository struct {
	vaultUsers map[int]map[int][]VaultUser // map[vaultId][userId][]VaultUser
}

func NewMockVaultRepository() *MockVaultRepository {
	return &MockVaultRepository{
		vaultUsers: make(map[int]map[int][]VaultUser),
	}
}

func (m *MockVaultRepository) AddVaultUser(vaultId, userId int, vu VaultUser) {
	if m.vaultUsers[vaultId] == nil {
		m.vaultUsers[vaultId] = make(map[int][]VaultUser)
	}
	m.vaultUsers[vaultId][userId] = append(m.vaultUsers[vaultId][userId], vu)
}

func (m *MockVaultRepository) GetVaultUsers(ctx context.Context, vaultId, userId int) ([]VaultUser, error) {
	if vault, exists := m.vaultUsers[vaultId]; exists {
		if users, hasUser := vault[userId]; hasUser {
			return users, nil
		}
	}
	return []VaultUser{}, nil
}

type VaultAccessMiddlewareTestSuite struct {
	suite.Suite
	app      *fiber.App
	mockRepo *MockVaultRepository
}

func (s *VaultAccessMiddlewareTestSuite) SetupTest() {
	s.app = fiber.New()
	s.mockRepo = NewMockVaultRepository()
}

func (s *VaultAccessMiddlewareTestSuite) setupRoute(requiredRole VaultRole, pathInRoute string) {
	s.app = fiber.New()

	// Use the actual production VaultAccessMiddleware with our mock
	middleware := VaultAccessMiddleware(s.mockRepo, requiredRole)

	s.app.Get(pathInRoute,
		func(c *fiber.Ctx) error {
			// Simulate JWT middleware setting userId
			c.Locals("userId", 1)
			return c.Next()
		},
		middleware,
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{
				"success":   true,
				"vaultId":   c.Locals("vaultId"),
				"vaultRole": c.Locals("vaultRole"),
			})
		},
	)
}

func (s *VaultAccessMiddlewareTestSuite) TestUserHasRequiredRole_FullAccess_ShouldAllow() {
	s.setupRoute(VaultRoleEditor, "/vault/:vaultId/files")

	// User 1 has Editor role with full vault access
	s.mockRepo.AddVaultUser(10, 1, VaultUser{
		VaultId: 10,
		UserId:  1,
		Role:    VaultRoleEditor,
		Path:    "/",
	})

	req := httptest.NewRequest("GET", "/vault/10/files", nil)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}

func (s *VaultAccessMiddlewareTestSuite) TestUserHasInsufficientRole_ShouldDeny() {
	s.setupRoute(VaultRoleEditor, "/vault/:vaultId/files")

	// User 1 only has Viewer role, but Editor is required
	s.mockRepo.AddVaultUser(10, 1, VaultUser{
		VaultId: 10,
		UserId:  1,
		Role:    VaultRoleViewer,
		Path:    "/",
	})

	req := httptest.NewRequest("GET", "/vault/10/files", nil)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusForbidden, resp.StatusCode)
}

func (s *VaultAccessMiddlewareTestSuite) TestUserHasHigherRole_ShouldAllow() {
	s.setupRoute(VaultRoleEditor, "/vault/:vaultId/files")

	// User 1 has Owner role (higher than required Editor)
	s.mockRepo.AddVaultUser(10, 1, VaultUser{
		VaultId: 10,
		UserId:  1,
		Role:    VaultRoleOwner,
		Path:    "/",
	})

	req := httptest.NewRequest("GET", "/vault/10/files", nil)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}

func (s *VaultAccessMiddlewareTestSuite) TestUserRestrictedToPath_AccessingAllowedPath_ShouldAllow() {
	s.app = fiber.New()
	middleware := VaultAccessMiddleware(s.mockRepo, VaultRoleEditor)

	s.app.Get("/vault/:vaultId/files/*",
		func(c *fiber.Ctx) error {
			c.Locals("userId", 1)
			return c.Next()
		},
		middleware,
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"success": true})
		},
	)

	// User 1 restricted to /documents
	s.mockRepo.AddVaultUser(10, 1, VaultUser{
		VaultId: 10,
		UserId:  1,
		Role:    VaultRoleEditor,
		Path:    "/documents",
	})

	req := httptest.NewRequest("GET", "/vault/10/files/documents/report.pdf", nil)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}

func (s *VaultAccessMiddlewareTestSuite) TestUserRestrictedToPath_AccessingDeniedPath_ShouldDeny() {
	s.app = fiber.New()
	middleware := VaultAccessMiddleware(s.mockRepo, VaultRoleEditor)

	s.app.Get("/vault/:vaultId/files/*",
		func(c *fiber.Ctx) error {
			c.Locals("userId", 1)
			return c.Next()
		},
		middleware,
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"success": true})
		},
	)

	// User 1 restricted to /documents
	s.mockRepo.AddVaultUser(10, 1, VaultUser{
		VaultId: 10,
		UserId:  1,
		Role:    VaultRoleEditor,
		Path:    "/documents",
	})

	// Trying to access /admin (outside allowed path)
	req := httptest.NewRequest("GET", "/vault/10/files/admin/config.txt", nil)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusForbidden, resp.StatusCode)
}

func (s *VaultAccessMiddlewareTestSuite) TestUserNotInVault_ShouldDeny() {
	s.setupRoute(VaultRoleViewer, "/vault/:vaultId/files")

	// No vault users added for user 1 in vault 10

	req := httptest.NewRequest("GET", "/vault/10/files", nil)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusForbidden, resp.StatusCode)
}

func (s *VaultAccessMiddlewareTestSuite) TestInvalidVaultId_ShouldReturn400() {
	s.setupRoute(VaultRoleViewer, "/vault/:vaultId/files")

	req := httptest.NewRequest("GET", "/vault/invalid/files", nil)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusBadRequest, resp.StatusCode)
}

func (s *VaultAccessMiddlewareTestSuite) TestMultipleVaultUserEntries_OneMatches_ShouldAllow() {
	s.app = fiber.New()
	middleware := VaultAccessMiddleware(s.mockRepo, VaultRoleEditor)

	s.app.Get("/vault/:vaultId/files/*",
		func(c *fiber.Ctx) error {
			c.Locals("userId", 1)
			return c.Next()
		},
		middleware,
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"success": true})
		},
	)

	// User has two entries: one for /documents (insufficient role), one for /public (sufficient role)
	s.mockRepo.AddVaultUser(10, 1, VaultUser{
		VaultId: 10,
		UserId:  1,
		Role:    VaultRoleViewer, // Too low
		Path:    "/documents",
	})
	s.mockRepo.AddVaultUser(10, 1, VaultUser{
		VaultId: 10,
		UserId:  1,
		Role:    VaultRoleEditor, // Sufficient
		Path:    "/public",
	})

	// Accessing /public with Editor requirement - should work
	req := httptest.NewRequest("GET", "/vault/10/files/public/data.json", nil)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}

func (s *VaultAccessMiddlewareTestSuite) TestPathTraversalAttempt_ShouldDeny() {
	s.app = fiber.New()
	middleware := VaultAccessMiddleware(s.mockRepo, VaultRoleEditor)

	s.app.Get("/vault/:vaultId/files/*",
		func(c *fiber.Ctx) error {
			c.Locals("userId", 1)
			return c.Next()
		},
		middleware,
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"success": true})
		},
	)

	// User restricted to /documents
	s.mockRepo.AddVaultUser(10, 1, VaultUser{
		VaultId: 10,
		UserId:  1,
		Role:    VaultRoleEditor,
		Path:    "/documents",
	})

	// Try to escape using ../
	req := httptest.NewRequest("GET", "/vault/10/files/documents/../admin/secret.txt", nil)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	// cleanPath should normalize this, and it won't match /documents
	assert.Equal(s.T(), fiber.StatusForbidden, resp.StatusCode)
}

func (s *VaultAccessMiddlewareTestSuite) TestRouteWithoutWildcard_NoPathValidation() {
	s.app = fiber.New()
	middleware := VaultAccessMiddleware(s.mockRepo, VaultRoleEditor)

	// Route without wildcard and no path query/body - shouldValidate=false
	s.app.Get("/vault/:vaultId/list",
		func(c *fiber.Ctx) error {
			c.Locals("userId", 1)
			return c.Next()
		},
		middleware,
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"success": true})
		},
	)

	// User restricted to /uploads (but path validation won't happen)
	s.mockRepo.AddVaultUser(10, 1, VaultUser{
		VaultId: 10,
		UserId:  1,
		Role:    VaultRoleEditor,
		Path:    "/uploads",
	})

	// Works because shouldValidate=false, so path restriction is not checked
	req := httptest.NewRequest("GET", "/vault/10/list", nil)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}

func (s *VaultAccessMiddlewareTestSuite) TestPathViaQueryParam_ShouldWork() {
	s.app = fiber.New()
	middleware := VaultAccessMiddleware(s.mockRepo, VaultRoleEditor)

	// Route without wildcard - should resolve path from query param
	s.app.Get("/vault/:vaultId/list",
		func(c *fiber.Ctx) error {
			c.Locals("userId", 1)
			return c.Next()
		},
		middleware,
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"success": true})
		},
	)

	s.mockRepo.AddVaultUser(10, 1, VaultUser{
		VaultId: 10,
		UserId:  1,
		Role:    VaultRoleEditor,
		Path:    "/uploads",
	})

	// Query param path matches allowed path
	req := httptest.NewRequest("GET", "/vault/10/list?path=/uploads/image.png", nil)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}

func (s *VaultAccessMiddlewareTestSuite) TestPathViaJSONBody_ShouldWork() {
	s.app = fiber.New()
	middleware := VaultAccessMiddleware(s.mockRepo, VaultRoleEditor)

	// Route without wildcard - should resolve path from JSON body
	s.app.Post("/vault/:vaultId/create",
		func(c *fiber.Ctx) error {
			c.Locals("userId", 1)
			return c.Next()
		},
		middleware,
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"success": true})
		},
	)

	s.mockRepo.AddVaultUser(10, 1, VaultUser{
		VaultId: 10,
		UserId:  1,
		Role:    VaultRoleEditor,
		Path:    "/uploads",
	})

	body := map[string]string{"path": "/uploads/newfile.txt"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/vault/10/create", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}

func (s *VaultAccessMiddlewareTestSuite) TestExactPathMatch_ShouldAllow() {
	s.app = fiber.New()
	middleware := VaultAccessMiddleware(s.mockRepo, VaultRoleEditor)

	s.app.Get("/vault/:vaultId/files/*",
		func(c *fiber.Ctx) error {
			c.Locals("userId", 1)
			return c.Next()
		},
		middleware,
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"success": true})
		},
	)

	// User restricted to /documents
	s.mockRepo.AddVaultUser(10, 1, VaultUser{
		VaultId: 10,
		UserId:  1,
		Role:    VaultRoleEditor,
		Path:    "/documents",
	})

	// Accessing exactly /documents (not a subdirectory)
	req := httptest.NewRequest("GET", "/vault/10/files/documents", nil)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
}

func (s *VaultAccessMiddlewareTestSuite) TestSimilarPathPrefix_ShouldNotMatch() {
	s.app = fiber.New()
	middleware := VaultAccessMiddleware(s.mockRepo, VaultRoleEditor)

	s.app.Get("/vault/:vaultId/files/*",
		func(c *fiber.Ctx) error {
			c.Locals("userId", 1)
			return c.Next()
		},
		middleware,
		func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"success": true})
		},
	)

	// User restricted to /doc
	s.mockRepo.AddVaultUser(10, 1, VaultUser{
		VaultId: 10,
		UserId:  1,
		Role:    VaultRoleEditor,
		Path:    "/doc",
	})

	// Trying to access /documents (starts with /doc but not a subdirectory)
	req := httptest.NewRequest("GET", "/vault/10/files/documents/file.txt", nil)
	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusForbidden, resp.StatusCode)
}

func TestVaultAccessMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(VaultAccessMiddlewareTestSuite))
}
