package handlers

import (
	cfg "Crawlzilla/logger"
	"Crawlzilla/services/bot/conversations"
	"Crawlzilla/services/bot/keyboards"
	"Crawlzilla/services/bot/menus"
	"Crawlzilla/services/cache"
	"Crawlzilla/services/super_admin"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
	"strconv"
)

func HandleConversation(ctx context.Context, update tgbotapi.Update) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	state := ctx.Value("state").(*cache.UserState)

	userState, err := state.GetUserState(ctx, update.Message.Chat.ID)

	if err != nil {
		botLogger.Error(
			"Error while reading user state from store",
			zap.Error(err),
			zap.String("user_id", strconv.Itoa(update.Message.From.ID)),
			zap.String("user_name", update.Message.From.UserName),
		)
	}

	if userState == (cache.State{}) {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "متوجه منظورت نشدم!. از منو زیر استفاده کن.")
		isAdmin := super_admin.IsSuperAdmin(int64(update.Message.From.ID))
		msg.ReplyMarkup = keyboards.InlineKeyboard(menus.MainMenu, isAdmin)
		bot.Send(msg)
		return
	}

	switch userState.Conversation {
	case "add_ad":
		conversations.AddAdConversation(ctx, userState, update)
	}
}
