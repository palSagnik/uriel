package config

import (
	"os"
	_ "github.com/joho/godotenv"
)

var (
	MONGO_DB_URI = os.Getenv("MONGO_DB_URI")
	DATABASE_NAME = os.Getenv("DATABASE_NAME")
)

// COLLECTIONS
const PLAYER_COLLECTION = "player"
const MAP_COLLECTION = "map"

// ROLES
const PLAYER = "player"
const ADMIN = "admin"
const GUEST = "guest"