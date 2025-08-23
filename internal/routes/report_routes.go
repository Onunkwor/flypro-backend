package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/onunkwor/flypro-backend/internal/config"
	"github.com/onunkwor/flypro-backend/internal/handlers"
	"github.com/onunkwor/flypro-backend/internal/repository"
	"github.com/onunkwor/flypro-backend/internal/services"
)

func RegisterReportRoutes(r *gin.Engine) {

	reportRepo := repository.NewReportRepository(config.DB)
	expenseRepo := repository.NewExpenseRepository(config.DB)
	currencySvc := services.NewCurrencyService(config.Redis, config.EnvCurrencyAPIKey())
	reportService := services.NewReportService(reportRepo, expenseRepo, currencySvc, config.Redis)
	reportHandler := handlers.NewReportHandler(reportService)

	api := r.Group("/api/reports")
	{
		api.POST("", reportHandler.CreateReport)
		api.POST("/:id/expenses", reportHandler.AddExpense)
		api.GET("", reportHandler.ListReports)
		api.PUT("/:id/submit", reportHandler.SubmitReport)
	}
}
