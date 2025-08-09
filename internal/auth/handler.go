package auth

import (
	"context"
	"fmt"
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

// RegisterUser is a Gin handler for User registration.
// It works between the HTTP request and the AuthService.
func (h *Handler) RegisterUser(c *gin.Context) {
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

	newUser, err := h.service.RegisterUserService(ctx, req)
	if err != nil {
		if err.Error() == "email already exists" || err.Error() == "username already exists" {
			c.JSON(http.StatusConflict, models.FailedResponse{
				Error: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.FailedResponse{
			Error: "Failed to register user",
		})
		return
	}

	c.JSON(http.StatusCreated, models.RegisterResponse{
		Message: "User registered succesfully",
		UserID:  newUser.ID.Hex(),
	})
}

func (h *Handler) LoginUser(c *gin.Context) {
	var req *models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	token, userId, err := h.service.LoginUserService(ctx, req.Username, req.Password)
	if err != nil {
		if err.Error() == "invalid username or password" {
			c.JSON(http.StatusUnauthorized, models.FailedResponse{
				Error: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.FailedResponse{
			Error: "Login failed due to internal server error",
		})
		return
	}

	// send it back
	authHeader := fmt.Sprintf("Bearer %v", token)
	c.Header("Authorization", authHeader)

	c.JSON(http.StatusOK, models.LoginResponse{
		Message: "User login successful",
		Token:   token,
		UserID:  userId,
	})
}
