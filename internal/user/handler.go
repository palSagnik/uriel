package user

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

func NewHandler(userService *Service) *Handler {
	return &Handler{service: userService}
}

func (h *Handler) UpdateUserAvatar(c *gin.Context) {
	// retrieve the userId
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please login"})
		return
	}

	var req *models.UpdateUserAvatarRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	msg, err := h.service.UpdateUserAvatar(ctx, userID.(string), req.AvatarId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": msg,
	})
}

func (h *Handler) GetAllUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	users, err := h.service.GetUsers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func (h *Handler) GetAllAvatars(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	avatars, err := h.service.GetAvatars(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.GetAvatarsResponse{
		Avatars: avatars,
	})
}
