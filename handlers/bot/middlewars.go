package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// mock data ------------------------------------------------------------
const (
	SuperAdminRole = "superAdmin"
	AdminRole      = "admin"
	UserRole       = "user"
)

var UserRoles = map[int64]string{
	5437262970: AdminRole,
	9876543218: UserRole,
}

func SetUserRole(userID int64, role string) {
	UserRoles[userID] = role
}

func GetUserRole(userID int64) string {
	role, exists := UserRoles[userID]
	if !exists {
		return UserRole
	}
	return role
}

// end mock data ------------------------------------------------------------
type MiddlewareFunc func(message *tgbotapi.Message, next func())

func UseMiddleware(message *tgbotapi.Message, middleware []MiddlewareFunc, handler func()) {
	var exec func(index int)
	exec = func(index int) {
		if index < len(middleware) {
			middleware[index](message, func() { exec(index + 1) })
		} else {
			handler()
		}
	}
	exec(0)
}

func( bs *BotServer) RoleMiddleware(requiredRole string) MiddlewareFunc {
	return func(message *tgbotapi.Message, next func()) {
		userID := message.From.ID
		userRole := GetUserRole(userID)
		if userRole != requiredRole {
			msg := tgbotapi.NewMessage(message.Chat.ID, "You don't have permission to perform this action.")
			bs.bot.Send(msg)
			
			return
		}
		next()
	}
}
