package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func SetupAuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")

	auth.POST("/register", register)
	auth.POST("/login", login)
}

// TODO: implement cmplete auth handlers
// handle all gin contexts here and all http management 
// then pass the data to the service layer
func register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}