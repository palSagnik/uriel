package user

import (
	"context"

	"github.com/palSagnik/uriel/internal/models"
)

type UserRepository interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	UpdateUserAvatar(ctx context.Context, id string, avatarUrl string) error
}
