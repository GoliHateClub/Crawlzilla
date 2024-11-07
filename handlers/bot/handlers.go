package bot

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bs *BotServer) StartCommandHandlerAdmin(chatID int64) {
	welcomeText := `‍‍سلااااام سوپرادمین عزیز😍`
	bs.TextView(chatID, welcomeText)
}
func (bs *BotServer) StartCommandHandlerUser(chatID int64) {
	welcomeText := `سلام! 👋 به ربات ادکرالر خوش اومدی! 🎉

برای شروع کار، می‌تونی از دستورات زیر استفاده کنی:

🔹 /start - آغاز به کار با ربات
🔹 /help - دریافت راهنمایی در مورد دستورات
🔹 /about - آشنایی بیشتر با ویژگی‌های ربات
🔹 /settings - تنظیمات ربات رو مدیریت کن`
	bs.TextView(chatID, welcomeText)
}
func (bs *BotServer) HelpCommandHandler(chatID int64) {
	bs.TextView(chatID, `Available commands:
	 /start,
	 /help,
	 /admin,
	 /set_filter,
	 /test,
	 /media,
	 /extension`)

}
func (bs *BotServer) AdminCommandHandler(message *tgbotapi.Message) {
	UseMiddleware(message, [] MiddlewareFunc{bs.RoleMiddleware( AdminRole)}, func() {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Hello, Admin! You have full access.")
		bs.bot.Send(msg)
	})
}

func (bs *BotServer) EchoMessageHandler(message *tgbotapi.Message, userState *UserState) {
	bs.TextView(message.Chat.ID, "Echo: "+message.Text)
	bs.mu.Lock()
	userState.Data["lastEcho"] = message.Text
	bs.mu.Unlock()
}
func (bs *BotServer) MinPriceMessageHandler(message *tgbotapi.Message, userState *UserState) {
	if parsed, err := strconv.Atoi(message.Text); err == nil {
		bs.SetMinPrice(message.Chat.ID, parsed)
		bs.mu.Lock()
		userState.CurrentCommand = "wait_max_price"
		bs.mu.Unlock()
		nextCommand := bs.String(message.Chat.ID) + "\n" + "Please enter the maximum price:"
		bs.TextView(message.Chat.ID, nextCommand)
	} else {
		bs.TextView(message.Chat.ID, "Invalid input. Please enter a valid minimum price.")
	}
}
func (bs *BotServer) MaxPriceMessageHandler(message *tgbotapi.Message, userState *UserState) {
	if parsed, err := strconv.Atoi(message.Text); err == nil {
		bs.SetMaxPrice(message.Chat.ID, parsed)
		bs.mu.Lock()
		userState.CurrentCommand = "wait_city"
		bs.mu.Unlock()
		nextCommand := bs.String(message.Chat.ID) + "\n" + "Please enter the city:"
		bs.TextView(message.Chat.ID, nextCommand)
	} else {
		bs.TextView(message.Chat.ID, "Invalid input. Please enter a valid maximum price.")
	}
}
func (bs *BotServer) CityMessageHandler(message *tgbotapi.Message, userState *UserState) {
	bs.SetCity(message.Chat.ID, message.Text)
	bs.mu.Lock()
	userState.CurrentCommand = "wait_neighborhood"
	bs.mu.Unlock()
	nextCommand := bs.String(message.Chat.ID) + "\n" + "Please enter the neighborhood:"
	bs.TextView(message.Chat.ID, nextCommand)
}
func (bs *BotServer) MinAreaMessageHandler(message *tgbotapi.Message, userState *UserState) {
	if parsed, err := strconv.Atoi(message.Text); err == nil {
		bs.SetMinArea(message.Chat.ID, parsed)
		bs.mu.Lock()
		userState.CurrentCommand = "wait_max_area"
		bs.mu.Unlock()
		nextCommand := bs.String(message.Chat.ID) + "\n" + "Please enter the maximum area:"
		bs.TextView(message.Chat.ID, nextCommand)
	} else {
		bs.TextView(message.Chat.ID, "Invalid input. Please enter a valid minimum area.")
	}
}
func (bs *BotServer) MaxAreaMessageHandler(message *tgbotapi.Message, userState *UserState) {
	if parsed, err := strconv.Atoi(message.Text); err == nil {
		bs.SetMaxArea(message.Chat.ID, parsed)
		bs.mu.Lock()
		userState.CurrentCommand = "wait_category"
		bs.mu.Unlock()
		nextCommand := bs.String(message.Chat.ID) + "\n" + "Please enter the category:"
		bs.TextView(message.Chat.ID, nextCommand)
	} else {
		bs.TextView(message.Chat.ID, "Invalid input. Please enter a valid maximum area.")
	}
}

func (bs *BotServer) NeighborhoodMessageHandler(message *tgbotapi.Message, userState *UserState) {
	bs.SetNeighborhood(message.Chat.ID, message.Text)
	bs.mu.Lock()
	userState.CurrentCommand = "wait_min_area"
	bs.mu.Unlock()
	nextCommand := bs.String(message.Chat.ID) + "\n" + "Please enter the minimum area:"
	bs.TextView(message.Chat.ID, nextCommand)
}
func (bs *BotServer) CategoryMessageHandler(message *tgbotapi.Message, userState *UserState) {
	bs.SetCategory(message.Chat.ID, message.Text)
	bs.mu.Lock()
	userState.CurrentCommand = "wait_min_age"
	bs.mu.Unlock()
	nextCommand := bs.String(message.Chat.ID) + "\n" + "Please enter the min Age:"
	bs.TextView(message.Chat.ID, nextCommand)
}
func (bs *BotServer) MinAgeMessageHandler(message *tgbotapi.Message, userState *UserState) {
	if parsed, err := strconv.Atoi(message.Text); err == nil {
		bs.SetMinAge(message.Chat.ID, parsed)
		bs.mu.Lock()
		userState.CurrentCommand = "wait_max_age"
		bs.mu.Unlock()
		nextCommand := bs.String(message.Chat.ID) + "\n" + "Please enter the max Age:"
		bs.TextView(message.Chat.ID, nextCommand)
	} else {
		bs.TextView(message.Chat.ID, "Invalid input. Please enter a valid minimum age.")
	}
}
func (bs *BotServer) MaxAgeMessageHandler(message *tgbotapi.Message, userState *UserState) {
	if parsed, err := strconv.Atoi(message.Text); err == nil {
		bs.SetMaxAge(message.Chat.ID, parsed)
		bs.mu.Lock()
		userState.CurrentCommand = "wait_min_floor"
		bs.mu.Unlock()
		nextCommand := bs.String(message.Chat.ID) + "\n" + "Please enter the minimum floor:"
		bs.TextView(message.Chat.ID, nextCommand)
	} else {
		bs.TextView(message.Chat.ID, "Invalid input. Please enter a valid maximum age.")
	}
}
func (bs *BotServer) MinFloorMessageHandler(message *tgbotapi.Message, userState *UserState) {
	if parsed, err := strconv.Atoi(message.Text); err == nil {
		bs.SetMinFloor(message.Chat.ID, parsed)
		bs.mu.Lock()
		userState.CurrentCommand = "wait_max_floor"
		bs.mu.Unlock()
		nextCommand := bs.String(message.Chat.ID) + "\n" + "Please enter the maximum floor:"
		bs.TextView(message.Chat.ID, nextCommand)
	} else {
		bs.TextView(message.Chat.ID, "Invalid input. Please enter a valid maximum floor.")
	}
}
func (bs *BotServer) MaxFloorMessageHandler(message *tgbotapi.Message, userState *UserState) {
	if parsed, err := strconv.Atoi(message.Text); err == nil {
		bs.SetMaxFloor(message.Chat.ID, parsed)
		bs.mu.Lock()
		userState.CurrentCommand = "wait_has_elevator"
		bs.mu.Unlock()
		nextCommand := bs.String(message.Chat.ID) + "\n" + "Elevator filter: Please enter yes or no:"
		bs.TextView(message.Chat.ID, nextCommand)
	} else {
		bs.TextView(message.Chat.ID, "Invalid input. Please enter a valid minimum floor.")
	}
}
func (bs *BotServer) ElevatorMessageHandler(message *tgbotapi.Message, userState *UserState) {
	ans := strings.ToLower(message.Text)
	if ans == "yes" {
		bs.SetHasElevator(message.Chat.ID, true)
		bs.mu.Lock()
		userState.CurrentCommand = "wait_has_storage"
		bs.mu.Unlock()
		nextCommand := bs.String(message.Chat.ID) + "\n" + "Storage filter: Please enter yes or no:"
		bs.TextView(message.Chat.ID, nextCommand)
	} else if ans == "no" {
		bs.SetHasElevator(message.Chat.ID, false)
		bs.mu.Lock()
		userState.CurrentCommand = "wait_has_storage"
		bs.mu.Unlock()
		nextCommand := bs.String(message.Chat.ID) + "\n" + "Storage filter: Please enter yes or no:"
		bs.TextView(message.Chat.ID, nextCommand)
	} else {
		bs.TextView(message.Chat.ID, "Invalid input. Please enter either yes or no.")
	}
}
func (bs *BotServer) StorageMessageHandler(message *tgbotapi.Message, userState *UserState) {
	ans := strings.ToLower(message.Text)
	if ans == "yes" {
		bs.SetHasStorage(message.Chat.ID, true)
		bs.mu.Lock()
		userState.CurrentCommand = "wait_min_date"
		bs.mu.Unlock()
		nextCommand := bs.String(message.Chat.ID) + "\n" + "Please enter the start date with format: 2024/11/25"
		bs.TextView(message.Chat.ID, nextCommand)
	} else if ans == "no" {
		bs.SetHasStorage(message.Chat.ID, false)
		bs.mu.Lock()
		userState.CurrentCommand = "wait_min_date"
		bs.mu.Unlock()
		nextCommand := bs.String(message.Chat.ID) + "\n" + "Please enter the start date with format: 2024/11/25"
		bs.TextView(message.Chat.ID, nextCommand)
	} else {
		bs.TextView(message.Chat.ID, "Invalid input. Please enter either yes or no.")
	}
}
func (bs *BotServer) MinDateMessageHandler(message *tgbotapi.Message, userState *UserState) {
	bs.SetMinDate(message.Chat.ID, message.Text)
	bs.mu.Lock()
	userState.CurrentCommand = "wait_max_date"
	bs.mu.Unlock()
	nextCommand := bs.String(message.Chat.ID) + "\n" + "Please enter the start date with format: 2024/11/25"
	bs.TextView(message.Chat.ID, nextCommand)
}
func (bs *BotServer) MaxDateMessageHandler(message *tgbotapi.Message, userState *UserState) {
	bs.SetMaxDate(message.Chat.ID, message.Text)
	bs.mu.Lock()
	userState.CurrentCommand = ""
	bs.mu.Unlock()
	nextCommand := bs.String(message.Chat.ID) + "\n" + "Please enter the start date with format: 2024/11/25"
	bs.TextView(message.Chat.ID, nextCommand)
}
