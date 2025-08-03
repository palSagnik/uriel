package auth

import (
	"context"

	"github.com/palSagnik/uriel/internal/models"
)

type AuthRepository interface {
	Createuser(ctx context.Context, user models.User) error
	GetuserByUsername(ctx context.Context, username string) (*models.User, error)
	GetuserById(ctx context.Context, id string) (*models.User, error)
	GetuserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateuserStatus(ctx context.Context, id string) error
}
