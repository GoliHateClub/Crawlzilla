package filters

import (
	cfg "Crawlzilla/logger"
	"Crawlzilla/models"
	"Crawlzilla/services/bot/constants"
	"Crawlzilla/services/cache"
	"Crawlzilla/utils"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var questions = map[string]struct {
	Name        string
	Description string
}{
	"Area": {
		Name:        "متراژ",
		Description: "متراژ خونه چقدر باشه؟",
	},
	"Price": {
		Name:        "قیمت",
		Description: "قیمت خونه چقدر باشه؟",
	},
	"Rent": {
		Name:        "اجاره",
		Description: "اجاره خونه چقدر باشه؟",
	},
	"Room": {
		Name:        "تعداد اتاق",
		Description: "تعداد اتاق خونه چقدر باشه؟",
	},
	"FloorNumber": {
		Name:        "تعداد طبقه",
		Description: "خونه تو کدوم طبقه باشه؟",
	},
}

func getNextQuestionKey(currentKey string) string {
	keys := []string{"Area", "Price", "Rent", "Room", "FloorNumber"}
	for i, key := range keys {
		if key == currentKey && i+1 < len(keys) {
			return keys[i+1]
		}
	}
	return ""
}

func AddFilterConversation(ctx context.Context, state cache.UserState, update tgbotapi.Update) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)
	userStates := ctx.Value("user_state").(*cache.UserCache)
	actionStates := ctx.Value("action_state").(*cache.ActionCache)
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	switch state.Stage {
	case "init":
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "نام فیلتر جدیدت رو بنویس")
		bot.Send(msg)

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
		userAction, err := actionStates.GetActionState(ctx, state.ChatId)

		if err != nil {
			cache.HandleActionStateError(botLogger, state, err)
			bot.Send(tgbotapi.NewMessage(state.ChatId, "خطایی رخ داد!"))
			return
		}

		if userAction.IsEmpty() {
			title := strings.Trim(update.Message.Text, " ")

			newFilter := models.Filters{
				Title: title,
			}

			err = actionStates.SetUserState(ctx, state.ChatId, cache.ActionState{
				ChatId:       state.ChatId,
				UserId:       state.UserId,
				Conversation: state.Conversation,
				Action:       "add_filter",
				ActionData: map[string]interface{}{
					"current_q":  "Area",
					"new_filter": newFilter,
				},
			})

			if err != nil {
				cache.HandleActionStateError(botLogger, state, err)
				bot.Send(tgbotapi.NewMessage(state.ChatId, "خطایی رخ داد!"))
				return
			}

			bot.Send(tgbotapi.NewMessage(state.ChatId, "حالا باید مقادیر عددی مثل قیمت، اجاره تعداد اتاق و... رو بفرستی."))
			bot.Send(tgbotapi.NewMessage(state.ChatId, "مقادیری که میفرستی باید در قالب زیر باشه:"))
			bot.Send(tgbotapi.NewMessage(state.ChatId, constants.NewFilterRange))
			bot.Send(tgbotapi.NewMessage(state.ChatId, constants.NewFilterRangeRules))
			bot.Send(tgbotapi.NewMessage(state.ChatId, questions["Area"].Description))
			return
		}

		filterData, err := actionStates.GetActionState(ctx, state.ChatId)
		filterMap, _ := filterData.ActionData["new_filter"]
		filterState, _ := filterData.ActionData["current_q"].(string)

		var filters models.Filters

		err = utils.MapToStruct(filterMap, &filters)

		nextQuestionKey := getNextQuestionKey(filterState)

		// Normalize Persian numbers to English
		input := strings.NewReplacer(
			"۰", "0", "۱", "1", "۲", "2", "۳", "3",
			"۴", "4", "۵", "5", "۶", "6", "۷", "7",
			"۸", "8", "۹", "9",
		).Replace(update.Message.Text)

		// Regex to match حداقل and حداکثر values
		minRegex := regexp.MustCompile(`(?i)حداقل[:：\s]*([\d]+)`)
		maxRegex := regexp.MustCompile(`(?i)حداکثر[:：\s]*([\d]+)`)

		// Find matches
		minMatch := minRegex.FindStringSubmatch(input)
		maxMatch := maxRegex.FindStringSubmatch(input)

		var minValue int
		var maxValue int

		// Parse matches to integers
		if len(minMatch) > 1 {
			minValue, err = strconv.Atoi(minMatch[1])
			if err != nil {
				cache.HandleActionStateError(botLogger, state, err)
				msg := tgbotapi.NewMessage(state.ChatId, fmt.Sprintf("ساختار مقادیر %s شما با ساختار تعریف شده مطابقت ندارد!", questions[filterState].Name))
				bot.Send(msg)
				return
			}
		}

		if len(maxMatch) > 1 {
			maxValue, err = strconv.Atoi(maxMatch[1])
			if err != nil {
				cache.HandleActionStateError(botLogger, state, err)
				msg := tgbotapi.NewMessage(state.ChatId, fmt.Sprintf("ساختار مقادیر %s شما با ساختار تعریف شده مطابقت ندارد!", questions[filterState].Name))
				bot.Send(msg)
				return
			}
		}

		if minValue != 0 {
			fieldName := "Min" + filterState
			reflect.ValueOf(&filters).Elem().FieldByName(fieldName).SetInt(int64(minValue))
		}
		if maxValue != 0 {
			fieldName := "Max" + filterState
			reflect.ValueOf(&filters).Elem().FieldByName(fieldName).SetInt(int64(maxValue))
		}

		err = actionStates.UpdateActionCache(ctx, state.ChatId, map[string]interface{}{
			"ActionData": map[string]interface{}{
				"new_filter": filters,
				"current_q":  nextQuestionKey,
			},
		})

		if err != nil {
			cache.HandleUserStateError(botLogger, state, err)
			msg := tgbotapi.NewMessage(state.ChatId, "خطایی رخ داد")
			bot.Send(msg)
			return
		}

		if nextQuestionKey != "" {
			nextQuestion := questions[nextQuestionKey]
			bot.Send(tgbotapi.NewMessage(state.ChatId, nextQuestion.Description))
			return
		}

		bot.Send(tgbotapi.NewMessage(state.ChatId, "سوالات تمام شد. لطفا اطلاعات خود را ذخیره کنید.")) // ToDO

	case "get_ranges_ask_text_details":

	}
}
