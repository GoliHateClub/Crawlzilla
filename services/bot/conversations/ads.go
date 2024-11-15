package conversations

import (
	"Crawlzilla/models"
	"Crawlzilla/services/bot/constants"
	"Crawlzilla/services/cache"
	"context"
	"go.uber.org/zap"
	"log"
	"regexp"
	"strconv"

	cfg "Crawlzilla/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AddAdConversation(ctx context.Context, state cache.UserState, update tgbotapi.Update) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)
	userStates := ctx.Value("user_state").(*cache.UserCache)
	actionStates := ctx.Value("action_state").(*cache.ActionCache)
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	switch state.Stage {
	case "init":
		msg := tgbotapi.NewMessage(state.ChatId, constants.NewAdText)
		msg.ReplyMarkup = nil
		bot.Send(msg)

		msg = tgbotapi.NewMessage(state.ChatId, constants.NewAdExampleText)
		bot.Send(msg)

		err := userStates.SetUserCache(ctx, state.ChatId, cache.UserState{
			ChatId:       state.ChatId,
			UserId:       state.UserId,
			Stage:        "check_data_get_desc",
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
	case "check_data_get_desc":
		// Define a regex pattern to extract each field
		pattern := `(?m)- نام\s*(.+)\n- آدرس لوکیشن در بلد\s*(.+)\n- آدرس URL محصول در سایت\s*(.+)\n- شهر\s*(.+)\n- محله\s*(.+)\n- منبع آگهی\s*(.+)\n- نوع آگهی\s*(.+)\n- نوع ملک\s*(.+)\n- متراژ\s*(\d+)\n- قیمت\(به ریال\)\s*(\d+)\n- اجاره\(ریال\)\s*(\d+)\n- تعداد اتاق\s*(\d+)\n- طبقه واحد\s*(\d+)\n- تعداد طبقات ملک\s*(\d+)\n- شماره تماس\s*(\d+)\n- آسانسور دارد؟\s*(بله|خیر)\n- انباری دارد؟\s*(بله|خیر)\n- پارکینگ دارد؟\s*(بله|خیر)\n- بالکن دارد؟\s*(بله|خیر)`

		re := regexp.MustCompile(pattern)

		// Match the input with the pattern
		matches := re.FindStringSubmatch(update.Message.Text)

		if matches == nil {
			msg := tgbotapi.NewMessage(state.ChatId, "ساختار آگهی شما با ساختار تعریف شده مطابقت ندارد!")
			bot.Send(msg)
			return
		}

		handeConvertError := func(data string) int {
			result, err := strconv.Atoi(matches[9])
			if err != nil {

			}
			return result
		}

		isBool := func(s string) bool {
			return s == "بله"
		}

		area := handeConvertError(matches[9])
		price := handeConvertError(matches[10])
		rent := handeConvertError(matches[11])
		room := handeConvertError(matches[12])
		floorNumber := handeConvertError(matches[13])
		totalFloors := handeConvertError(matches[14])

		newAd := models.Ads{
			Title:        matches[1],
			LocationURL:  matches[2],
			URL:          matches[3],
			City:         matches[4],
			Neighborhood: matches[5],
			Reference:    matches[6],
			CategoryType: matches[7],
			PropertyType: matches[8],
			Area:         area,
			Price:        price,
			Rent:         rent,
			Room:         room,
			FloorNumber:  floorNumber,
			TotalFloors:  totalFloors,
			VisitCount:   0,
			HasElevator:  isBool(matches[16]),
			HasStorage:   isBool(matches[17]),
			HasParking:   isBool(matches[18]),
			HasBalcony:   isBool(matches[19]),
		}

		err := actionStates.SetUserState(ctx, state.ChatId, cache.ActionState{
			ChatId:       state.ChatId,
			UserId:       state.UserId,
			Conversation: state.Conversation,
			Action:       "add_add_st1",
			ActionData: map[string]interface{}{
				"new_ad": newAd,
			},
		})

		if err != nil {
			botLogger.Error(
				"Error updating action state",
				zap.Error(err),
				zap.String("user_id", strconv.Itoa(int(state.UserId))),
				zap.String("chat_id", strconv.Itoa(int(state.ChatId))),
				zap.String("action", "add_add_st1"),
			)
			log.Printf("Error updating user state: %v", err)
			return
		}

		err = userStates.UpdateUserCache(ctx, state.ChatId, map[string]interface{}{
			"Stage": "check_data_get_location",
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

		msg := tgbotapi.NewMessage(state.ChatId, "حالا توضیحات آگهی خودت رو در قالب یک پیام بنویس و ارسال کن")
		bot.Send(msg)
	case "check_data_get_location":

	}
}
