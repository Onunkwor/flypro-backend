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
	config.ConnectRedis()
	log.Println("✅ Successfully connected to the database and Redis")
}

func main() {
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	router := gin.Default()

	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("failed to set trusted proxies: %v", err)
	}

	// Register routes
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
