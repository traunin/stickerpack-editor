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
	telegramToken string
	port          string
	botName       string
	domainCORS    string
	secretKey     string
	dbConn        *db.Postgres
}

var (
	cfg  *Config
	once sync.Once
)

func (c *Config) TelegramToken() string { return c.telegramToken }
func (c *Config) Port() string          { return c.port }
func (c *Config) BotName() string       { return c.botName }
func (c *Config) DomainCORS() string    { return c.domainCORS }
func (c *Config) SecretKey() string     { return c.secretKey }
func (c *Config) DBConn() *db.Postgres  { return c.dbConn }

func Load() *Config {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Println("no local .env found")
		}

		cfg = &Config{
			port:          env.Fallback("PORT", "8080"),
			domainCORS:    env.Fallback("DOMAIN_CORS", "*"),
			telegramToken: env.Must("TELEGRAM_TOKEN"),
			botName:       env.Must("BOT_NAME"),
			secretKey:     env.Must("SECRET_KEY"),
			dbConn:        db.NewPostgres(),
		}
	})

	return cfg
}
