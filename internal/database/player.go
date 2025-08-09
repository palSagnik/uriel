package database

import (
	"context"

	"github.com/palSagnik/uriel/internal/config"
	"github.com/palSagnik/uriel/internal/user"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongouserRepository struct {
	collection *mongo.Collection
}

func NewuserRepository(mongo *MongoDB) user.UserRepository {
	userCollection := mongo.GetCollection(config.USER_COLLECTION)

	return &mongouserRepository{collection: userCollection}
}

func (repo *mongouserRepository) GetuserLocations(ctx context.Context) error {

	return nil
}
