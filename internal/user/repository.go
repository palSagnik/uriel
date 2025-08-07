package user

import "context"

type UserRepository interface {
	GetUserLocations(ctx context.Context) error
}
