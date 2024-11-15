package handlers

import (
	"Crawlzilla/services/bot/commands"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleCommands(ctx context.Context, update tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		commands.CommandStart(ctx, update)
	case "addAd":
		//startAddAdConversation(bot, update.Message)
	}
	return
}
