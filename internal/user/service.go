package user

type Service struct {
	repo UserRepository
}

func NewService(userRepo UserRepository) *Service {
	return &Service{
		repo: userRepo,
	}
}
