package database

import (
	"context"
	"log"
	"time"

	"github.com/palSagnik/uriel/internal/auth"
	"github.com/palSagnik/uriel/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoAuthRepository struct {
	collection *mongo.Collection
}

func NewMongoAuthRepository(client *mongo.Client) auth.AuthRepository {
	playerCollection := client.Database("uriel").Collection("player")

	// ensuring the index on username for efficient login and preventing duplicates
	// this helps in data integrity
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "username", Value: "1"}},
		Options: options.Index().SetUnique(true),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	_, err := playerCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("Warning: The unique index on username could not be created: %v", err)
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

	filter := bson.M{"player_id": id}
	err := repo.collection.FindOne(ctx, filter).Decode(&player)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err 
	}
	return &player, nil
}