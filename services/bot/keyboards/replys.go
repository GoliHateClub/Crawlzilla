package keyboards

import (
	"Crawlzilla/services/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func ReplyKeyboardMain(isAdmin bool) tgbotapi.InlineKeyboardMarkup {
	menu := bot.MainMenu
	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	for _, items := range menu {
		var row []tgbotapi.InlineKeyboardButton

		for _, item := range items {

			if item.IsAdmin && !isAdmin {
				continue
			}
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(item.Name, item.Path))
		}

		if len(row) > 0 {
			keyboardRows = append(keyboardRows, row)
		}
	}

	// Return the complete keyboard markup
	return tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)
}
