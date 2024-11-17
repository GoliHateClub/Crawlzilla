package conversations

import (
	"Crawlzilla/database"
	"Crawlzilla/models"
	"Crawlzilla/services/ads"
	"Crawlzilla/services/bot/constants"
	"Crawlzilla/services/cache"
	"Crawlzilla/services/super_admin"
	"Crawlzilla/utils"
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"go.uber.org/zap"

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
			return
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

func GetAdDetailsConversation(ctx context.Context, update tgbotapi.Update) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	// Extract the ad ID from the callback query
	if update.CallbackQuery == nil || update.CallbackQuery.Data == "" {
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "❌ خطا: شناسه آگهی موجود نیست!"))
		botLogger.Error("Callback query data is empty or missing.")
		return
	}

	action := update.CallbackQuery.Data
	if len(action) <= len("/view_ad:") || action[:len("/view_ad:")] != "/view_ad:" {
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "❌ خطا: فرمت شناسه آگهی نامعتبر است!"))
		botLogger.Error("Invalid callback query format.")
		return
	}

	adID := action[len("/view_ad:"):]

	// Fetch ad details using the service layer
	ad, err := ads.GetAdById(database.DB, adID)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "❌ خطا در دریافت جزئیات آگهی!"))
		botLogger.Error(
			"Error fetching ad details",
			zap.String("ad_id", adID),
			zap.Error(err),
		)
		return
	}

	// Format ad details into a user-friendly message with emojis
	response := fmt.Sprintf(
		"📋 *جزئیات آگهی:*\n\n"+
			"🏷️ *عنوان:* %s\n"+
			"📝 *توضیحات:* %s\n"+
			"📍 *شهر:* %s\n"+
			"🏘️ *محله:* %s\n"+
			"📐 *مساحت:* %d متر مربع\n"+
			"💰 *قیمت:* %d تومان\n"+
			"📞 *شماره تماس:* %s",
		ad.Title, ad.Description, ad.City, ad.Neighborhood, ad.Area, ad.Price, ad.ContactNumber,
	)

	// Decide the message type based on the presence of an image URL
	chatID := update.CallbackQuery.Message.Chat.ID
	if ad.ImageURL != "" {
		// Send a photo message with details in the caption
		photoMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(ad.ImageURL))
		photoMsg.Caption = response
		photoMsg.ParseMode = "Markdown" // Enable Markdown for formatting
		bot.Send(photoMsg)
	} else {
		// Send a regular text message
		msg := tgbotapi.NewMessage(chatID, response)
		msg.ParseMode = "Markdown" // Enable Markdown for formatting
		bot.Send(msg)
	}

	// Acknowledge the callback to prevent loading spinner in the UI
	bot.Send(tgbotapi.NewCallback(update.CallbackQuery.ID, "جزئیات آگهی ارسال شد."))
}
