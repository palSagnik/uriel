package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	UserID  string `json:"player_id"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	UserID  string `json:"player_id"`
}

type Claims struct {
	UserID   string
	Username string
	Role     string
	jwt.RegisteredClaims
}
