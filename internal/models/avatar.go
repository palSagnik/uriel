package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Avatar struct {
	ID        primitive.ObjectID `bson:"_id"`
	AvatarUrl string             `bson:"avatar_url"`
	Name      string             `bson:"name"`
}
