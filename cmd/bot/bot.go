package bot

import (
	cfg "Crawlzilla/logger"
	"Crawlzilla/services/bot"
	"context"
	"fmt"
)

func StartBot(ctx context.Context) {
	// Load loggers
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	botLogger.Info("Bot started successfully")

	bot.Init()

	// Listen for the shutdown signal from the context
	select {
	case <-ctx.Done(): // triggered if the server's context is canceled
		fmt.Println("Bot received shutdown signal from context, stopping...")
		return
	}
}
