package config

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client

func ConnectRedis() error {
	Redis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return Redis.Ping(context.Background()).Err()
}
