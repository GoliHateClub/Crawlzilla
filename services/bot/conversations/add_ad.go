package conversations

import (
	"Crawlzilla/services/cache"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AddFilterConversation(ctx context.Context, state cache.UserState, update tgbotapi.Update) {
	//bot := ctx.Value("bot").(*tgbotapi.BotAPI)
	//userStates := ctx.Value("user_state").(*cache.UserCache)
	//actionStates := ctx.Value("action_state").(*cache.ActionCache)
	//configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	//botLogger, _ := configLogger("bot")

	switch state.Stage {
	case "init":

	}
}
