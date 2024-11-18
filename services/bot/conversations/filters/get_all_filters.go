package filters

import (
	"Crawlzilla/database"
	cfg "Crawlzilla/logger"
	"Crawlzilla/models"
	"Crawlzilla/services/cache"
	filterService "Crawlzilla/services/filters"
	"Crawlzilla/services/users"
	"context"
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func GetAllFilterConversation(ctx context.Context, state cache.UserState, update tgbotapi.Update) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)
	userStates := ctx.Value("user_state").(*cache.UserCache)
	actionStates := ctx.Value("action_state").(*cache.ActionCache)
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	// Retrieve database user ID from Telegram ID
	userID, err := users.GetUserIDByTelegramID(database.DB, strconv.FormatInt(state.UserId, 10))
	if err != nil {
		bot.Send(tgbotapi.NewMessage(state.ChatId, "خطا در شناسایی کاربر!"))
		botLogger.Error(
			"Error retrieving user ID from Telegram ID",
			zap.Error(err),
			zap.String("telegram_id", strconv.FormatInt(state.UserId, 10)),
		)
		return
	}

	// Retrieve user role
	user, err := users.GetUserByIDService(database.DB, userID)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(state.ChatId, "خطا در دریافت اطلاعات کاربر!"))
		botLogger.Error(
			"Error retrieving user role",
			zap.Error(err),
			zap.String("user_id", userID),
		)
		return
	}

	// Extract page number from callback data (if provided)
	page := 1 // Default page
	if update.CallbackQuery != nil && update.CallbackQuery.Data != "" {
		action := update.CallbackQuery.Data
		if len(action) > len("/see_all_filters:") && action[:len("/see_all_filters:")] == "/see_all_filters:" {
			pageStr := action[len("/see_all_filters:"):]
			if p, err := strconv.Atoi(pageStr); err == nil {
				page = p
			}
		}
	}

	// Define page size
	pageSize := 2

	// Fetch filters based on user role
	var filterData filterService.PaginatedFilters
	if user.Role == models.RoleUser {
		// Call GetFiltersByUserID for normal users
		filterData, err = filterService.GetFiltersByUserID(database.DB, userID, page, pageSize)
	} else {
		// Call GetAllFilters for admins and super admins
		filterData, err = filterService.GetAllFilters(database.DB, userID, page, pageSize)
	}

	if err != nil {
		bot.Send(tgbotapi.NewMessage(state.ChatId, "خطا در دریافت فیلترها"))
		botLogger.Error("Error fetching filters", zap.Error(err))
		return
	}

	// If no filters are found
	if len(filterData.Data) == 0 {
		msg := tgbotapi.NewMessage(state.ChatId, "فیلتری برای نمایش یافت نشد.")
		bot.Send(msg)
		return
	}

	// Generate response
	response := fmt.Sprintf("📋 *فیلترهای موجود (صفحه %d از %d):*\n\n", filterData.PageIndex, filterData.TotalPages)
	var buttons [][]tgbotapi.InlineKeyboardButton

	for _, filter := range filterData.Data {
		response += fmt.Sprintf(
			"🏷️ *عنوان:* %s\n"+
				"🏙️ *شهر:* %s\n"+
				"📍 *محله:* %s\n"+
				"🏢 *نوع ملک:* %s\n"+
				"📐 *حداقل متراژ:* %d\n"+
				"📐 *حداکثر متراژ:* %d\n"+
				"💰 *حداقل قیمت:* %d\n"+
				"💰 *حداکثر قیمت:* %d\n\n",
			filter.Title,
			filter.City,
			filter.Neighborhood,
			filter.PropertyType,
			filter.MinArea, filter.MaxArea,
			filter.MinPrice, filter.MaxPrice,
		)

		// Add an inline button for filter actions (e.g., View, Edit, Delete)
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("🔍 مشاهده: %s", filter.Title), fmt.Sprintf("/view_filter:%s", filter.ID)),
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("🗑️ حذف %s", filter.Title), fmt.Sprintf("/delete_filter:%s", filter.ID)),
		))
	}

	// Add pagination buttons
	if filterData.TotalPages > 1 {
		if filterData.PageIndex > 1 {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("⬅️ صفحه قبلی", fmt.Sprintf("/see_all_filters:%d", filterData.PageIndex-1)),
			))
		}
		if filterData.PageIndex < filterData.TotalPages {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("➡️ صفحه بعدی", fmt.Sprintf("/see_all_filters:%d", filterData.PageIndex+1)),
			))
		}
	}

	// Send the message with inline buttons
	msg := tgbotapi.NewMessage(state.ChatId, response)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
	msg.ParseMode = "Markdown"
	bot.Send(msg)

	// Update user state and action state
	err = userStates.SetUserCache(ctx, state.ChatId, cache.UserState{
		ChatId:       state.ChatId,
		UserId:       state.UserId,
		Stage:        "get_filters",
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
		Action:       "view_filters",
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
