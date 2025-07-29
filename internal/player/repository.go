package player

import "context"

type PlayerRepository interface {
	GetPlayerLocations(ctx context.Context) error
}