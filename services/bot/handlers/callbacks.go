package handlers

import (
	"Crawlzilla/services/bot/conversations/ads"
	"Crawlzilla/services/bot/conversations/configs"
	"Crawlzilla/services/bot/conversations/filters"
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
		ads.GetAdDetailsConversation(ctx, update)
	case action == "/add_admin":
		bot.Send(tgbotapi.NewMessage(chatID, "Adding admin..."))
	case action == "/remove_admin":
		bot.Send(tgbotapi.NewMessage(chatID, "Removing admin..."))
	case action == "/add_ad":
		ads.AddAdConversation(ctx, cache.CreateNewUserState("add_ad", update.CallbackQuery), update)
	case action == "/remove_ad":
		bot.Send(tgbotapi.NewMessage(chatID, "Removing admin..."))
	case action == "/update_ad":
		bot.Send(tgbotapi.NewMessage(chatID, "Removing admin..."))
	case len(action) >= len("/see_all_ads") && action[:len("/see_all_ads")] == "/see_all_ads":
		ads.GetAllAdConversation(ctx, cache.CreateNewUserState("see_all_ads", update.CallbackQuery), update)
	case action == "/get_admin":
		bot.Send(tgbotapi.NewMessage(chatID, "Removing admin..."))
	case action == "/get_all_users":
		bot.Send(tgbotapi.NewMessage(chatID, "Removing admin..."))
	case action == "/filters":
		bot.Send(tgbotapi.NewMessage(chatID, "Removing admin..."))
	case action == "/add_filter":
		filters.AddFilterConversation(ctx, cache.CreateNewUserState("add_filter", update.CallbackQuery), update)
	case len(action) >= len("/see_all_filters") && action[:len("/see_all_filters")] == "/see_all_filters":
		filters.GetAllFilterConversation(ctx, cache.CreateNewUserState("see_all_filters", update.CallbackQuery), update)
	case len(action) > len("/view_filter:") && action[:len("/view_filter:")] == "/view_filter:":
		filters.ViewFilterDetailsConversation(ctx, cache.CreateNewUserState("view_filter_details", update.CallbackQuery), update)
	case len(action) > len("/delete_filter:") && action[:len("/delete_filter:")] == "/delete_filter:":
		filters.DeleteFilterConversation(ctx, update)
	case len(action) > len("/apply_filter:") && action[:len("/apply_filter:")] == "/apply_filter:":
		filters.ApplyFilterConversation(ctx, cache.CreateNewUserState("apply_filter", update.CallbackQuery), update)
	case action == "/config":
		configs.ConfigCrawlerConversation(ctx, cache.CreateNewUserState("config_crawler", update.CallbackQuery), update)
	case len(action) > len("/export_filter:") && action[:len("/export_filter:")] == "/export_filter:":
		filters.ExportFilteredResultsConversation(ctx, cache.CreateNewUserState("export_filter", update.CallbackQuery), update)
	case action == "/remove_all_filters":
		filters.RemoveAllFiltersConversation(ctx, cache.CreateNewUserState("remove_all_filters", update.CallbackQuery), update)
	case len(action) >= len("/most_filtered_ads") && action[:len("/most_filtered_ads")] == "/most_filtered_ads":
		ads.GetMostFilteredAdsConversation(ctx, cache.CreateNewUserState("most_filtered_ads", update.CallbackQuery), update)
	case action == "/start_crawler":
		configs.StartCrawlerConversation(ctx, update)
	}

	// Acknowledge the callback to prevent the loading indicator
	bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, "Action received"))
}
