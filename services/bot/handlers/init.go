package handlers

import (
	"Crawlzilla/services/cache"
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	bot, ok := ctx.Value("bot").(*tgbotapi.BotAPI)
	if !ok || bot == nil {
		log.Println("Bot instance is missing in context.")
		return
	}

	userState := cache.CreateUserCache(ctx)
	actionState := cache.CreateActionCache(ctx)
	ctx = context.WithValue(ctx, "user_state", userState)
	ctx = context.WithValue(ctx, "action_state", actionState)

	u := tgbotapi.NewUpdate(Configuration.NewUpdateOffset)
	u.Timeout = Configuration.Timeout

	updates := bot.GetUpdatesChan(u)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Received shutdown signal, stopping update processing...")
				return
			case update, ok := <-updates:
				if !ok {
					log.Println("Updates channel closed.")
					return
				}

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
	}()
}
