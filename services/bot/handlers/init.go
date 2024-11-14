package handlers

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type ConfigurationType struct {
	NewUpdateOffset int
	Timeout         int
}

var Configuration = &ConfigurationType{
	NewUpdateOffset: 0,
	Timeout:         60,
}

func Init(ctx context.Context) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)
	u := tgbotapi.NewUpdate(Configuration.NewUpdateOffset)
	u.Timeout = Configuration.Timeout

	updates, _ := bot.GetUpdatesChan(u)

	// Loop through each update.
	for update := range updates {
		if update.Message != nil && update.Message.IsCommand() {
			HandleCommands(ctx, update)
		}

		// Handle messages in conversation
		if update.Message != nil {
			HandleConversation(ctx, update)
		}

		// Handle callback queries (from InlineKeyboard)
		if update.CallbackQuery != nil {
			HandleCallbacks(ctx, update)
		}
	}
}
