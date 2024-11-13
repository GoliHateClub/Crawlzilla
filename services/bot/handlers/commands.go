package handlers

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func HandleCommands(ctx context.Context, update tgbotapi.Update) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)

	switch update.Message.Command() {
	case "start":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome! Use /addAd to add an ad.")
		bot.Send(msg)
	case "addAd":
		//startAddAdConversation(bot, update.Message)
	}
	return
}
