package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/onunkwor/flypro-backend/internal/config"
	"github.com/onunkwor/flypro-backend/internal/handlers"
	"github.com/onunkwor/flypro-backend/internal/services"
)

func RegisterCurrencyRoutes(router *gin.Engine) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.EnvRedisAddr(),
		Password: config.EnvRedisPassword(),
		DB:       0,
	})
	svc := services.NewCurrencyService(rdb, config.EnvCurrencyAPIKey())
	handler := handlers.NewCurrencyHandler(svc)

	router.POST("/api/expenses/convert", handler.ConvertCurrency)
}
