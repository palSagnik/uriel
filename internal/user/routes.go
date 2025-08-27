package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, handler *Handler, middleware gin.HandlerFunc) {
	users := router.Group("/users")
	{
		users.POST("/avatar", middleware, handler.UpdateUserAvatar)
		users.GET("/avatar", middleware, handler.GetAllAvatars)
		users.GET("/user", middleware, handler.GetAllUsers)
	}
}
