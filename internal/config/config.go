package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TELEGRAM_APITOKEN string `env:"TELEGRAM_APITOKEN"`
	QDRANT_URL        string `env:"QDRANT_URL"`
	QDRANT_COLLECTION string `env:"QDRANT_COLLECTION"`
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("there was an error laoding the env file %v", err)
	}

	return &Config{
		TELEGRAM_APITOKEN: getEnv("TELEGRAM_APITOKEN", ""),
		QDRANT_URL:        getEnv("QDRANT_URL", ""),
		QDRANT_COLLECTION: getEnv("QDRANT_COLLECTION", ""),
	}
}

func getEnv(key, fallbackValue string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}
	return fallbackValue
}
