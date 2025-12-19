package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort  string
	DbDriver string
	DbName   string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	return &Config{
		AppPort:  getEnv("APP_PORT", "8080"),
		DbDriver: getEnv("DB_DRIVER", "sqlite"),
		DbName:   getEnv("DB_NAME", "yk.db"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
