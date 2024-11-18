package ads

import (
	"Crawlzilla/database"
	cfg "Crawlzilla/logger"
	"Crawlzilla/services/ads"
	"Crawlzilla/services/cache"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"log"
	"strconv"
)

func GetAllAdConversation(ctx context.Context, state cache.UserState, update tgbotapi.Update) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)
	userStates := ctx.Value("user_state").(*cache.UserCache)
	actionStates := ctx.Value("action_state").(*cache.ActionCache)
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	// Extract page number from the action (if provided)
	page := 1 // Default page

	// action := update.CallbackQuery.Data
	if update.CallbackQuery != nil && update.CallbackQuery.Data != "" {
		action := update.CallbackQuery.Data
		if len(action) > len("/see_all_ads:") && action[:len("/see_all_ads:")] == "/see_all_ads:" {
			pageStr := action[len("/see_all_ads:"):]
			if p, err := strconv.Atoi(pageStr); err == nil {
				page = p
			}
		}
	}

	// Define page size
	pageSize := 2

	// Fetch ads using the service layer
	adData, err := ads.GetAllAds(database.DB, page, pageSize)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(state.ChatId, "خطا در دریافت آگهی‌ها"))
		return
	}

	// If no ads are found
	if len(adData.Data) == 0 {
		msg := tgbotapi.NewMessage(state.ChatId, "آگهی‌ای برای نمایش یافت نشد.")
		bot.Send(msg)
		return
	}

	response := fmt.Sprintf("📋 *آگهی‌های موجود (صفحه %d از %d):*\n\n", adData.Page, adData.Pages)
	var buttons [][]tgbotapi.InlineKeyboardButton

	for _, ad := range adData.Data {
		// Format each ad with emojis and details
		response += fmt.Sprintf(
			"🏷️ *عنوان:* %s\n"+
				"🆔 *شناسه:* `%s`\n"+
				"🔍 [مشاهده جزئیات](%s)\n\n",
			ad.Title, ad.ID, fmt.Sprintf("/view_ad:%s", ad.ID),
		)

		// Add an inline button for each ad to view its details
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("🔍 مشاهده: %s", ad.Title), fmt.Sprintf("/view_ad:%s", ad.ID)),
		))
	}

	if adData.Pages > 1 {
		response += "\n📄 *انتخاب صفحه:*"
		if adData.Page > 1 {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️ صفحه قبلی", fmt.Sprintf("/see_all_ads:%d", adData.Page-1)),
			))
		}
		if adData.Page < adData.Pages {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("➡️ صفحه بعدی", fmt.Sprintf("/see_all_ads:%d", adData.Page+1)),
			))
		}
	}

	// Send the message with inline buttons
	msg := tgbotapi.NewMessage(state.ChatId, response)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
	msg.ParseMode = "Markdown" // Enables better formatting for bold text, links, etc.
	bot.Send(msg)

	// Update user state and action state
	err = userStates.SetUserCache(ctx, state.ChatId, cache.UserState{
		ChatId:       state.ChatId,
		UserId:       state.UserId,
		Stage:        "get_ads",
		Conversation: state.Conversation,
	})
	if err != nil {
		botLogger.Error(
			"Error updating user state",
			zap.Error(err),
			zap.String("user_id", strconv.Itoa(int(state.UserId))),
			zap.String("chat_id", strconv.Itoa(int(state.ChatId))),
		)
		log.Printf("Error updating user state: %v", err)
	}

	err = actionStates.SetUserState(ctx, state.ChatId, cache.ActionState{
		ChatId:       state.ChatId,
		UserId:       state.UserId,
		Conversation: state.Conversation,
		Action:       "view_ads",
		ActionData: map[string]interface{}{
			"page": page,
		},
	})
	if err != nil {
		cache.HandleActionStateError(botLogger, state, err)
		bot.Send(tgbotapi.NewMessage(state.ChatId, "خطایی رخ داد!"))
		return
	}
}
