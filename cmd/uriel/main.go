package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/palSagnik/uriel/internal/auth"
	"github.com/palSagnik/uriel/internal/config"
	"github.com/palSagnik/uriel/internal/database"
)


func main() {

	// Loading config
	cfg := config.LoadConfig()

	mongodb, err := database.NewMongoClient(cfg.MongoDBURI)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	// --- Initialise Repositories ---
	authRepo := database.NewMongoAuthRepository(mongodb)

	// --- Initialise Services ---
	authService := auth.NewService(authRepo)

	// --- Initialise Handlers ---
	authHandler := auth.NewHandler(authService)

	v1 := router.Group("/api/v1")
	{
		auth.RegisterRoutes(v1, authHandler)
	}
}