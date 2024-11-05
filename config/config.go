package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadConfig loads environment variables from the .env file
func LoadConfig() error {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
		return err
	}
	return nil
}
