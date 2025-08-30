package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/palSagnik/uriel/internal/avatar"
	"github.com/palSagnik/uriel/internal/config"
	"github.com/palSagnik/uriel/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoAvatarRepository struct {
	collection *mongo.Collection
}

func NewAvatarRepository(mongo *MongoDB) avatar.AvatarRepository {
	avatarCollection := mongo.GetCollection(config.AVATAR_COLLECTION)
	return &mongoAvatarRepository{collection: avatarCollection}
}

func (repo *mongoAvatarRepository) GetAvatarUrlById(ctx context.Context, id string) (string, error) {
	avatarId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	var avatar models.Avatar
	filter := bson.M{"_id": avatarId}
	if err := repo.collection.FindOne(ctx, filter).Decode(&avatar); err != nil {
		if err == mongo.ErrNoDocuments {
			return "", errors.New("avatar not found")
		}
		return "", fmt.Errorf("failed to find avatar: %v", err)
	}

	return avatar.AvatarUrl, nil
}

func (repo *mongoAvatarRepository) GetAvatars(ctx context.Context) ([]models.Avatar, error) {
	var avatars []models.Avatar

	cursor, err := repo.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.New("avatar list not found")
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &avatars); err != nil {
		return nil, fmt.Errorf("failed to decode cursor: %w", err)
	}

	return avatars, nil
}