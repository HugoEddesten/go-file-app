package vault

import (
	"context"
	"testing"

	"go-file-api/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type VaultRepositoryTestSuite struct {
	suite.Suite
	pgContainer *testutil.PostgresContainer
	repo        *Repository
	ctx         context.Context
}

func (s *VaultRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()

	// Start PostgreSQL container
	pgContainer, err := testutil.SetupPostgres(s.ctx)
	require.NoError(s.T(), err, "failed to setup postgres container")

	s.pgContainer = pgContainer
	s.repo = &Repository{DB: pgContainer.Pool}
}

func (s *VaultRepositoryTestSuite) TearDownSuite() {
	if s.pgContainer != nil {
		err := s.pgContainer.Teardown(s.ctx)
		assert.NoError(s.T(), err, "failed to teardown postgres container")
	}
}

func (s *VaultRepositoryTestSuite) SetupTest() {
	// Clean tables before each test
	err := s.pgContainer.CleanupTables(s.ctx)
	require.NoError(s.T(), err, "failed to cleanup tables")
}

// Helper to create a test user
func (s *VaultRepositoryTestSuite) createTestUser(email string) int {
	var userId int
	err := s.pgContainer.Pool.QueryRow(s.ctx,
		`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`,
		email, "test_hash",
	).Scan(&userId)
	require.NoError(s.T(), err)
	return userId
}

// Test Create
func (s *VaultRepositoryTestSuite) TestCreate_Success() {
	userId := s.createTestUser("owner@example.com")
	vaultName := "My Test Vault"

	// Create vault
	vault, err := s.repo.Create(s.ctx, vaultName, userId)

	// Assertions
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), vault)
	assert.Greater(s.T(), vault.Id, 0)
	assert.Equal(s.T(), vaultName, vault.Name)
	assert.NotZero(s.T(), vault.CreatedAt)
	assert.NotZero(s.T(), vault.UpdatedAt)

	// Verify vault_users entry was created with owner role
	vaultUsers, err := s.repo.GetVaultUsers(s.ctx, vault.Id, userId)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), vaultUsers, 1)
	assert.Equal(s.T(), VaultRoleOwner, vaultUsers[0].Role)
	assert.Equal(s.T(), "/", vaultUsers[0].Path)
}

func (s *VaultRepositoryTestSuite) TestCreate_TransactionRollback() {
	// Use invalid userId to trigger rollback
	invalidUserId := 99999

	vault, err := s.repo.Create(s.ctx, "Test Vault", invalidUserId)

	// Should fail due to foreign key constraint
	assert.Error(s.T(), err)
	assert.Nil(s.T(), vault)

	// Verify no vault was created
	var count int
	err = s.pgContainer.Pool.QueryRow(s.ctx, `SELECT COUNT(*) FROM vaults`).Scan(&count)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, count)
}

// Test AddUserToVault
func (s *VaultRepositoryTestSuite) TestAddUserToVault_Success() {
	owner := s.createTestUser("owner@example.com")
	editor := s.createTestUser("editor@example.com")

	vault, err := s.repo.Create(s.ctx, "Shared Vault", owner)
	require.NoError(s.T(), err)

	// Add editor to vault
	vaultUser, err := s.repo.AddUserToVault(
		s.ctx,
		vault.Id,
		editor,
		"/documents",
		VaultRoleEditor,
	)

	// Assertions
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), vaultUser)
	assert.Equal(s.T(), vault.Id, vaultUser.VaultId)
	assert.Equal(s.T(), editor, vaultUser.UserId)
	assert.Equal(s.T(), "/documents", vaultUser.Path)
	assert.Equal(s.T(), VaultRoleEditor, vaultUser.Role)
}

func (s *VaultRepositoryTestSuite) TestAddUserToVault_DuplicateEntry() {
	owner := s.createTestUser("owner@example.com")
	vault, err := s.repo.Create(s.ctx, "Test Vault", owner)
	require.NoError(s.T(), err)

	// Try to add owner again with same path
	_, err = s.repo.AddUserToVault(s.ctx, vault.Id, owner, "/", VaultRoleEditor)

	// Should fail due to unique constraint
	assert.Error(s.T(), err)
}

// Test GetVaultUsers
func (s *VaultRepositoryTestSuite) TestGetVaultUsers_MultipleEntries() {
	owner := s.createTestUser("owner@example.com")
	vault, err := s.repo.Create(s.ctx, "Test Vault", owner)
	require.NoError(s.T(), err)

	// Add another access for same user with different path
	_, err = s.repo.AddUserToVault(s.ctx, vault.Id, owner, "/private", VaultRoleEditor)
	require.NoError(s.T(), err)

	// Get vault users
	vaultUsers, err := s.repo.GetVaultUsers(s.ctx, vault.Id, owner)

	// Should return both entries
	assert.NoError(s.T(), err)
	assert.Len(s.T(), vaultUsers, 2)
}

func (s *VaultRepositoryTestSuite) TestGetVaultUsers_NoAccess() {
	owner := s.createTestUser("owner@example.com")
	otherUser := s.createTestUser("other@example.com")
	vault, err := s.repo.Create(s.ctx, "Private Vault", owner)
	require.NoError(s.T(), err)

	// Try to get vault users for user without access
	vaultUsers, err := s.repo.GetVaultUsers(s.ctx, vault.Id, otherUser)

	assert.NoError(s.T(), err)
	assert.Empty(s.T(), vaultUsers)
}

func (s *VaultRepositoryTestSuite) TestGetVaultUsers_DatabaseError() {
	// Create a cancelled context to force a database error
	cancelledCtx, cancel := context.WithCancel(s.ctx)
	cancel() // Immediately cancel the context

	// Try to query with cancelled context
	vaultUsers, err := s.repo.GetVaultUsers(cancelledCtx, 1, 1)

	// Should return error due to cancelled context
	assert.Error(s.T(), err)
	assert.Nil(s.T(), vaultUsers)
}

// Test UpdateVaultUser
func (s *VaultRepositoryTestSuite) TestUpdateVaultUser_Success() {
	owner := s.createTestUser("owner@example.com")
	editor := s.createTestUser("editor@example.com")

	vault, err := s.repo.Create(s.ctx, "Test Vault", owner)
	require.NoError(s.T(), err)

	vaultUser, err := s.repo.AddUserToVault(s.ctx, vault.Id, editor, "/docs", VaultRoleViewer)
	require.NoError(s.T(), err)

	// Update role and path
	vaultUser.Role = VaultRoleEditor
	vaultUser.Path = "/documents"

	updated, err := s.repo.UpdateVaultUser(s.ctx, vaultUser)

	// Assertions
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), VaultRoleEditor, updated.Role)
	assert.Equal(s.T(), "/documents", updated.Path)
	assert.True(s.T(), updated.UpdatedAt.After(updated.CreatedAt))
}

func (s *VaultRepositoryTestSuite) TestUpdateVaultUser_NonExistent() {
	// Try to update a non-existent vault user
	nonExistentVaultUser := &VaultUser{
		Id:      99999,
		VaultId: 1,
		UserId:  1,
		Path:    "/test",
		Role:    VaultRoleEditor,
	}

	updated, err := s.repo.UpdateVaultUser(s.ctx, nonExistentVaultUser)

	// Should return error (no rows to update)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), updated)
}

// Test GetVault
func (s *VaultRepositoryTestSuite) TestGetVault_Success() {
	owner := s.createTestUser("owner@example.com")
	editor := s.createTestUser("editor@example.com")

	vault, err := s.repo.Create(s.ctx, "Multi-User Vault", owner)
	require.NoError(s.T(), err)

	// Add editor to vault
	_, err = s.repo.AddUserToVault(s.ctx, vault.Id, editor, "/docs", VaultRoleEditor)
	require.NoError(s.T(), err)

	// Get vault with all users
	vaultWithUsers, err := s.repo.GetVault(s.ctx, vault.Id)

	// Assertions
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), vaultWithUsers)
	assert.Equal(s.T(), vault.Id, vaultWithUsers.Id)
	assert.Equal(s.T(), "Multi-User Vault", vaultWithUsers.Name)
	assert.Len(s.T(), vaultWithUsers.Users, 2) // Owner and Editor
}

func (s *VaultRepositoryTestSuite) TestGetVault_AggregatesMultipleUsers() {
	owner := s.createTestUser("owner@example.com")
	user2 := s.createTestUser("user2@example.com")
	user3 := s.createTestUser("user3@example.com")

	vault, err := s.repo.Create(s.ctx, "Shared Vault", owner)
	require.NoError(s.T(), err)

	// Add multiple users
	_, err = s.repo.AddUserToVault(s.ctx, vault.Id, user2, "/", VaultRoleEditor)
	require.NoError(s.T(), err)
	_, err = s.repo.AddUserToVault(s.ctx, vault.Id, user3, "/public", VaultRoleViewer)
	require.NoError(s.T(), err)

	// Get vault
	vaultWithUsers, err := s.repo.GetVault(s.ctx, vault.Id)

	// Should properly aggregate all users
	assert.NoError(s.T(), err)
	assert.Len(s.T(), vaultWithUsers.Users, 3)

	// Verify each user has correct properties
	for _, user := range vaultWithUsers.Users {
		assert.NotZero(s.T(), user.Id)
		assert.NotEmpty(s.T(), user.Email)
		assert.NotZero(s.T(), user.Role)
		assert.NotEmpty(s.T(), user.Path)
	}
}

func (s *VaultRepositoryTestSuite) TestGetVault_NonExistent() {
	nonExistentVaultId := 99999

	// Try to get non-existent vault
	vault, err := s.repo.GetVault(s.ctx, nonExistentVaultId)

	// Should return error (ErrNoRows)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), vault)
}

// Test GetVaultsForUser
func (s *VaultRepositoryTestSuite) TestGetVaultsForUser_Success() {
	user := s.createTestUser("user@example.com")
	otherUser := s.createTestUser("other@example.com")

	// Create multiple vaults
	_, err := s.repo.Create(s.ctx, "Vault 1", user)
	require.NoError(s.T(), err)
	_, err = s.repo.Create(s.ctx, "Vault 2", user)
	require.NoError(s.T(), err)

	// Add user to a third vault owned by someone else
	vault3, err := s.repo.Create(s.ctx, "Vault 3", otherUser)
	require.NoError(s.T(), err)
	_, err = s.repo.AddUserToVault(s.ctx, vault3.Id, user, "/shared", VaultRoleEditor)
	require.NoError(s.T(), err)

	// Get all vaults for user
	vaults, err := s.repo.GetVaultsForUser(s.ctx, user)

	// Assertions
	assert.NoError(s.T(), err)
	assert.Len(s.T(), vaults, 3) // User has access to all 3 vaults

	// Verify vault names
	vaultNames := make(map[string]bool)
	for _, v := range vaults {
		vaultNames[v.Name] = true
	}
	assert.True(s.T(), vaultNames["Vault 1"])
	assert.True(s.T(), vaultNames["Vault 2"])
	assert.True(s.T(), vaultNames["Vault 3"])
}

func (s *VaultRepositoryTestSuite) TestGetVaultsForUser_AggregatesUsers() {
	owner := s.createTestUser("owner@example.com")
	editor := s.createTestUser("editor@example.com")
	viewer := s.createTestUser("viewer@example.com")

	// Create vault with multiple users
	vault, err := s.repo.Create(s.ctx, "Team Vault", owner)
	require.NoError(s.T(), err)
	_, err = s.repo.AddUserToVault(s.ctx, vault.Id, editor, "/", VaultRoleEditor)
	require.NoError(s.T(), err)
	_, err = s.repo.AddUserToVault(s.ctx, vault.Id, viewer, "/public", VaultRoleViewer)
	require.NoError(s.T(), err)

	// Get vaults for owner
	vaults, err := s.repo.GetVaultsForUser(s.ctx, owner)

	// Should return vault with all users properly aggregated
	assert.NoError(s.T(), err)
	require.Len(s.T(), vaults, 1)
	assert.Len(s.T(), vaults[0].Users, 3) // Owner, Editor, Viewer
}

func (s *VaultRepositoryTestSuite) TestGetVaultsForUser_NoVaults() {
	user := s.createTestUser("isolated@example.com")

	// Get vaults for user with no vault access
	vaults, err := s.repo.GetVaultsForUser(s.ctx, user)

	// Should return empty slice, no error
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), vaults)
}

func TestVaultRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(VaultRepositoryTestSuite))
}
