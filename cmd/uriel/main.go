package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/palSagnik/uriel/internal/auth"
	"github.com/palSagnik/uriel/internal/config"
	"github.com/palSagnik/uriel/internal/database"
	"github.com/palSagnik/uriel/internal/player"
)


func main() {

	// Loading config
	cfg := config.LoadConfig()

	mongodb, err := database.NewMongoClient(cfg.MongoDBURI)
	if err != nil {
		log.Fatalf("Failed to connect to mongo: %v", err)
	}

	router := gin.Default()
	router.Use(gin.Logger(), gin.Recovery())

	// --- Initialise Repositories ---
	authRepo := database.NewAuthRepository(mongodb)
	playerRepo := database.NewPlayerRepository(mongodb)

	// --- Initialise Services ---
	authService := auth.NewService(authRepo, []byte(cfg.JWTSecret))
	playerService := player.NewService(playerRepo)

	// --- Initialise Handlers ---
	authHandler := auth.NewHandler(authService)
	playerHandler := player.NewHandler(playerService)

	// --- Initialise Middleware ---
	authMiddleware := authService.AuthMiddleware()

	v1 := router.Group("/api/v1")
	{
		auth.RegisterRoutes(v1, authHandler)
		player.RegisterRoutes(v1, playerHandler, authMiddleware)
	}

	// --- Running the server ---
	router.Run(cfg.ServerPort)
}