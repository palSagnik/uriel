package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/palSagnik/uriel/internal/auth"
	"github.com/palSagnik/uriel/internal/config"
	"github.com/palSagnik/uriel/internal/database"
	"github.com/palSagnik/uriel/internal/user"
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
	userRepo := database.NewuserRepository(mongodb)

	// --- Initialise Services ---
	authService := auth.NewService(authRepo, []byte(cfg.JWTSecret))
	userService := user.NewService(userRepo)

	// --- Initialise Handlers ---
	authHandler := auth.NewHandler(authService)
	userHandler := user.NewHandler(userService)

	// --- Initialise Middleware ---
	authMiddleware := authService.AuthMiddleware()

	v1 := router.Group("/api/v1")
	{
		auth.RegisterRoutes(v1, authHandler)
		user.RegisterRoutes(v1, userHandler, authMiddleware)
	}

	// --- Running the server ---
	router.Run(cfg.ServerPort)
}
