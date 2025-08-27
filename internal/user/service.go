package user

import (
	"context"

	"github.com/palSagnik/uriel/internal/avatar"
	"github.com/palSagnik/uriel/internal/models"
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

func (s *Service) UpdateUserAvatar(ctx context.Context, userId string, avatarId string) (string, error) {
	avatarUrl, err := s.avatarRepo.GetAvatarUrlById(ctx, avatarId)
	if err != nil {
		return "failed to get avatar url", err
	}

	if err := s.userRepo.UpdateUserAvatar(ctx, userId, avatarUrl); err != nil {
		return "failed to update avatar", err
	}

	return "updated avatar succesfully", nil
}

func (s *Service) GetAvatars(ctx context.Context) ([]models.Avatar, error) {
	avatars, err := s.avatarRepo.GetAvatars(ctx)
	if err != nil {
		return nil, err
	}

	return avatars, nil
}

func (s *Service) GetUsers(ctx context.Context) ([]models.User, error) {
	users, err := s.userRepo.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}
