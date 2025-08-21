package handlers

import "github.com/gin-gonic/gin"

func RegisterUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/api/users")
	{
		userGroup.POST("", CreateUser)
	}
}
