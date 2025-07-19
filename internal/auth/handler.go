package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func SetupAuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")

	auth.POST("/signup", signup)
	auth.POST("/login", login)
}

// TODO: implement cmplete auth handlers
func signup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

