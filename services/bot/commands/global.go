package commands

import (
	cfg "Crawlzilla/logger"
	"Crawlzilla/services/bot/keyboards"
	"Crawlzilla/services/bot/menus"
	"Crawlzilla/services/super_admin"
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CommandStart(ctx context.Context, update tgbotapi.Update) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	isAdmin := super_admin.IsSuperAdmin(int64(update.Message.From.ID))

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "سلام به بات ما خوش اومدی!")

	msg.ReplyMarkup = keyboards.InlineKeyboard(menus.MainMenu, isAdmin)

	_, err := bot.Send(msg)

	if err != nil {
		println(err.Error())
		botLogger.Error("Error while sending welcome message.",
			zap.Error(err), zap.String("user_id", strconv.Itoa(update.Message.From.ID)),
		)
	}
}
