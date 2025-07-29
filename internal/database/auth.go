package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/palSagnik/uriel/internal/auth"
	"github.com/palSagnik/uriel/internal/config"
	"github.com/palSagnik/uriel/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoAuthRepository struct {
	collection *mongo.Collection
}

func NewAuthRepository(mongodb *MongoDB) auth.AuthRepository {
	var err error

	playerCollection := mongodb.GetCollection(config.PLAYER_COLLECTION)

	// ensuring proper indexes efficient login and preventing duplicates
	// this helps in data integrity
	// USERNAME (INDEX)
	usernameIndexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	// EMAIL (INDEX)
	emailIndexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// creating index on username
	_, err = playerCollection.Indexes().CreateOne(ctx, usernameIndexModel)
	if err != nil {
		log.Printf("Warning: The unique index on username could not be created: %v", err)
	}

	// creating index on email
	_, err = playerCollection.Indexes().CreateOne(ctx, emailIndexModel)
	if err != nil {
		log.Printf("Warning: The unique index on email could not be created: %v", err)
	}
	return &mongoAuthRepository{collection: playerCollection}
}

// MongoAuthRepository implementing AuthRepository interface
func (repo *mongoAuthRepository) CreatePlayer(ctx context.Context, player models.Player) error {
	_, err := repo.collection.InsertOne(ctx, player)
	return err
}

func (repo *mongoAuthRepository) GetPlayerByUsername(ctx context.Context, username string) (*models.Player, error) {
	var player models.Player

	filter := bson.M{"username": username}
	err := repo.collection.FindOne(ctx, filter).Decode(&player)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &player, nil
}

func (repo *mongoAuthRepository) GetPlayerById(ctx context.Context, id string) (*models.Player, error) {
	var player models.Player

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id %v", err)
	}

	filter := bson.M{"_id": objectId}
	err = repo.collection.FindOne(ctx, filter).Decode(&player)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &player, nil
}

func (repo *mongoAuthRepository) GetPlayerByEmail(ctx context.Context, email string) (*models.Player, error) {
	var player models.Player

	filter := bson.M{"email": email}
	err := repo.collection.FindOne(ctx, filter).Decode(&player)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &player, nil
}

func (repo *mongoAuthRepository) UpdatePlayerStatus(ctx context.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectId}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "is_online", Value: true}}}}
	
	_, err = repo.collection.UpdateOne(ctx, filter, update)
	return err
}