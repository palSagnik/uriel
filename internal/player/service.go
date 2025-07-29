package player

type Service struct{
	repo PlayerRepository
}

func NewService(playerRepo PlayerRepository) *Service {
	return &Service{
		repo: playerRepo,
	}
}