package user

import "context"

type UserRepository interface {
	GetuserLocations(ctx context.Context) error
}
