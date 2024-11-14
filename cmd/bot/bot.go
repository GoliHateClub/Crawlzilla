package bot

import (
	cfg "Crawlzilla/logger"
	tgBot "Crawlzilla/services/bot"
	"Crawlzilla/services/bot/handlers"
	"context"
	"fmt"
)

func StartBot(ctx context.Context) {
	// Listen for the shutdown signal from the context

	for {
		select {
		case <-ctx.Done(): // triggered if the server's context is canceled
			fmt.Println("Bot received shutdown signal from context, stopping...")
			return
		default:
			// Load loggers
			configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
			botLogger, _ := configLogger("bot")

			botLogger.Info("Bot started successfully")

			bot := tgBot.Init()

			ctx = context.WithValue(ctx, "bot", bot)

			handlers.Init(ctx)
		}
	}
}
