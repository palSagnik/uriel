package user

import (
	"context"

	"github.com/palSagnik/uriel/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

type MockAvatarRepository struct {
	mock.Mock
}

// Mocking user repository methods
// GetUsers(ctx context.Context) ([]models.User, error)
func (m *MockUserRepository) GetUsers(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]models.User), args.Error(1)
}

// UpdateUserAvatar(ctx context.Context, id string, avatarUrl string) error
func (m *MockUserRepository) UpdateUserAvatar(ctx context.Context, id, avatarUrl string) error {
	args := m.Called(ctx, id, avatarUrl)
	return args.Error(0)
}

// Mocking avatar repository methods
// GetAvatarUrlById(ctx context.Context, id string) (string, error)
func (m *MockAvatarRepository) GetAvatarUrlById(ctx context.Context, id string) (string, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}

	return args.Get(0).(string), args.Error(1)
}

// GetAvatars(ctx context.Context) ([]models.Avatar, error)
func (m *MockAvatarRepository) GetAvatars(ctx context.Context) ([]models.Avatar, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]models.Avatar), args.Error(1)
}