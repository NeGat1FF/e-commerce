package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all the configuration values for the application
type Config struct {
	LogLevel              string
	ServerPort            string
	DatabaseURL           string
	CacheURL              string
	MessageBrokerURL      string
	MessageBrokerExchange string
	GRPCPort              string
}

// LoadConfig reads configuration from config file and environment variables
func LoadConfig() *Config {
	godotenv.Load()

	cfg := Config{
		LogLevel:              os.Getenv("LOG_LEVEL"),
		ServerPort:            os.Getenv("SERVER_PORT"),
		DatabaseURL:           os.Getenv("DATABASE_URL"),
		CacheURL:              os.Getenv("CACHE_URL"),
		MessageBrokerURL:      os.Getenv("MESSAGE_BROKER_URL"),
		MessageBrokerExchange: os.Getenv("MESSAGE_BROKER_EXCHANGE"),
		GRPCPort:              os.Getenv("GRPC_PORT"),
	}
	return &cfg
}
