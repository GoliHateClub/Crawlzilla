package bot

import (
	cfg "Crawlzilla/logger"
	"context"
	"fmt"
)

func StartBot(ctx context.Context) {
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	botLogger.Info("Bot started successfully")

	// Listen for the shutdown signal from the context
	select {
	case <-ctx.Done(): // triggered if the server's context is canceled
		fmt.Println("Bot received shutdown signal from context, stopping...")
		return
	}
}
