package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/palSagnik/uriel/internal/auth"
	"github.com/palSagnik/uriel/internal/config"
	"github.com/palSagnik/uriel/internal/database"
)


func main() {
	mongoClient, err := database.NewMongoClient(config.MONGO_DB_URI)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	// --- Initialise Repositories ---
	authRepo := database.NewMongoAuthRepository(mongoClient)

	// --- Initialise Services ---
	authService := auth.NewService(authRepo)

	// --- Initialise Handlers ---
	authHandler := auth.NewHandler(authService)

	v1 := router.Group("/api/v1")
	{
		auth.RegisterRoutes(v1, authHandler)
	}
}