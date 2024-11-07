package bot

import (
	"bytes"
	"github.com/GoliHateClub/Crawlzilla/config"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (bs *BotServer) SendPhoto(chatID int64, photoURL, caption, button1Text, button1Data, button2Text, button2Data string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	url := "https://api.telegram.org/bot" + cfg.BotToken + "/sendPhoto"

	inlineKeyboard := [][]map[string]interface{}{
		{
			{
				"text":          button1Text,
				"callback_data": button1Data,
			},
			{
				"text":          button2Text,
				"callback_data": button2Data,
			},
		},
	}

	requestBody := map[string]interface{}{
		"chat_id":      chatID,
		"photo":        photoURL,
		"caption":      caption,
		"reply_markup": map[string]interface{}{"inline_keyboard": inlineKeyboard},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return retry(func() error {
		resp, err := bs.client.Do(req)
		if err != nil {
			return fmt.Errorf("request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("non-200 status: %s", resp.Status)
		}
		return nil
	}, 3, 2*time.Second)
}

// retry function attempts an operation with exponential backoff
func retry(operation func() error, attempts int, backoff time.Duration) error {
	for i := 0; i < attempts; i++ {
		err := operation()
		if err == nil {
			return nil
		}
		fmt.Printf("Attempt %d failed: %v. Retrying in %v...\n", i+1, err, backoff)
		time.Sleep(backoff)
		backoff *= 2
	}
	return fmt.Errorf("operation failed after %d attempts", attempts)
}
