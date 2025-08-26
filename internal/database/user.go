package database

import (
	"context"

	"github.com/palSagnik/uriel/internal/config"
	"github.com/palSagnik/uriel/internal/models"
	"github.com/palSagnik/uriel/internal/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoUserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(mongo *MongoDB) user.UserRepository {
	userCollection := mongo.GetCollection(config.USER_COLLECTION)

	return &mongoUserRepository{collection: userCollection}
}

func (repo *mongoUserRepository) UpdateAvatar(ctx context.Context, userId string, avatarUrl string) error {
	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	filterUser := bson.M{"_id": userObjectId}
	updateUser := bson.D{{Key: "$set", Value: bson.D{{Key: "avatar_url", Value: avatarUrl}}}}

	_, err = repo.collection.UpdateOne(ctx, filterUser, updateUser)
	return err
}

func (repo *mongoUserRepository) GetUserLocations(ctx context.Context) error {
	return nil
}

func (repo *mongoUserRepository) GetUsers(ctx context.Context) ([]models.User, error) {
	return nil, nil
}
