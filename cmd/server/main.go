package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/onunkwor/flypro-backend/internal/config"
	"github.com/onunkwor/flypro-backend/internal/handlers"
)

func init() {
	config.LoadEnvVarialbles()
	if err := config.ConnectDatabase(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
}
func main() {
	router := gin.Default()
	handlers.RegisterUserRoutes(router)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 Server running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
