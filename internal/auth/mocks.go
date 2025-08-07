package auth

import (
	"context"

	"github.com/palSagnik/uriel/internal/models"
	"github.com/stretchr/testify/mock"
)

// Mocking Auth Repository
type MockAuthRepository struct {
	mock.Mock
}

// Mocking the repository methods
// CreateUser(ctx context.Context, user models.User) error
func (m *MockAuthRepository) CreateUser(ctx context.Context, user models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// GetUserByUsername(ctx context.Context, username string) (*models.User, error)
func (m *MockAuthRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.User), args.Error(1)
}

// GetUserByEmail(ctx context.Context, email string) (*models.User, error)
func (m *MockAuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.User), args.Error(1)
}

// GetUserById(ctx context.Context, id string) (*models.User, error)
func (m *MockAuthRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	return nil, nil
}

// UpdateUserStatus(ctx context.Context, id string) error
func (m *MockAuthRepository) UpdateUserStatus(ctx context.Context, id string) error {
	return nil
}

