package handlers

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func HandleConversation(ctx context.Context, update tgbotapi.Update) {
	_ = ctx.Value("bot").(*tgbotapi.BotAPI)
}