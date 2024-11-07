package views

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bs *BotServer) TextView(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bs.bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
	return err
}
func (bs *BotServer) PairButtonView(btns [][2]string, msgTxt string, chatID int64) error {
	var rows [][]tgbotapi.InlineKeyboardButton
	for i := 0; i < len(btns); i += 2 {
		var row []tgbotapi.InlineKeyboardButton
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(btns[i][0], btns[i][1]))
		if i+1 < len(btns) {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(btns[i+1][0], btns[i+1][1]))
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(row...))
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg := tgbotapi.NewMessage(chatID, msgTxt)
	msg.ReplyMarkup = keyboard
	_, err := bs.bot.Send(msg)
	if err != nil {
		log.Printf("error sending pair buttons message: %v", err.Error())
	}
	return err
}
func (bs *BotServer) MediaView(chatID int64) error {
	cfg := tgbotapi.NewMediaGroup(chatID, []interface{}{
		tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL("https://picsum.photos/400")),
		tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL("https://picsum.photos/536/354")),
	})
	_, err := bs.bot.SendMediaGroup(cfg)
	if err != nil {
		log.Printf("error sending media group message: %v", err)
		return err
	}

	return nil
}
