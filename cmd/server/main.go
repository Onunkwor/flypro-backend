package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/onunkwor/flypro-backend/internal/config"
	"github.com/onunkwor/flypro-backend/internal/routes"
	"github.com/onunkwor/flypro-backend/internal/validators"
)

func init() {
	config.LoadEnvVarialbles()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validators.CurrencyValidator)
	}
	if err := config.ConnectDatabase(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
}
func main() {
	router := gin.Default()
	routes.RegisterUserRoutes(router)
	routes.RegisterExpenseRoutes(router)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
