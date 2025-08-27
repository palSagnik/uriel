package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Email     string             `bson:"email"`
	Username  string             `bson:"username"`
	Password  string             `bson:"password"`
	Role      string             `bson:"role"`
	AvatarUrl string             `bson:"avatar_url"`
	IsOnline  bool               `bson:"is_online"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}

type UpdateUserAvatarRequest struct {
	AvatarId string `json:"avatar_id"`
}
