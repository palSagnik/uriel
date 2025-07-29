package player

import "github.com/gin-gonic/gin"



func RegisterRoutes(router *gin.RouterGroup, handler *Handler, middleware gin.HandlerFunc) {
	players := router.Group("/players")
	{
		players.GET("/locations", middleware, handler.GetPlayerLocations)
	}
}