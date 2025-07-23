package database

import (
	"context"

	"github.com/palSagnik/uriel/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	database *mongo.Database
}

func Connect() (*DB, error) {

	// Connect to the database.
    clientOptions := options.Client().ApplyURI(config.MONGO_DB_URI)
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        return nil, err
    }
    db := client.Database(config.DATABASE_NAME)

   return &DB{database: db}, nil
}

func (db *DB) Ping() error {
    if err := db.database.Client().Ping(context.Background(), nil); err != nil {
		return err
	}

	return nil
}
