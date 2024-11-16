package commands

import (
	"Crawlzilla/services/bot/keyboards"
	"Crawlzilla/services/bot/menus"
	"Crawlzilla/services/super_admin"
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CommandStart(ctx context.Context, update tgbotapi.Update) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)

	isAdmin := super_admin.IsSuperAdmin(update.Message.From.ID)

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
