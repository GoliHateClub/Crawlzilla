package bot

import (
	cfg "Crawlzilla/logger"
	"Crawlzilla/services/bot/handlers"
	"Crawlzilla/services/cache"
	"context"
)

func StartBot(ctx context.Context) {
	// Load loggers
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	botLogger.Info("Bot started successfully")

	redis := cache.InitRedis(ctx)

	ctx = context.WithValue(ctx, "redis", redis)

	handlers.Init(ctx)

	// Wait for shutdown signal
	<-ctx.Done()
	botLogger.Info("Bot shutdown initiated.")
	return
}
