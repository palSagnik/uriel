package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/palSagnik/uriel/internal/config"
	"github.com/palSagnik/uriel/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo AuthRepository
}

func NewService(repo AuthRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) RegisterPlayerService(ctx context.Context, req *models.RegisterRequest) (*models.Player, error) {

	// check if this username already exists
	existingPlayerByUsername, err := s.repo.GetPlayerByUsername(ctx, req.Username)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("service: error checking existing username %w", err)
	}
	if existingPlayerByUsername != nil {
		return nil, errors.New("username already exists")
	}

	// check if this email already exists
	existingPlayerByEmail, err := s.repo.GetPlayerByEmail(ctx, req.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("service: error checking existing email %w", err)
	}
	if existingPlayerByEmail != nil {
		return nil, errors.New("email already exists")
	}

	// hash password
	// does not accept more than 72 bytes
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("service: error hashing password %w", err)
	}

	// create player
	// TODO: Errors should be ENUMS
	newPlayer := models.Player{
		Username: req.Username,
		Email: req.Email,
		Password: string(hashedPassword),
		Role: config.PLAYER,
		IsOnline: false,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.repo.CreatePlayer(ctx, newPlayer); err != nil {
		return nil, fmt.Errorf("service: error in creating new player %w", err)
	}

	return &newPlayer, nil
}
