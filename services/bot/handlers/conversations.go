package handlers

import (
	cfg "Crawlzilla/logger"
	"Crawlzilla/services/bot/conversations/ads"
	"Crawlzilla/services/bot/conversations/configs"
	"Crawlzilla/services/bot/conversations/filters"
	"Crawlzilla/services/bot/keyboards"
	"Crawlzilla/services/bot/menus"
	"Crawlzilla/services/cache"
	"Crawlzilla/services/super_admin"
	"context"
	"strconv"

	"go.uber.org/zap"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleConversation(ctx context.Context, update tgbotapi.Update) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	state := ctx.Value("user_state").(*cache.UserCache)

	userState, err := state.GetUserCache(ctx, update.Message.Chat.ID)

	if err != nil {
		botLogger.Error(
			"Error while reading user state from store",
			zap.Error(err),
			zap.String("user_id", strconv.Itoa(int(update.Message.From.ID))),
			zap.String("user_name", update.Message.From.UserName),
		)
	}

	if userState == (cache.UserState{}) {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "متوجه منظورت نشدم!. از منو زیر استفاده کن.")
		isAdmin := super_admin.IsSuperAdmin(update.Message.From.ID)
		msg.ReplyMarkup = keyboards.InlineKeyboard(menus.MainMenu, isAdmin)
		bot.Send(msg)
		return
	}

	switch userState.Conversation {
	case "add_ad":
		ads.AddAdConversation(ctx, userState, update)
	case "see_all_ads":
		ads.GetAllAdConversation(ctx, userState, update)
	case "add_filter":
		filters.AddFilterConversation(ctx, userState, update)
	case "see_all_filters":
		filters.GetAllFilterConversation(ctx, userState, update)
	case "view_filter_details":
		filters.ViewFilterDetailsConversation(ctx, userState, update)
	case "config_crawler":
		configs.ConfigCrawlerConversation(ctx, userState, update)
	}
}
