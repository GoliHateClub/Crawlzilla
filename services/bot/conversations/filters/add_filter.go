package filters

import (
	"Crawlzilla/database"
	cfg "Crawlzilla/logger"
	"Crawlzilla/models"
	"Crawlzilla/services/cache"
	filterService "Crawlzilla/services/filters"
	"Crawlzilla/services/users"
	"context"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func AddFilterConversation(ctx context.Context, state cache.UserState, update tgbotapi.Update) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)
	userStates := ctx.Value("user_state").(*cache.UserCache)
	actionStates := ctx.Value("action_state").(*cache.ActionCache)
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	switch state.Stage {
	case "init":
		// Inform the user to provide a filter title
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "نام فیلتر جدیدت رو بنویس")
		bot.Send(msg)

		// Update the user state to move to the next stage
		err := userStates.SetUserCache(ctx, state.ChatId, cache.UserState{
			ChatId:       state.ChatId,
			UserId:       state.UserId,
			Stage:        "get_title_ask_ranges",
			Conversation: state.Conversation,
		})

		if err != nil {
			cache.HandleUserStateError(botLogger, state, err)
			log.Printf("Error updating user state: %v", err)
			return
		}

	case "get_title_ask_ranges":
		// Step 1: Receive only the title
		title := strings.TrimSpace(update.Message.Text)
		if title == "" {
			bot.Send(tgbotapi.NewMessage(state.ChatId, "عنوان نمی‌تواند خالی باشد! لطفاً یک عنوان معتبر وارد کنید."))
			return
		}

		// Step 2: Fix user state and action state
		err := actionStates.SetUserState(ctx, state.ChatId, cache.ActionState{
			ChatId:       state.ChatId,
			UserId:       state.UserId,
			Conversation: "add_filter", // Conversation name
			Action:       "add_filter",
			ActionData: map[string]interface{}{
				"title": title, // Save the title in action state
			},
		})
		if err != nil {
			cache.HandleActionStateError(botLogger, state, err)
			bot.Send(tgbotapi.NewMessage(state.ChatId, "خطایی رخ داد! لطفاً دوباره تلاش کنید."))
			return
		}

		err = userStates.SetUserCache(ctx, state.ChatId, cache.UserState{
			ChatId:       state.ChatId,
			UserId:       state.UserId,
			Stage:        "ask_text_details", // Move to the next stage
			Conversation: state.Conversation,
		})
		if err != nil {
			cache.HandleUserStateError(botLogger, state, err)
			bot.Send(tgbotapi.NewMessage(state.ChatId, "خطایی رخ داد! لطفاً دوباره تلاش کنید."))
			return
		}

		// Step 3: Inform the user
		bot.Send(tgbotapi.NewMessage(state.ChatId, "مرسی ازت! حالا می‌خوام که یکم جزئیات متنی بهم بدی."))
		bot.Send(tgbotapi.NewMessage(state.ChatId, `لطفاً اطلاعات را به این صورت وارد کن:
		
شهر: تهران  
مرجع: دیوار/شیپور/ادمین 
محله: تجریش  
نوع آگهی: فروش  
نوع ملک: آپارتمانی  
حداقل متراژ: 50  
حداکثر متراژ: 200  
حداقل قیمت: 2000000000  
حداکثر قیمت: 10000000000  
حداقل اجاره: 2,000,000  
حداکثر اجاره: 10,000,000  
حداقل تعداد اتاق: 2  
حداکثر تعداد اتاق: 4  
حداقل تعداد طبقه: 1  
حداکثر تعداد طبقه: 5  
آسانسور داشته باشد؟ بله  
انباری داشته باشد؟ خیر  
پارکینگ داشته باشد؟ بله  
بالکن داشته باشد؟ بله  
مرتب سازی: قیمت | اجاره | مساحت | اتاق | طبقه | تعداد بازدید | تاریخ ایجاد  
ترتیب: سعودی | نزولی`))

	case "ask_text_details":
		// Retrieve title from action state
		actionState, err := actionStates.GetActionState(ctx, state.ChatId)
		if err != nil {
			cache.HandleActionStateError(botLogger, state, err)
			bot.Send(tgbotapi.NewMessage(state.ChatId, "خطایی رخ داد! لطفاً دوباره تلاش کنید."))
			return
		}

		title, ok := actionState.ActionData["title"].(string)
		if !ok || title == "" {
			bot.Send(tgbotapi.NewMessage(state.ChatId, "خطایی رخ داد! عنوان یافت نشد. لطفاً دوباره تلاش کنید."))
			return
		}

		// Parse user input for text details
		input := strings.TrimSpace(update.Message.Text)
		fields := map[string]string{
			"City":         `(?i)شهر[:：\s]*(.+)`,
			"Reference":    `(?i)مرجع[:：\s]*(.+)`,
			"Neighborhood": `(?i)محله[:：\s]*(.+)`,
			"CategoryType": `(?i)نوع آگهی[:：\s]*(.+)`,
			"PropertyType": `(?i)نوع ملک[:：\s]*(.+)`,
			"Sort":         `(?i)مرتب سازی[:：\s]*(.+)`,
			"Order":        `(?i)ترتیب[:：\s]*(سعودی|نزولی)`,
		}

		numericFields := map[string]string{
			"MinArea":        `(?i)حداقل متراژ[:：\s]*(\d+)`,
			"MaxArea":        `(?i)حداکثر متراژ[:：\s]*(\d+)`,
			"MinPrice":       `(?i)حداقل قیمت[:：\s]*([\d,]+)`,
			"MaxPrice":       `(?i)حداکثر قیمت[:：\s]*([\d,]+)`,
			"MinRent":        `(?i)حداقل اجاره[:：\s]*([\d,]+)`,
			"MaxRent":        `(?i)حداکثر اجاره[:：\s]*([\d,]+)`,
			"MinRoom":        `(?i)حداقل تعداد اتاق[:：\s]*(\d+)`,
			"MaxRoom":        `(?i)حداکثر تعداد اتاق[:：\s]*(\d+)`,
			"MinFloorNumber": `(?i)حداقل تعداد طبقه[:：\s]*(\d+)`,
			"MaxFloorNumber": `(?i)حداکثر تعداد طبقه[:：\s]*(\d+)`,
		}

		booleanFields := map[string]string{
			"HasElevator": `(?i)آسانسور داشته باشد؟[:：\s]*(بله|خیر)`,
			"HasStorage":  `(?i)انباری داشته باشد؟[:：\s]*(بله|خیر)`,
			"HasParking":  `(?i)پارکینگ داشته باشد؟[:：\s]*(بله|خیر)`,
			"HasBalcony":  `(?i)بالکن داشته باشد؟[:：\s]*(بله|خیر)`,
		}

		var filter models.Filters
		filter.Title = title // Use the title from the previous step

		// Map textual fields
		for field, pattern := range fields {
			value := extractField(pattern, input)
			if field == "Sort" {
				value = mapSortField(value)
			} else if field == "Order" {
				value = mapOrderField(value)
			}
			reflect.ValueOf(&filter).Elem().FieldByName(field).SetString(value)
		}

		// Map numeric fields
		for field, pattern := range numericFields {
			value := parseToInt(pattern, input)
			reflect.ValueOf(&filter).Elem().FieldByName(field).SetInt(int64(value))
		}

		// Map boolean fields
		for field, pattern := range booleanFields {
			value := extractBoolean(pattern, input)
			reflect.ValueOf(&filter).Elem().FieldByName(field).SetBool(value)
		}

		// Set the user ID
		userId, err := users.GetUserIDByTelegramID(database.DB, strconv.FormatInt(state.UserId, 10))
		if err != nil {
			botLogger.Error("Error while reading user ID", zap.Error(err))
			bot.Send(tgbotapi.NewMessage(state.ChatId, "خطایی رخ داد! لطفاً دوباره تلاش کنید."))
			return
		}
		filter.USER_ID = userId

		// Save the filter
		_, err = filterService.CreateOrUpdateFilter(database.DB, filter)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(state.ChatId, "خطایی هنگام ذخیره‌سازی اطلاعات رخ داد! لطفاً دوباره تلاش کنید."))
			return
		}

		// Clear states and inform the user
		userStates.ClearUserCache(ctx, state.ChatId)
		actionStates.ClearActionState(ctx, state.ChatId)
		bot.Send(tgbotapi.NewMessage(state.ChatId, "فیلتر جدید با موفقیت ذخیره شد!"))
	}
}

// Utility Functions
func mapSortField(value string) string {
	switch strings.TrimSpace(value) {
	case "قیمت":
		return "price"
	case "اجاره":
		return "rent"
	case "مساحت":
		return "area"
	case "اتاق":
		return "room"
	case "طبقه":
		return "floor_number"
	case "تعداد بازدید":
		return "visit_count"
	case "تاریخ ایجاد":
		return "created_at"
	default:
		return ""
	}
}

func mapOrderField(value string) string {
	switch strings.TrimSpace(value) {
	case "سعودی":
		return "asc"
	case "نزولی":
		return "desc"
	default:
		return ""
	}
}

// Helper to extract a field value using regex
func extractField(pattern string, input string) string {
	regex := regexp.MustCompile(pattern)
	match := regex.FindStringSubmatch(input)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

// Helper to parse an integer from a regex-matched field
func parseToInt(pattern string, input string) int {
	valueStr := extractField(pattern, input)
	valueStr = strings.ReplaceAll(valueStr, ",", "") // Remove commas for large numbers
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0
	}
	return value
}

// Helper to extract a boolean field
func extractBoolean(pattern string, input string) bool {
	value := extractField(pattern, input)
	return value == "بله"
}
