package filters

import (
	"Crawlzilla/database"
	cfg "Crawlzilla/logger"
	"Crawlzilla/services/cache"
	"Crawlzilla/services/search"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func ApplyFilterConversation(ctx context.Context, state cache.UserState, update tgbotapi.Update) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)
	userStates := ctx.Value("user_state").(*cache.UserCache)
	actionStates := ctx.Value("action_state").(*cache.ActionCache)
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	// Extract the filter ID from the callback data
	action := update.CallbackQuery.Data
	parts := strings.Split(action[len("/apply_filter:"):], ":")
	filterID := parts[0]

	// Extract page number if available
	page := 1 // Default page
	if len(parts) > 1 {
		if p, err := strconv.Atoi(parts[1]); err == nil {
			page = p
		}
	}
	pageSize := 5

	// Fetch filtered ads using the provided service
	adsData, err := search.GetFilteredAds(database.DB, filterID, page, pageSize)
	log.Println(filterID, page, pageSize)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(state.ChatId, "خطا در دریافت آگهی‌ها با استفاده از فیلتر"))
		return
	}

	// If no ads match the filter
	if len(adsData.Data) == 0 {
		bot.Send(tgbotapi.NewMessage(state.ChatId, "هیچ آگهی‌ای مطابق با فیلتر یافت نشد."))
		return
	}

	// Generate response message
	response := fmt.Sprintf("📋 *نتایج فیلتر (صفحه %d از %d):*\n\n", adsData.Page, adsData.Pages)
	var buttons [][]tgbotapi.InlineKeyboardButton

	for _, ad := range adsData.Data {
		response += fmt.Sprintf(
			"🏷️ *عنوان:* %s\n"+
				"🏙️ *شهر:* %s\n"+
				"📍 *محله:* %s\n"+
				"💰 *قیمت:* %d\n"+
				"🚪 *اتاق‌ها:* %d\n"+
				"📐 *متراژ:* %d\n"+
				"🔍 [مشاهده جزئیات](%s)\n\n",
			ad.Title,
			ad.City,
			ad.Neighborhood,
			ad.Price,
			ad.Room,
			ad.Area,
			fmt.Sprintf("/view_ad:%s", ad.ID),
		)

		// Add a button to view the ad details
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("🔍 مشاهده: %s", ad.Title), fmt.Sprintf("/view_ad:%s", ad.ID)),
		))
	}

	// Add pagination buttons
	if adsData.Pages > 1 {
		if adsData.Page > 1 {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️ صفحه قبلی", fmt.Sprintf("/apply_filter:%s:%d", filterID, adsData.Page-1)),
			))
		}
		if adsData.Page < adsData.Pages {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("➡️ صفحه بعدی", fmt.Sprintf("/apply_filter:%s:%d", filterID, adsData.Page+1)),
			))
		}
	}

	// Send the response
	msg := tgbotapi.NewMessage(state.ChatId, response)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
	msg.ParseMode = "Markdown"
	bot.Send(msg)

	// Update the user's state
	err = userStates.SetUserCache(ctx, state.ChatId, cache.UserState{
		ChatId:       state.ChatId,
		UserId:       state.UserId,
		Stage:        "apply_filter",
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

	// Save the action state for pagination
	err = actionStates.SetUserState(ctx, state.ChatId, cache.ActionState{
		ChatId:       state.ChatId,
		UserId:       state.UserId,
		Conversation: state.Conversation,
		Action:       "apply_filter",
		ActionData: map[string]interface{}{
			"filter_id": filterID,
			"page":      page,
		},
	})
	if err != nil {
		cache.HandleActionStateError(botLogger, state, err)
		bot.Send(tgbotapi.NewMessage(state.ChatId, "خطایی رخ داد!"))
		return
	}
}
