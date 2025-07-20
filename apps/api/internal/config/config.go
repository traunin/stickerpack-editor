package config

import (
	"log"
	"sync"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/db"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/env"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	TelegramToken string
	Port          string
	BotName       string
	DomainCORS    string
	DBConn        *db.Postgres
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
			Port:          env.Fallback("PORT", "8080"),
			DomainCORS:    env.Fallback("DOMAIN_CORS", "*"),
			TelegramToken: env.Must("TELEGRAM_TOKEN"),
			BotName:       env.Must("BOT_NAME"),
			DBConn:        db.NewPostgres(),
		}
	})

	return cfg
}
