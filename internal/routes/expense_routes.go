package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/onunkwor/flypro-backend/internal/config"
	"github.com/onunkwor/flypro-backend/internal/handlers"
	"github.com/onunkwor/flypro-backend/internal/repository"
	"github.com/onunkwor/flypro-backend/internal/services"
)

func RegisterExpenseRoutes(router *gin.Engine) {
	repo := repository.NewExpenseRepository(config.DB)
	svc := services.NewExpenseService(repo)
	handler := handlers.NewExpenseHandler(svc)

	expenses := router.Group("/api/expenses")
	{
		expenses.POST("", handler.CreateExpense)
	}
}
