package ads

import (
	"Crawlzilla/database"
	cfg "Crawlzilla/logger"
	"Crawlzilla/services/cache"
	"Crawlzilla/services/search"
	"context"
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func GetMostFilteredAdsConversation(ctx context.Context, state cache.UserState, update tgbotapi.Update) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)
	userStates := ctx.Value("user_state").(*cache.UserCache)
	actionStates := ctx.Value("action_state").(*cache.ActionCache)
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	// Extract page number from the callback data (if provided)
	page := 1 // Default page
	if update.CallbackQuery != nil && update.CallbackQuery.Data != "" {
		action := update.CallbackQuery.Data
		if len(action) > len("/most_filtered_ads:") && action[:len("/most_filtered_ads:")] == "/most_filtered_ads:" {
			pageStr := action[len("/most_filtered_ads:"):]
			if p, err := strconv.Atoi(pageStr); err == nil {
				page = p
			}
		}
	}

	// Define page size
	pageSize := 5

	// Fetch ads using the service layer
	adData, err := search.GetMostFilteredAds(database.DB, page, pageSize)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(state.ChatId, "خطا در دریافت آگهی‌ها"))
		botLogger.Error("Error fetching most-filtered ads", zap.Error(err))
		return
	}

	// If no ads are found
	if len(adData.Data) == 0 {
		msg := tgbotapi.NewMessage(state.ChatId, "آگهی‌ای برای نمایش یافت نشد.")
		bot.Send(msg)
		return
	}

	// Format the response message
	response := fmt.Sprintf("📋 *آگهی‌های پربازدید (صفحه %d از %d):*\n\n", adData.Page, adData.Pages)
	var buttons [][]tgbotapi.InlineKeyboardButton

	for _, ad := range adData.Data {
		// Format each ad
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

	// Add pagination buttons if there are multiple pages
	if adData.Pages > 1 {
		response += "\n📄 *انتخاب صفحه:*"
		if adData.Page > 1 {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️ صفحه قبلی", fmt.Sprintf("/most_filtered_ads:%d", adData.Page-1)),
			))
		}
		if adData.Page < adData.Pages {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("➡️ صفحه بعدی", fmt.Sprintf("/most_filtered_ads:%d", adData.Page+1)),
			))
		}
	}

	// Send the response with inline buttons
	msg := tgbotapi.NewMessage(state.ChatId, response)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
	msg.ParseMode = "Markdown"
	bot.Send(msg)

	// Update user state
	err = userStates.SetUserCache(ctx, state.ChatId, cache.UserState{
		ChatId:       state.ChatId,
		UserId:       state.UserId,
		Stage:        "get_most_filtered_ads",
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

	// Update action state
	err = actionStates.SetUserState(ctx, state.ChatId, cache.ActionState{
		ChatId:       state.ChatId,
		UserId:       state.UserId,
		Conversation: state.Conversation,
		Action:       "most_filtered_ads",
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
