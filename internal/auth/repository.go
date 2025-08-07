package auth

import (
	"context"

	"github.com/palSagnik/uriel/internal/models"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUserStatus(ctx context.Context, id string) error
}