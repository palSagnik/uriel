package database

import (
	"context"

	"github.com/palSagnik/uriel/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient(uri string) (*mongo.Client, error) {

    clientOptions := options.Client().ApplyURI(uri)
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        return nil, err
    }
    db := client.Database(config.DATABASE_NAME)

   return db.Client(), nil
}

// func (db *DB) Ping() error {
//     if err := db.database.Client().Ping(context.Background(), nil); err != nil {
// 		return err
// 	}

// 	return nil
// }
