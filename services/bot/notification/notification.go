package notification

import (
	"Crawlzilla/database"
	"Crawlzilla/services/users"
	"errors"
	"fmt"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

// NotifySuperAdmin sends a message to the super admin
func NotifySuperAdmin(bot *tgbotapi.BotAPI, message string) error {
	// Retrieve SUPER_ADMIN_ID from environment variables
	superAdminIDStr := os.Getenv("SUPER_ADMIN_ID")
	if superAdminIDStr == "" {
		return errors.New("SUPER_ADMIN_ID environment variable is not set")
	}

	// Convert SUPER_ADMIN_ID to int64
	superAdminID, err := strconv.ParseInt(superAdminIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse SUPER_ADMIN_ID: %w", err)
	}

	// Fetch super admin user details from the database
	superAdmin, err := users.GetUserByTelegramIDService(database.DB, superAdminID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("super admin not found in the database")
		}
		return fmt.Errorf("error fetching super admin: %w", err)
	}

	// Validate that the super admin has a valid ChatID
	if superAdmin.ChatID == 0 {
		return errors.New("super admin does not have a valid chat ID")
	}

	// Validate the message
	if message == "" {
		return errors.New("message cannot be empty")
	}

	// Construct the message
	msg := tgbotapi.NewMessage(superAdmin.ChatID, message)

	// Send the message
	_, err = bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send notification to super admin: %w", err)
	}

	return nil
}
