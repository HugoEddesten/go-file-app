package users

import (
	"context"
	"testing"

	"go-file-api/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UsersRepositoryTestSuite struct {
	suite.Suite
	pgContainer *testutil.PostgresContainer
	repo        *Repository
	ctx         context.Context
}

func (s *UsersRepositoryTestSuite) SetupSuite() {
	s.ctx = context.Background()

	// Start PostgreSQL container
	pgContainer, err := testutil.SetupPostgres(s.ctx)
	require.NoError(s.T(), err, "failed to setup postgres container")

	s.pgContainer = pgContainer
	s.repo = &Repository{DB: pgContainer.Pool}
}

func (s *UsersRepositoryTestSuite) TearDownSuite() {
	if s.pgContainer != nil {
		err := s.pgContainer.Teardown(s.ctx)
		assert.NoError(s.T(), err, "failed to teardown postgres container")
	}
}

func (s *UsersRepositoryTestSuite) SetupTest() {
	// Clean tables before each test
	err := s.pgContainer.CleanupTables(s.ctx)
	require.NoError(s.T(), err, "failed to cleanup tables")
}

// Test Create
func (s *UsersRepositoryTestSuite) TestCreate_Success() {
	email := "test@example.com"
	passwordHash := "hashed_password_123"

	// Create user
	userId, err := s.repo.Create(email, passwordHash)

	// Assertions
	assert.NoError(s.T(), err)
	assert.Greater(s.T(), userId, 0)

	// Verify user exists in database
	user, err := s.repo.FindByEmail(email)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), userId, user.Id)
	assert.Equal(s.T(), email, user.Email)
	assert.Equal(s.T(), passwordHash, user.PasswordHash)
}

func (s *UsersRepositoryTestSuite) TestCreate_DuplicateEmail_ReturnsError() {
	email := "duplicate@example.com"
	passwordHash := "hashed_password_123"

	// Create first user
	_, err := s.repo.Create(email, passwordHash)
	require.NoError(s.T(), err)

	// Try to create duplicate
	_, err = s.repo.Create(email, passwordHash)

	// Should return error due to unique constraint
	assert.Error(s.T(), err)
}

// Test FindByEmail
func (s *UsersRepositoryTestSuite) TestFindByEmail_UserExists() {
	email := "find@example.com"
	passwordHash := "hashed_password_456"

	// Create user first
	userId, err := s.repo.Create(email, passwordHash)
	require.NoError(s.T(), err)

	// Find user
	user, err := s.repo.FindByEmail(email)

	// Assertions
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), user)
	assert.Equal(s.T(), userId, user.Id)
	assert.Equal(s.T(), email, user.Email)
	assert.Equal(s.T(), passwordHash, user.PasswordHash)
}

func (s *UsersRepositoryTestSuite) TestFindByEmail_UserDoesNotExist() {
	email := "nonexistent@example.com"

	// Try to find non-existent user
	user, err := s.repo.FindByEmail(email)

	// Should return error (pgx.ErrNoRows)
	assert.Error(s.T(), err)
	assert.NotNil(s.T(), user) // The function still returns a User struct
}

func TestUsersRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UsersRepositoryTestSuite))
}
