package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	Port          string
	BotName       string
	DomainCORS    string
}

var (
	cfg  *Config
	once sync.Once
)

func Load() *Config {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Println("no local .env found")
		}

		cfg = &Config{
			Port:          fallbackEnv("PORT", "8080"),
			DomainCORS:    fallbackEnv("DOMAIN_CORS", "*"),
			TelegramToken: mustEnv("TELEGRAM_TOKEN"),
			BotName:       mustEnv("BOT_NAME"),
		}
	})

	return cfg
}

func mustEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("missing required %s\n", key)
	}

	return val
}

func fallbackEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return fallback
}
