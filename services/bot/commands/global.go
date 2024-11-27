package commands

import (
	"Crawlzilla/database"
	cfg "Crawlzilla/logger"
	"Crawlzilla/services/bot/keyboards"
	"Crawlzilla/services/bot/menus"
	"Crawlzilla/services/super_admin"
	"Crawlzilla/services/users"
	"context"
	"strconv"

	"go.uber.org/zap"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CommandStart(ctx context.Context, update tgbotapi.Update) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	isAdmin := super_admin.IsSuperAdmin(update.Message.From.ID)

	if isAdmin {
		err := users.UpdateChatID(database.DB, update.SentFrom().ID, update.Message.Chat.ID)
		if err != nil {
			botLogger.Error(
				"Error while updating chatID for User",
				zap.Error(err),
				zap.String("user_id", strconv.FormatInt(update.SentFrom().ID, 10)),
			)
		}
	}
	_, err := users.LoginUser(database.DB, update.SentFrom().ID, update.Message.Chat.ID)

	if err != nil {
		botLogger.Error(
			"Error while calling login user",
			zap.Error(err),
			zap.String("user_id", strconv.FormatInt(update.SentFrom().ID, 10)),
		)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "خطا نگام دریافت اطلاعات کاربر")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "سلام به بات ما خوش اومدی!")

	msg.ReplyMarkup = keyboards.InlineKeyboard(menus.MainMenu, isAdmin)

	bot.Send(msg)

	cmdCfg := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "/start",
			Description: "شروع بات",
		},
		tgbotapi.BotCommand{
			Command:     "/menu",
			Description: "منو بات",
		},
	)

	bot.Send(cmdCfg)
}

func ShowMenu(ctx context.Context, update tgbotapi.Update) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)

	isAdmin := super_admin.IsSuperAdmin(update.Message.From.ID)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "منو خدمت شما")

	msg.ReplyMarkup = keyboards.InlineKeyboard(menus.MainMenu, isAdmin)

	bot.Send(msg)
}
