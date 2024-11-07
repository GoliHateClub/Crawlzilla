package routers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bs *BotServer) CommandHandler(message *tgbotapi.Message) {
	bs.SetUserState(message.Chat.ID, message.Command())
	switch message.Command() {
	case "start":
		if GetUserRole(message.From.ID) == AdminRole {
			bs.StartCommandHandlerAdmin(message.Chat.ID)
		} else {
			bs.StartCommandHandlerUser(message.From.ID)
		}
	case "help":
		bs.HelpCommandHandler(message.Chat.ID)
	case "echo":
		bs.TextView(message.Chat.ID, "Echo command received! Type something to see it echoed back.")
	case "admin":
		bs.AdminCommandHandler(message)
	case "set_filter":
		bs.SetUserState(message.Chat.ID, "wait_min_price")
		bs.TextView(message.Chat.ID, "Please enter the minimum price:")
	case "test":
		buttons := [][2]string{{"hello🖐🏻", "hello"}, {"bey👋🏻", "bey"}, {"hello2222🖐🏻", "hello2"}, {"bey222👋🏻", "bey2"}}
		msg := "یه پیام همینجوری صرفا جهت تست👀"
		bs.PairButtonView(buttons, msg, message.Chat.ID)
	case "media":
		bs.MediaView(message.Chat.ID)
	case "extension":
		err := bs.SendPhoto(message.Chat.ID, "https://picsum.photos/400", "Photo caption", "Button1", "command1", "Button2", "command2")
		if err != nil {
			log.Println("Error sending photo:", err)
		}
	default:
		bs.TextView(message.Chat.ID, "I'm here to respond to your commands! Try /help.")
	}
}

func (bs *BotServer) MessageHandler(message *tgbotapi.Message) {
	userState := bs.GetUserState(message.Chat.ID)
	switch userState.CurrentCommand {
	case "echo":
		bs.EchoMessageHandler(message, userState)
	case "wait_min_price":
		bs.MinPriceMessageHandler(message, userState)
	case "wait_max_price":
		bs.MaxPriceMessageHandler(message, userState)
	case "wait_city":
		bs.CityMessageHandler(message, userState)
	case "wait_neighborhood":
		bs.NeighborhoodMessageHandler(message, userState)
	case "wait_min_area":
		bs.MinAreaMessageHandler(message, userState)
	case "wait_max_area":
		bs.MaxAreaMessageHandler(message, userState)
	case "wait_category": //
		bs.CategoryMessageHandler(message, userState)
	case "wait_min_age":
		bs.MinAgeMessageHandler(message, userState)
	case "wait_max_age":
		bs.MaxAgeMessageHandler(message, userState)
	case "wait_min_floor":
		bs.MinFloorMessageHandler(message, userState)
	case "wait_max_floor":
		bs.MaxFloorMessageHandler(message, userState)
	case "wait_has_elevator":
		bs.ElevatorMessageHandler(message, userState)
	case "wait_has_storage":
		bs.StorageMessageHandler(message, userState)
	case "wait_min_date":
		bs.MinDateMessageHandler(message, userState)
	case "wait_max_date":
		bs.MaxDateMessageHandler(message, userState)
	default:
		bs.TextView(message.Chat.ID, "I'm here to respond to your commands! Try /help.")
	}
}

func (bs *BotServer) HandleCallbackQuery(callback *tgbotapi.CallbackQuery) {
	data := callback.Data
	var response string

	switch data {
	case "hello":
		response = "سلاااام کلاینت گل گلاب😘"
	case "bey":
		response = "خوش اومدی، زود به زود سربزن🙋🏻‍♂️"
	default:
		response = "Command not recognized!"
	}

	// Respond to the callback query (e.g., acknowledging button press)
	ack := tgbotapi.NewCallback(callback.ID, response)
	if _, err := bs.bot.Request(ack); err != nil {
		log.Printf("error acknowledging callback query: %v", err)
	}

	// Optionally, send a new message as a response
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, response)
	if _, err := bs.bot.Send(msg); err != nil {
		log.Printf("error sending response message: %v", err)
	}
}
