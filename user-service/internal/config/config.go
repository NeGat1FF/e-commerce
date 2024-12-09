package config

import "os"

type Config struct {
	Addr             string
	Port             string
	JWTSecret        string
	DATABASE_URL     string
	NOTIFICATION_URL string
}

func LoadConfig() *Config {
	return &Config{
		JWTSecret:        os.Getenv("JWT_SECRET"),
		DATABASE_URL:     os.Getenv("DATABASE_URL"),
		NOTIFICATION_URL: os.Getenv("NOTIFICATION_SERVICE_URL"),
	}
}
