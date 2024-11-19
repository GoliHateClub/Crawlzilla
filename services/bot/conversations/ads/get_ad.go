package ads

import (
	"Crawlzilla/database"
	cfg "Crawlzilla/logger"
	"Crawlzilla/services/ads"
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

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
			"💰 *اجاره:* %d تومان\n"+
			"📞 *شماره تماس:* %s\n"+
			"💰 *تاریخ:* %v \n"+
			"*مرجع:* %v \n"+
			"*طبقه:* %v \n"+
			"*کل طبقه:* %v \n"+
			"*تغداد اتاق:* %v \n"+
			"*نوع آگهی:* %v \n"+
			"*نوع ملک:* %v \n",
		ad.Title, ad.Description, ad.City, ad.Neighborhood, ad.Area, ad.Price, ad.Rent, ad.ContactNumber, ad.CreatedAt, ad.Reference, ad.FloorNumber, ad.TotalFloors, ad.Room, ad.CategoryType, ad.PropertyType,
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
