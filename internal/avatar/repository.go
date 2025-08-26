package avatar

import "context"

type AvatarRepository interface {
	GetAvatarUrlById(ctx context.Context, id string) (string, error)
}