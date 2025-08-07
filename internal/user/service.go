package user

import "context"

type Service struct {
	repo UserRepository
}

func NewService(repo UserRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) UpdateMetadata(ctx context.Context, avatarUrl string) (string, error) {
	return "", nil
}