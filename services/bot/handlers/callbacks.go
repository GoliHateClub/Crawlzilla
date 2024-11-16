package handlers

import (
	"Crawlzilla/services/bot/conversations"
	"Crawlzilla/services/cache"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleCallbacks(ctx context.Context, update tgbotapi.Update) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)
	action := update.CallbackQuery.Data
	chatID := update.CallbackQuery.Message.Chat.ID

	switch action {
	case "/add_admin":
		bot.Send(tgbotapi.NewMessage(chatID, "Adding admin..."))
	case "/remove_admin":
		bot.Send(tgbotapi.NewMessage(chatID, "Removing admin..."))
	case "/get_admin":
		bot.Send(tgbotapi.NewMessage(chatID, "Removing admin..."))
	case "/get_all_users":
		bot.Send(tgbotapi.NewMessage(chatID, "Removing admin..."))
	case "/filters":
		bot.Send(tgbotapi.NewMessage(chatID, "Removing admin..."))
	case "/add_filter":
		conversations.AddFilterConversation(ctx, cache.CreateNewUserState("add_filter", update.CallbackQuery), update)
	case "/search":
		bot.Send(tgbotapi.NewMessage(chatID, "Removing admin..."))
	case "/remove_ad":
		bot.Send(tgbotapi.NewMessage(chatID, "Removing admin..."))
	case "/update_ad":
		bot.Send(tgbotapi.NewMessage(chatID, "Removing admin..."))
	case "/add_ad":
		conversations.AddAdConversation(ctx, cache.CreateNewUserState("add_ad", update.CallbackQuery), update)
	}

	// Acknowledge the callback to prevent the loading indicator
	bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, "Action received"))
}
