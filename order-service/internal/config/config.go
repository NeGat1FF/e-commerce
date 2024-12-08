package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SERVER_PORT   string
	DB_URL        string
	JWTSecret     string
	PRICE_SERVICE string
}

var config *Config

func InitConfig() {
	godotenv.Load()

	config = &Config{
		SERVER_PORT:   os.Getenv("SERVER_PORT"),
		DB_URL:        os.Getenv("DATABASE_URL"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		PRICE_SERVICE: os.Getenv("PRICE_SERVICE"),
	}
}

func GetConfig() *Config {
	return config
}
