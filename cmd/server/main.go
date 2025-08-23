package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/onunkwor/flypro-backend/internal/config"
	"github.com/onunkwor/flypro-backend/internal/middleware"
	"github.com/onunkwor/flypro-backend/internal/routes"
	"github.com/onunkwor/flypro-backend/internal/validators"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
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

	// Create custom gin engine (not gin.Default(), so we control middlewares)
	router := gin.New()

	// zap logger for structured logs
	logger, _ := zap.NewProduction()
	router.Use(middleware.RequestLogger(logger))

	// recover from panics
	router.Use(gin.Recovery())

	// CORS
	router.Use(middleware.CORSMiddleware())

	// Rate limit: 5 requests/sec, burst of 10
	router.Use(middleware.RateLimiter(rate.Limit(5), 10))

	if err := router.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("failed to set trusted proxies: %v", err)
	}

	// Register routes
	routes.RegisterUserRoutes(router)
	routes.RegisterExpenseRoutes(router)
	routes.RegisterReportRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
