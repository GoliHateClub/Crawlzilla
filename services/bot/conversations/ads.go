package conversations

import (
	"Crawlzilla/services/bot/constants"
	"Crawlzilla/services/cache"
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"regexp"
	"strconv"

	cfg "Crawlzilla/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func AddAdConversation(ctx context.Context, state cache.State, update tgbotapi.Update) {
	bot := ctx.Value("bot").(*tgbotapi.BotAPI)
	states := ctx.Value("state").(*cache.UserState)
	configLogger := ctx.Value("configLogger").(cfg.ConfigLoggerType)
	botLogger, _ := configLogger("bot")

	switch state.Stage {
	case "init":
		msg := tgbotapi.NewMessage(state.ChatId, constants.NewAdText)
		msg.ReplyMarkup = nil
		bot.Send(msg)

		msg = tgbotapi.NewMessage(state.ChatId, constants.NewAdExampleText)
		bot.Send(msg)

		err := states.SetUserState(ctx, state.ChatId, cache.State{
			ChatId:       state.ChatId,
			UserId:       state.UserId,
			Stage:        "check_data_get_dec",
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
	case "check_data_get_dec":
		// Define a regex pattern to extract each field
		pattern := `(?m)- نام\s*(.+)\n- آدرس لوکیشن در بلد\s*(.+)\n- آدرس URL محصول در سایت\s*(.+)\n- شهر\s*(.+)\n- محله\s*(.+)\n- منبع آگهی\s*(.+)\n- نوع آگهی\s*(.+)\n- نوع ملک\s*(.+)\n- متراژ\s*(\d+)\n- قیمت\(به ریال\)\s*(\d+)\n- اجاره\(ریال\)\s*(\d+)\n- تعداد اتاق\s*(\d+)\n- طبقه واحد\s*(\d+)\n- تعداد طبقات ملک\s*(\d+)\n- شماره تماس\s*(\d+)\n- آسانسور دارد؟\s*(بله|خیر)\n- انباری دارد؟\s*(بله|خیر)\n- پارکینگ دارد؟\s*(بله|خیر)\n- بالکن دارد؟\s*(بله|خیر)`

		re := regexp.MustCompile(pattern)

		// Match the input with the pattern
		matches := re.FindStringSubmatch(update.Message.Text)

		if matches == nil {
			fmt.Println("No match found or input is invalid!")
			return
		}

		fmt.Printf("%v\n", matches)
	}
}
