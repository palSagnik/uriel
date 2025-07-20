package auth

import (
	"context"

	"github.com/palSagnik/uriel/internal/models"
)

type AuthRepository interface {
	CreatePlayer(ctx context.Context, player models.Player) error
	GetPlayerByUsername(ctx context.Context, username string) (*models.Player, error)
	GetPlayerById(ctx context.Context, id string) (*models.Player, error)
}