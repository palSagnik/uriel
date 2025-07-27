package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/palSagnik/uriel/internal/models"
)

type Handler struct {
	service *Service
}

func NewHandler(authService *Service) *Handler {
	return &Handler{service: authService}
}

// RegisterPlayer is a Gin handler for player registration.
// It works between the HTTP request and the AuthService.
func (h *Handler) RegisterPlayer(c *gin.Context) {
	var req *models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Password != req.Confirm {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password and confirm password fields do not match"})
		return
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	newPlayer, err := h.service.RegisterPlayerService(ctx, req)
	if err != nil {
		if err.Error() == "email already exists" || err.Error() == "username already exists" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register player"})
		return
	}

	c.JSON(http.StatusCreated, models.RegisterResponse{
		Message: "Player registered succesfully",
		PlayerID: newPlayer.ID.String(),
	})
}

func (h *Handler) LoginPlayer(c *gin.Context) {
	var req *models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	token, playerId, err := h.service.LoginPlayerService(ctx, req.Username, req.Password)
	if err != nil {
		if err.Error() == "invalid username or password" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed due to internal server error"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Message: "Player login successful",
		Token: token,
		PlayerID: playerId,
	})
}
