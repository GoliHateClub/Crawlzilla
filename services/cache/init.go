package cache

import (
	cfg "Crawlzilla/logger"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"log"
	"os"
)

func InitRedis(ctx context.Context) *redis.Client {
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	RedisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT"),
		),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	// Test connection
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		botLogger.Error("Failed to connect to Redis", zap.Error(err))
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return RedisClient
}
