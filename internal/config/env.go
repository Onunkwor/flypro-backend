package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVarialbles() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func EnvRedisAddr() string {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	return addr
}

func EnvRedisPassword() string {
	return os.Getenv("REDIS_PASSWORD")
}

func EnvCurrencyAPIKey() string {
	apiKey := os.Getenv("CURRENCY_API_KEY")
	if apiKey == "" {
		log.Println("⚠️ CURRENCY_API_KEY not set, defaulting to free API without key")
	}
	return apiKey
}
