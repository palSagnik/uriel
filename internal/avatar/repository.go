package avatar

import (
	"context"

	"github.com/palSagnik/uriel/internal/models"
)

type AvatarRepository interface {
	GetAvatarUrlById(ctx context.Context, id string) (string, error)
	GetAvatars(ctx context.Context) ([]models.Avatar, error)
}