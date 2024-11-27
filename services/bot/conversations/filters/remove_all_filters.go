package filters

import (
	"Crawlzilla/database"
	cfg "Crawlzilla/logger"
	"Crawlzilla/services/cache"
	"Crawlzilla/services/filters"
	"Crawlzilla/services/users"
	"context"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func RemoveAllFiltersConversation(ctx context.Context, state cache.UserState, update tgbotapi.Update) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	// Retrieve the user ID from the state
	userID := state.UserId

	// Convert Telegram ID to User ID from the database
	dbUserID, err := users.GetUserIDByTelegramID(database.DB, strconv.Itoa(int(userID)))
	if err != nil {
		bot.Send(tgbotapi.NewMessage(state.ChatId, "خطا در شناسایی کاربر! لطفاً دوباره تلاش کنید."))
		botLogger.Error("Error retrieving user ID", zap.Error(err))
		return
	}

	// Call the service to remove all filters for the user
	err = filters.RemoveAllFilters(database.DB, dbUserID)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(state.ChatId, "خطا در حذف فیلترها! لطفاً دوباره تلاش کنید."))
		botLogger.Error("Error removing all filters", zap.Error(err))
		return
	}

	// Notify the user of successful removal
	bot.Send(tgbotapi.NewMessage(state.ChatId, "✅ تمامی فیلترها با موفقیت حذف شدند."))

	// Log the action
	botLogger.Info("All filters removed successfully for user", zap.Int64("telegram_id", userID))
}
