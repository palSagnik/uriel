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

func (h *Handler) GetUserLocations(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome",
	})
}

func (h *Handler) UpdateMetadata(c *gin.Context) {
	var req *models.UpdateMetadataRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	if msg, err := h.service.UpdateMetadata(ctx, req.AvatarUrl); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "avatar_url updated successfully",
	})
}
