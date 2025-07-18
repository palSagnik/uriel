package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: struct will be extended further, this is enough for authentication
type Player struct {
	ID       primitive.ObjectID `bson:"_id"`
	Email    string             `bson:"email"`
	Username string             `bson:"username"`
	Password string             `bson:"password"`
	Role     string             `bson:"role"`
}
