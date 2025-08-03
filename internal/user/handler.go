package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(userService *Service) *Handler {
	return &Handler{service: userService}
}

func (h *Handler) GetuserLocations(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome",
	})
}
