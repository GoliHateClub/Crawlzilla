package filters

import (
	"Crawlzilla/database"
	cfg "Crawlzilla/logger"
	"Crawlzilla/models"
	"Crawlzilla/services/cache"
	"Crawlzilla/services/search"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func ExportFilteredResultsConversation(ctx context.Context, state cache.UserState, update tgbotapi.Update) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	// Extract the filter ID from the callback data
	action := update.CallbackQuery.Data
	filterID := strings.TrimPrefix(action, "/export_filter:")

	// Fetch all filtered ads
	page := 1
	pageSize := 100         // Fetch 100 records per page (adjustable for large datasets)
	var allAds []models.Ads // Corrected type

	for {
		adsData, err := search.GetFilteredAds(database.DB, filterID, page, pageSize)
		if err != nil {
			bot.Send(tgbotapi.NewMessage(state.ChatId, "خطا در دریافت نتایج جستجو!"))
			botLogger.Error("Error fetching filtered ads", zap.Error(err))
			return
		}

		allAds = append(allAds, adsData.Data...)

		// If we've fetched all pages, break
		if page >= adsData.Pages {
			break
		}

		page++
	}

	// If no ads found
	if len(allAds) == 0 {
		bot.Send(tgbotapi.NewMessage(state.ChatId, "هیچ آگهی‌ای مطابق با فیلتر یافت نشد."))
		return
	}

	// Create a temporary CSV file
	fileName := fmt.Sprintf("filtered_results_%s.csv", filterID)
	file, err := os.CreateTemp("", fileName) // Use `os.CreateTemp` for safer temp file creation
	if err != nil {
		bot.Send(tgbotapi.NewMessage(state.ChatId, "خطا در ایجاد فایل CSV!"))
		botLogger.Error("Error creating CSV file", zap.Error(err))
		return
	}
	defer file.Close()

	// Write data to CSV
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	err = writer.Write([]string{"ID", "Title", "City", "Neighborhood", "Price", "Rooms", "Area", "Details URL"})
	if err != nil {
		bot.Send(tgbotapi.NewMessage(state.ChatId, "خطا در نوشتن به فایل CSV!"))
		botLogger.Error("Error writing CSV header", zap.Error(err))
		return
	}

	// Write ad data
	for _, ad := range allAds {
		err := writer.Write([]string{
			ad.ID,
			ad.Title,
			ad.City,
			ad.Neighborhood,
			strconv.Itoa(ad.Price),
			strconv.Itoa(ad.Room),
			strconv.Itoa(ad.Area),
			fmt.Sprintf("https://yourdomain.com/view_ad:%s", ad.ID), // Replace with actual ad URL format
		})
		if err != nil {
			botLogger.Error("Error writing to CSV", zap.Error(err))
		}
	}

	// Flush writer to ensure all data is written
	writer.Flush()
	if err := writer.Error(); err != nil {
		bot.Send(tgbotapi.NewMessage(state.ChatId, "خطا در نوشتن داده‌ها به فایل CSV!"))
		botLogger.Error("Error flushing CSV writer", zap.Error(err))
		return
	}

	// Send the CSV file to the user
	msg := tgbotapi.NewDocument(state.ChatId, tgbotapi.FilePath(file.Name()))
	msg.Caption = "📄 فایل نتایج فیلتر"
	if _, err := bot.Send(msg); err != nil {
		botLogger.Error("Error sending CSV file to user", zap.Error(err))
		bot.Send(tgbotapi.NewMessage(state.ChatId, "خطا در ارسال فایل به کاربر!"))
		return
	}

	// Delete the CSV file after successful send
	if err := os.Remove(file.Name()); err != nil {
		botLogger.Error("Error deleting CSV file after sending", zap.Error(err))
	} else {
		botLogger.Info("CSV file deleted successfully", zap.String("file", file.Name()))
	}

	// Log successful file creation and sending
	botLogger.Info("CSV file successfully sent to user", zap.String("file", file.Name()), zap.Int("ads_count", len(allAds)))
}
