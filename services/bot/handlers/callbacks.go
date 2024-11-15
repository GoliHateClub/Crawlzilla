package handlers

import (
	"Crawlzilla/services/bot/conversations"
	"Crawlzilla/services/cache"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
	case "/add_ad":
		conversations.AddAdConversation(ctx, cache.CreateNewState("add_ad", update.CallbackQuery), update)
	case "/remove_ad":
		bot.Send(tgbotapi.NewMessage(chatID, "Removing admin..."))
	case "/update_ad":
		bot.Send(tgbotapi.NewMessage(chatID, "Removing admin..."))
	case "/see_all_ads":
		bot.Send(tgbotapi.NewMessage(chatID, "Removing admin..."))
		// Add other cases as needed
	}

	// Acknowledge the callback to prevent the loading indicator
	bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Action received"))
}
