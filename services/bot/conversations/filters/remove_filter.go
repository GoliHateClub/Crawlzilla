package filters

import (
	"Crawlzilla/database"
	cfg "Crawlzilla/logger"
	filterService "Crawlzilla/services/filters"
	usersService "Crawlzilla/services/users"
	"context"
	"errors"

	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func DeleteFilterConversation(ctx context.Context, update tgbotapi.Update) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	// Extract filter ID from the callback data
	callbackData := update.CallbackQuery.Data
	if len(callbackData) <= len("/delete_filter:") {
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "خطای غیرمنتظره‌ای رخ داد!"))
		return
	}
	filterID := callbackData[len("/delete_filter:"):]

	// Extract Telegram user ID
	telegramUserID := update.CallbackQuery.From.ID
	userID, err := usersService.GetUserIDByTelegramID(database.DB, strconv.FormatInt(telegramUserID, 10))
	if err != nil {
		botLogger.Error("Error retrieving user ID by Telegram ID", zap.Error(err))
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "خطایی رخ داد! لطفاً دوباره تلاش کنید."))
		return
	}

	// Attempt to delete the filter
	err = filterService.RemoveFilter(database.DB, userID, filterID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "فیلتر یافت نشد!"))
		} else if errors.Is(err, errors.New("unauthorized to delete this filter")) {
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "شما مجاز به حذف این فیلتر نیستید!"))
		} else {
			botLogger.Error("Error deleting filter", zap.Error(err))
			bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "خطایی هنگام حذف فیلتر رخ داد!"))
		}
		return
	}

	// Send confirmation message
	bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "فیلتر با موفقیت حذف شد!"))
}
