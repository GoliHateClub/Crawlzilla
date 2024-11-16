package conversations

import (
	"Crawlzilla/database"
	"Crawlzilla/models"
	"Crawlzilla/services/bot/constants"
	"Crawlzilla/services/cache"
	"Crawlzilla/services/super_admin"
	"Crawlzilla/utils"
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"regexp"
	"strconv"
	"strings"

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
		pattern := `(?m)- نام\s*(.+)\n- شهر\s*(.+)\n- محله\s*(.+)\n- نوع آگهی \(فروش/خرید\)\s*(.+)\n- نوع ملک \(آپارتمانی/ویلایی\)\s*(.+)\n- متراژ\s*(\d+)\n- قیمت \(به تومان\)\s*(\d+)\n- اجاره \(به تومان\)\s*(\d+)\n- تعداد اتاق\s*(\d+)\n- طبقه واحد\s*(\d+)\n- تعداد طبقات ملک\s*(\d+)\n- شماره تماس \(همراه با 0\)\s*(\d+)\n- آسانسور دارد؟ \(فقط آپارتمانی\)\s*(بله|خیر)\n- انباری دارد؟\s*(بله|خیر)\n- پارکینگ دارد؟\s*(بله|خیر)\n- بالکن دارد؟\s*(بله|خیر)`

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

		area := handeConvertError(matches[6])
		price := handeConvertError(matches[7])
		rent := handeConvertError(matches[8])
		room := handeConvertError(matches[9])
		floorNumber := handeConvertError(matches[10])
		totalFloors := handeConvertError(matches[11])

		newAd := models.Ads{
			Title:         matches[1],
			City:          matches[2],
			Neighborhood:  matches[3],
			CategoryType:  matches[4],
			PropertyType:  matches[5],
			Area:          area,
			Price:         price,
			Rent:          rent,
			Room:          room,
			FloorNumber:   floorNumber,
			TotalFloors:   totalFloors,
			VisitCount:    0,
			ContactNumber: matches[12],
			HasElevator:   isBool(matches[13]),
			HasStorage:    isBool(matches[14]),
			HasParking:    isBool(matches[15]),
			HasBalcony:    isBool(matches[16]),
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
			cache.HandleActionStateError(botLogger, state, err)
			bot.Send(tgbotapi.NewMessage(state.ChatId, "خطایی رخ داد!"))
			return
		}

		err = userStates.UpdateUserCache(ctx, state.ChatId, map[string]interface{}{
			"Stage": "check_data_get_location",
		})

		if err != nil {
			cache.HandleUserStateError(botLogger, state, err)
			bot.Send(tgbotapi.NewMessage(state.ChatId, "خطایی رخ داد!"))
			return
		}

		msg := tgbotapi.NewMessage(state.ChatId, "حالا توضیحات آگهی خودت رو در قالب یک پیام بنویس و ارسال کن")
		bot.Send(msg)
	case "check_data_get_location":
		description := strings.Trim(update.Message.Text, " ")

		acData, err := actionStates.GetActionState(ctx, state.ChatId)
		adMap, _ := acData.ActionData["new_ad"]
		var ad models.Ads

		err = utils.MapToStruct(adMap, &ad)

		if err != nil {
			bot.Send(tgbotapi.NewMessage(state.ChatId, "خطایی هنگام پردازش داده رخ داد"))
			botLogger.Error(
				"Error while parsing ad data",
				zap.String("user_id", strconv.FormatInt(state.UserId, 10)),
				zap.Error(err),
			)
			return
		}

		fmt.Printf("%v\n", ad)

		ad.Description = description

		err = actionStates.UpdateActionCache(ctx, state.ChatId, map[string]interface{}{
			"ActionData": map[string]interface{}{
				"new_ad": ad,
			},
		})

		if err != nil {
			cache.HandleActionStateError(botLogger, state, err)
			bot.Send(tgbotapi.NewMessage(state.ChatId, "خطایی رخ داد!"))
			return
		}

		err = userStates.UpdateUserCache(ctx, state.ChatId, map[string]interface{}{
			"Stage": "save",
		})

		if err != nil {
			cache.HandleUserStateError(botLogger, state, err)
			bot.Send(tgbotapi.NewMessage(state.ChatId, "خطایی رخ داد!"))
			return
		}

		msg := tgbotapi.NewMessage(state.ChatId, "حالا با استفاده از لوکیشن تلگرام، لوکیشن مکان مورد نظر رو بفرست.")
		bot.Send(msg)
	case "save":
		if update.Message.Location == nil {
			msg := tgbotapi.NewMessage(state.ChatId, "آدرس ارسالی معتبر نبود! از امکان ارسال لوکیشن تلگرام استفاده کنید!")
			bot.Send(msg)
			return
		}

		acData, err := actionStates.GetActionState(ctx, state.ChatId)
		adMap, _ := acData.ActionData["new_ad"]
		var ad models.Ads

		err = utils.MapToStruct(adMap, &ad)

		if err != nil {
			bot.Send(tgbotapi.NewMessage(state.ChatId, "خطایی هنگام پردازش داده رخ داد"))
			botLogger.Error(
				"Error while parsing ad data",
				zap.String("user_id", strconv.FormatInt(state.UserId, 10)),
				zap.Error(err),
			)
			return
		}

		ad.Latitude = update.Message.Location.Latitude
		ad.Longitude = update.Message.Location.Longitude

		err = super_admin.CreateAd(database.DB, &ad)

		if err != nil {
			println(err)
			msg := tgbotapi.NewMessage(state.ChatId, "خطایی هنگام ذخیره سازی اطلاعات رخ داد! لطفا دوباره تلاش کنید.")
			bot.Send(msg)
			msg = tgbotapi.NewMessage(state.ChatId, err.Error())
			bot.Send(msg)
		}

		err = userStates.ClearUserCache(ctx, state.ChatId)

		if err != nil {
			cache.HandleUserStateError(botLogger, state, err)
			msg := tgbotapi.NewMessage(state.ChatId, "خطایی رخ داد")
			bot.Send(msg)
			return
		}

		err = actionStates.ClearActionState(ctx, state.ChatId)

		if err != nil {
			cache.HandleActionStateError(botLogger, state, err)
			msg := tgbotapi.NewMessage(state.ChatId, "خطایی رخ داد")
			bot.Send(msg)
			return
		}

		msg := tgbotapi.NewMessage(state.ChatId, "آگهی جدید با موفقیت ذخیره شد!")
		bot.Send(msg)
	}
}
