package database

import (
	"context"

	"github.com/palSagnik/uriel/internal/config"
	"github.com/palSagnik/uriel/internal/user"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoUserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(mongo *MongoDB) user.UserRepository {
	userCollection := mongo.GetCollection(config.USER_COLLECTION)

	return &mongoUserRepository{collection: userCollection}
}

func (repo *mongoUserRepository) GetUserLocations(ctx context.Context) error {

	return nil
}
