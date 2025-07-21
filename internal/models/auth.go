package models

type RegisterRequest struct {
	Email           string
	Username        string
	Password        string
	ConfirmPassword string
}

type Credentials struct {
	Email    string
	Password string
}
