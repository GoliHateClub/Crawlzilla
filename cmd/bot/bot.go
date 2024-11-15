package bot

import (
	cfg "Crawlzilla/logger"
	tgBot "Crawlzilla/services/bot"
	"Crawlzilla/services/bot/handlers"
	"Crawlzilla/services/cache"
	"context"
	"log"
)

func StartBot(ctx context.Context) {
	// Load loggers
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	botLogger.Info("Bot started successfully")

	redis := cache.InitRedis(ctx)

	bot := tgBot.Init()

	ctx = context.WithValue(ctx, "bot", bot)
	ctx = context.WithValue(ctx, "redis", redis)

	handlers.Init(ctx)

	// Wait for shutdown signal
	<-ctx.Done()
	botLogger.Info("Bot shutdown initiated.")
	log.Println("Bot shutdown initiated.")
	return
}
