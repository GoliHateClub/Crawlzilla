package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
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

// GetBoolean Return type safe boolean config value
func GetBoolean(name string) bool {
	return os.Getenv(name) == "true"
}
