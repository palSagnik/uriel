package user

import (
	"context"
	"fmt"

	"github.com/palSagnik/uriel/internal/avatar"
)

type Service struct {
	userRepo   UserRepository
	avatarRepo avatar.AvatarRepository
}

func NewService(userRepo UserRepository, avatarRepo avatar.AvatarRepository) *Service {
	return &Service{
		userRepo:   userRepo,
		avatarRepo: avatarRepo,
	}
}

func (s *Service) UpdateMetadata(ctx context.Context, userId string, avatarId string) (string, error) {
	avatarUrl, err := s.avatarRepo.GetAvatarUrlById(ctx, avatarId)
	if err != nil {
		return fmt.Sprintf("failed to get avatar url"), err
	}


	return "", nil
}
