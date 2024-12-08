package config

import "os"

type Config struct {
	Addr             string
	Port             string
	JWTSecret        string
	DB_URL           string
	NOTIFICATION_URL string
}

func LoadConfig() *Config {
	return &Config{
		JWTSecret:        os.Getenv("JWT_SECRET"),
		DB_URL:           os.Getenv("DB_URL"),
		NOTIFICATION_URL: os.Getenv("NOTIFICATION_SERVICE_URL"),
	}
}
