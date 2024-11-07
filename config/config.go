package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type RobotConfig struct {
	BotToken       string
	WebhookURL     string
	SocksProxyAddr string
	ListenAddr     string
}

func LoadConfig() (*RobotConfig, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	config := &RobotConfig{
		BotToken:       os.Getenv("BOT_TOKEN"),
		WebhookURL:     os.Getenv("WEBHOOK_URL"),
		SocksProxyAddr: os.Getenv("SOCKS_PROXY_ADDR"),
		ListenAddr:     os.Getenv("LISTEN_ADDR"),
	}

	return config, nil
}
