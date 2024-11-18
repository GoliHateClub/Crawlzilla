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

func ViewFilterDetailsConversation(ctx context.Context, state cache.UserState, update tgbotapi.Update) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)
	userStates := ctx.Value("user_state").(*cache.UserCache)
	actionStates := ctx.Value("action_state").(*cache.ActionCache)
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	// Extract filter ID from callback data
	var filterID string
	if update.CallbackQuery != nil && update.CallbackQuery.Data != "" {
		action := update.CallbackQuery.Data
		if len(action) > len("/view_filter:") && action[:len("/view_filter:")] == "/view_filter:" {
			filterID = action[len("/view_filter:"):]
		}
	}

	if filterID == "" {
		bot.Send(tgbotapi.NewMessage(state.ChatId, "خطا در دریافت شناسه فیلتر!"))
		return
	}

	// Get database user ID from Telegram ID
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

	// Fetch filter details by ID
	filter, err := filterService.GetFilterByID(database.DB, filterID)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(state.ChatId, "خطا در دریافت اطلاعات فیلتر!"))
		botLogger.Error("Error fetching filter details", zap.Error(err))
		return
	}

	// Fetch the role of the requesting user
	requestingUser, err := users.GetUserByIDService(database.DB, userID)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(state.ChatId, "خطا در شناسایی نقش کاربر!"))
		botLogger.Error("Error fetching user role", zap.Error(err))
		return
	}

	// Check access permissions
	if requestingUser.Role == models.RoleUser && filter.USER_ID != userID {
		bot.Send(tgbotapi.NewMessage(state.ChatId, "شما اجازه مشاهده این فیلتر را ندارید!"))
		return
	}

	// Prepare response based on role
	response := fmt.Sprintf(
		"🏷️ *عنوان:* %s\n"+
			"🏙️ *شهر:* %s\n"+
			"📍 *محله:* %s\n"+
			"📏 *مرجع:* %s\n"+
			"🏡 *نوع آگهی:* %s\n"+
			"🏢 *نوع ملک:* %s\n"+
			"📐 *حداقل متراژ:* %d\n"+
			"📐 *حداکثر متراژ:* %d\n"+
			"💰 *حداقل قیمت:* %d\n"+
			"💰 *حداکثر قیمت:* %d\n"+
			"💸 *حداقل اجاره:* %d\n"+
			"💸 *حداکثر اجاره:* %d\n"+
			"🚪 *حداقل تعداد اتاق:* %d\n"+
			"🚪 *حداکثر تعداد اتاق:* %d\n"+
			"🏗️ *حداقل تعداد طبقات:* %d\n"+
			"🏗️ *حداکثر تعداد طبقات:* %d\n"+
			"🚪 *آسانسور:* %s\n"+
			"📦 *انباری:* %s\n"+
			"🚗 *پارکینگ:* %s\n"+
			"🌳 *بالکن:* %s\n"+
			"🕓 *مرتب‌سازی بر اساس:* %s\n"+
			"🔀 *ترتیب:* %s\n"+
			"🆔 *تاریخ ایجاد:* %s\n",
		filter.Title,
		filter.City,
		filter.Neighborhood,
		filter.Reference,
		filter.CategoryType,
		filter.PropertyType,
		filter.MinArea, filter.MaxArea,
		filter.MinPrice, filter.MaxPrice,
		filter.MinRent, filter.MaxRent,
		filter.MinRoom, filter.MaxRoom,
		filter.MinFloorNumber, filter.MaxFloorNumber,
		boolToEmoji(filter.HasElevator),
		boolToEmoji(filter.HasStorage),
		boolToEmoji(filter.HasParking),
		boolToEmoji(filter.HasBalcony),
		sortKeyToName(filter.Sort),
		orderKeyToName(filter.Order),
		filter.CreatedAt.Format("2006-01-02 15:04:05"),
	)

	if requestingUser.Role == models.RoleSuperAdmin {
		// Include the Telegram ID of the filter owner for super-admins
		ownerTelegramID, err := users.GetUserByIDService(database.DB, filter.USER_ID)
		if err == nil {
			response += fmt.Sprintf("👤 *تلگرام کاربر:* `%d`\n", ownerTelegramID.Telegram_ID)
		} else {
			botLogger.Error("Error fetching Telegram ID for filter owner", zap.Error(err))
		}
	}

	// Add action buttons (Apply and Delete)
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("🗑️ حذف", fmt.Sprintf("/delete_filter:%s", filter.ID)),
			tgbotapi.NewInlineKeyboardButtonData("✅ اعمال", fmt.Sprintf("/apply_filter:%s", filter.ID)),
		},
	}

	// Send the response
	msg := tgbotapi.NewMessage(state.ChatId, response)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
	msg.ParseMode = "Markdown"
	bot.Send(msg)

	// Update user state and action state
	err = userStates.SetUserCache(ctx, state.ChatId, cache.UserState{
		ChatId:       state.ChatId,
		UserId:       state.UserId,
		Stage:        "view_filter_details",
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
		Action:       "view_filter_details",
		ActionData: map[string]interface{}{
			"filter_id": filterID,
		},
	})
	if err != nil {
		cache.HandleActionStateError(botLogger, state, err)
		bot.Send(tgbotapi.NewMessage(state.ChatId, "خطایی رخ داد!"))
		return
	}
}
func boolToEmoji(value bool) string {
	if value {
		return "✅ بله"
	}
	return "❌ خیر"
}

func sortKeyToName(sortKey string) string {
	switch sortKey {
	case "price":
		return "قیمت"
	case "rent":
		return "اجاره"
	case "area":
		return "مساحت"
	case "room":
		return "اتاق"
	case "floor_number":
		return "تعداد طبقات"
	case "visit_count":
		return "تعداد بازدید"
	case "created_at":
		return "تاریخ ایجاد"
	default:
		return "نامشخص"
	}
}

func orderKeyToName(order string) string {
	switch order {
	case "asc":
		return "صعودی"
	case "desc":
		return "نزولی"
	default:
		return "نامشخص"
	}
}
