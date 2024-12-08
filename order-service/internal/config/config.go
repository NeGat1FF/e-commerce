package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_URL        string
	JWTSecret     string
	PRICE_SERVICE string
}

var config *Config

func InitConfig() {
	godotenv.Load()

	config = &Config{
		DB_URL:        os.Getenv("DATABASE_URL"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		PRICE_SERVICE: os.Getenv("PRICE_SERVICE"),
	}
}

func GetConfig() *Config {
	return config
}
