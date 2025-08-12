package user

import (
	"context"

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

func (s *Service) UpdateAvatar(ctx context.Context, userId string, avatarId string) (string, error) {
	avatarUrl, err := s.avatarRepo.GetAvatarUrlById(ctx, avatarId)
	if err != nil {
		return "failed to get avatar url", err
	}

	if err := s.userRepo.UpdateAvatar(ctx, userId, avatarUrl); err != nil {
		return "failed to update avatar", err
	}

	return "updated avatar succesfully", nil
}
