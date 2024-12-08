package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all the configuration values for the application
type Config struct {
	JWTSecret string
}

// LoadConfig reads configuration from config file and environment variables
func LoadConfig() *Config {
	godotenv.Load()

	cfg := Config{
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
	return &cfg
}
