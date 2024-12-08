package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all the configuration values for the application
type Config struct {
	ElasticURL       string
	ElasticUsername  string
	ElasticPassword  string
	MessageBrokerURL string
}

// LoadConfig reads configuration from config file and environment variables
func LoadConfig() *Config {
	godotenv.Load()

	cfg := Config{
		ElasticURL:       os.Getenv("ELASTICSEARCH_URL"),
		ElasticUsername:  os.Getenv("ELASTICSEARCH_USERNAME"),
		ElasticPassword:  os.Getenv("ELASTICSEARCH_PASSWORD"),
		MessageBrokerURL: os.Getenv("RABBITMQ_URL"),
	}
	return &cfg
}
