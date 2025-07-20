package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: implement only business logic here
func signupService(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func loginService(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

