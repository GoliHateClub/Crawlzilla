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

	switch {
	case len(action) > len("/view_ad:") && action[:len("/view_ad:")] == "/view_ad:":
		conversations.GetAdDetailsConversation(ctx, update)
	case action == "/add_admin":
		bot.Send(tgbotapi.NewMessage(chatID, "Adding admin..."))
	case action == "/remove_admin":
		bot.Send(tgbotapi.NewMessage(chatID, "Removing admin..."))
	case action == "/add_ad":
		conversations.AddAdConversation(ctx, cache.CreateNewUserState("add_ad", update.CallbackQuery), update)
	case action == "/remove_ad":
		bot.Send(tgbotapi.NewMessage(chatID, "Removing admin..."))
	case action == "/update_ad":
		bot.Send(tgbotapi.NewMessage(chatID, "Removing admin..."))
	case len(action) >= len("/see_all_ads") && action[:len("/see_all_ads")] == "/see_all_ads":
		conversations.GetAllAdConversation(ctx, cache.CreateNewUserState("see_all_ads", update.CallbackQuery), update)
	case "/get_admin":
		bot.Send(tgbotapi.NewMessage(chatID, "Removing admin..."))
	case "/get_all_users":
		bot.Send(tgbotapi.NewMessage(chatID, "Removing admin..."))
	case "/filters":
		bot.Send(tgbotapi.NewMessage(chatID, "Removing admin..."))
	case "/add_filter":
		conversations.AddFilterConversation(ctx, cache.CreateNewUserState("add_filter", update.CallbackQuery), update)
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
