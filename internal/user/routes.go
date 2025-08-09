package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler, middleware gin.HandlerFunc) {
	users := router.Group("/users")
	{
    users.GET("/locations", middleware, handler.GetUserLocations)
		users.POST("/metadata", middleware, handler.UpdateMetadata)
	}
}
