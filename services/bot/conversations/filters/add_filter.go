package filters

import (
	"Crawlzilla/database"
	cfg "Crawlzilla/logger"
	"Crawlzilla/models"
	"Crawlzilla/services/bot/constants"
	"Crawlzilla/services/cache"
	filterService "Crawlzilla/services/filters"
	"Crawlzilla/services/users"
	"context"
	"log"
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
		bot.Send(tgbotapi.NewMessage(state.ChatId, "عین همون قالبی که بالا دیدی، برای فرم زیر رو بفرست:"))
		bot.Send(tgbotapi.NewMessage(state.ChatId, constants.NewFilterTextDetails))

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
		CityRegex := regexp.MustCompile(`(?i)شهر[:：\s]*(.+)`)
		NeighborhoodRegex := regexp.MustCompile(`(?i)محله[:：\s]*(.+)`)
		CategoryTypeRegex := regexp.MustCompile(`(?i)نوع آگهی[:：\s]*(.+)`)
		PropertyTypeRegex := regexp.MustCompile(`(?i)نوع ملک[:：\s]*(.+)`)
		HasElevatorRegex := regexp.MustCompile(`(?i)آسانسور داشته باشید؟[:：\s]*(بله|خیر)`)
		HasStorageRegex := regexp.MustCompile(`(?i)انباری داشته باشید؟[:：\s]*(بله|خیر)`)
		HasParkingRegex := regexp.MustCompile(`(?i)پارکینگ داشته باشید؟[:：\s]*(بله|خیر)`)
		HasBalconyRegex := regexp.MustCompile(`(?i)بالکن داشته باشید؟[:：\s]*(بله|خیر)`)

		City := CityRegex.FindStringSubmatch(input)
		Neighborhood := NeighborhoodRegex.FindStringSubmatch(input)
		CategoryType := CategoryTypeRegex.FindStringSubmatch(input)
		PropertyType := PropertyTypeRegex.FindStringSubmatch(input)
		HasElevator := HasElevatorRegex.FindStringSubmatch(input)
		HasStorage := HasStorageRegex.FindStringSubmatch(input)
		HasParking := HasParkingRegex.FindStringSubmatch(input)
		HasBalcony := HasBalconyRegex.FindStringSubmatch(input)

		// Create a filter structure
		var filter models.Filters
		filter.Title = title // Use the title from the previous step

		if City != nil {
			filter.City = City[1]
		}
		if Neighborhood != nil {
			filter.Neighborhood = Neighborhood[1]
		}
		if CategoryType != nil {
			filter.CategoryType = CategoryType[1]
		}
		if PropertyType != nil {
			filter.PropertyType = PropertyType[1]
		}
		if HasElevator != nil {
			filter.HasElevator = HasElevator[1] == "بله"
		}
		if HasStorage != nil {
			filter.HasStorage = HasStorage[1] == "بله"
		}
		if HasParking != nil {
			filter.HasParking = HasParking[1] == "بله"
		}
		if HasBalcony != nil {
			filter.HasBalcony = HasBalcony[1] == "بله"
		}

		// Set the user ID
		userId, err := users.GetUserIDByTelegramID(database.DB, strconv.FormatInt(state.UserId, 10))
		if err != nil {
			botLogger.Error(
				"Error while reading user id",
				zap.Error(err),
				zap.String("user_id", strconv.FormatInt(state.UserId, 10)),
			)
			bot.Send(tgbotapi.NewMessage(state.ChatId, "خطایی رخ داد! لطفاً دوباره تلاش کنید."))
			return
		}
		filter.USER_ID = userId

		// Save the filter
		_, err = filterService.CreateOrUpdateFilter(database.DB, filter)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(state.ChatId, "خطایی هنگام ذخیره سازی اطلاعات رخ داد! لطفاً دوباره تلاش کنید."))
			return
		}

		// Clear states and inform the user
		err = userStates.ClearUserCache(ctx, state.ChatId)
		if err != nil {
			cache.HandleUserStateError(botLogger, state, err)
			return
		}
		err = actionStates.ClearActionState(ctx, state.ChatId)
		if err != nil {
			cache.HandleActionStateError(botLogger, state, err)
			return
		}

		bot.Send(tgbotapi.NewMessage(state.ChatId, "فیلتر جدید با موفقیت ذخیره شد!"))
	}
}
