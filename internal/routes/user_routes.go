package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/onunkwor/flypro-backend/internal/config"
	"github.com/onunkwor/flypro-backend/internal/handlers"
	"github.com/onunkwor/flypro-backend/internal/repository"
	"github.com/onunkwor/flypro-backend/internal/services"
)

func RegisterUserRoutes(router *gin.Engine) {
	userRepo := repository.NewUserRepository(config.DB)
	userSvc := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userSvc)
	userGroup := router.Group("/api/users")
	{
		userGroup.POST("", userHandler.CreateUser)
		userGroup.GET("/:id", userHandler.GetUserByID)
	}
}
