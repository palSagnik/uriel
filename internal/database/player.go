package database

import (
	"context"

	"github.com/palSagnik/uriel/internal/config"
	"github.com/palSagnik/uriel/internal/player"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoPlayerRepository struct {
	collection *mongo.Collection
}

func NewPlayerRepository(mongo *MongoDB) player.PlayerRepository {
	playerCollection := mongo.GetCollection(config.PLAYER_COLLECTION)

	return &mongoPlayerRepository{collection: playerCollection}
}

func (repo *mongoPlayerRepository) GetPlayerLocations(ctx context.Context) error {

	return nil
}