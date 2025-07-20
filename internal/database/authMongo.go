package database

import (
	"context"

	"github.com/palSagnik/uriel/internal/auth"
	"github.com/palSagnik/uriel/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoAuthRepository struct {
	collection *mongo.Collection
}

func newAuthRepository(client *mongo.Client) auth.AuthRepository {
	playerCollection := client.Database("uriel").Collection("player")

	return &mongoAuthRepository{
		collection: playerCollection,
	}
}

// MongoAuthRepository implementing AuthRepository interface
func (repo *mongoAuthRepository) CreatePlayer(ctx context.Context, player models.Player) error {
	return nil
}

func (repo *mongoAuthRepository) GetPlayerByUsername(ctx context.Context, username string) (*models.Player, error) {
	return nil, nil
}

func (repo *mongoAuthRepository) GetPlayerById(ctx context.Context, id string) (*models.Player, error) {
	return nil, nil
}