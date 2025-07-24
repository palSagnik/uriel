package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	MongoDBURI string
	JWTSecret string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables.")
	}

	cfg := &Config{
		ServerPort: getEnv("SERVER_PORT", ":8080"),
		MongoDBURI: getEnv("MONGO_URI", "mongodb://localhost:27017/uriel?authSource=admin"),
		JWTSecret: getEnv("JWT_SECRET", "super_secret_jwt_key"),
	}

	if cfg.JWTSecret == "super_secret_default_key" {
		log.Println("WARNING: JWT_SECRET is using a default, insecure value. Please set it in environment variables.")
	}
	if strings.Contains(cfg.MongoDBURI, "localhost") && getEnv("MONGO_URI", "") == "" {
		log.Println("INFO: MONGO_URI is using a default 'localhost' value. Ensure MongoDB is running locally or via Docker Compose.")
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}