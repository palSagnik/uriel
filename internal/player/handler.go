package player

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(playerService *Service) *Handler {
	return &Handler{service: playerService}
}

func (h *Handler) GetPlayerLocations(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome",
	})
}