package user

import (
	"context"

	"github.com/palSagnik/uriel/internal/models"
)

type UserRepository interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserLocations(ctx context.Context) error
	UpdateAvatar(ctx context.Context, id string, avatarUrl string) error
}
