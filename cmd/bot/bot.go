package bot

import (
	cfg "Crawlzilla/logger"
	tgBot "Crawlzilla/services/bot"
	"Crawlzilla/services/bot/handlers"
	"context"
	"log"
)

func StartBot(ctx context.Context) {
	// Load loggers
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	botLogger.Info("Bot started successfully")

	bot := tgBot.Init()

	ctx = context.WithValue(ctx, "bot", bot)

	handlers.Init(ctx)

	// Wait for shutdown signal
	<-ctx.Done()
	botLogger.Info("Bot shutdown initiated.")
	log.Println("Bot shutdown initiated.")
	return
}
