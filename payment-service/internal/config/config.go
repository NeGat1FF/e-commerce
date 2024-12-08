package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all the configuration values for the application
type Config struct {
	ServerPort          string
	DatabaseURL         string
	StripeSecret        string
	StripeWebhookSecret string
}

// LoadConfig reads configuration from config file and environment variables
func LoadConfig() *Config {
	godotenv.Load()

	cfg := Config{
		ServerPort:          os.Getenv("SERVER_PORT"),
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		StripeSecret:        os.Getenv("STRIPE_SECRET"),
		StripeWebhookSecret: os.Getenv("STRIPE_WEBHOOK_SECRET"),
	}
	return &cfg
}
