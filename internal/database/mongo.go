package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/palSagnik/uriel/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
    Client *mongo.Client
}

func NewMongoClient(uri string) (*MongoDB, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

    clientOptions := options.Client().ApplyURI(uri)
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return nil, err
    }

    // Ping the primary to verify connection and credentials
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		// Ensure to disconnect if ping fails, as Connect might still return a client
		if disconnectErr := client.Disconnect(context.Background()); disconnectErr != nil {
			log.Printf("Error during MongoDB ping failure disconnect: %v", disconnectErr)
		}
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

   return &MongoDB{Client: client}, nil
}

func (mongo *MongoDB) GetCollection(collectionName string) *mongo.Collection {
    return mongo.Client.Database(config.DATABASE_NAME).Collection(collectionName)
}
