package config

import (
	"log"
	"strconv"
	"sync"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/db"
	"github.com/Traunin/stickerpack-editor/apps/api/internal/env"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	telegramToken   string
	port            string
	botName         string
	domain          string
	secretKey       string
	downloadRetries int
	dbConn          *db.Postgres
}

var (
	cfg  *Config
	once sync.Once
)

func (c *Config) TelegramToken() string { return c.telegramToken }
func (c *Config) Port() string          { return c.port }
func (c *Config) BotName() string       { return c.botName }
func (c *Config) Domain() string        { return c.domain }
func (c *Config) SecretKey() string     { return c.secretKey }
func (c *Config) DownloadRetries() int  { return c.downloadRetries }
func (c *Config) DBConn() *db.Postgres  { return c.dbConn }

func Load() *Config {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Println("no local .env found")
		}

		secretKey := env.Must("SECRET_KEY")
		if len(secretKey) < 32 {
			log.Fatalln("SECRET_KEY must be >= 32 characters long")
		}
		downloadRetries, err := strconv.Atoi(env.Must("DOWNLOAD_RETRIES"))
		if err != nil {
			log.Fatalln("DOWNLOAD_RETRIES is not a number")
		}

		cfg = &Config{
			port:            env.Fallback("PORT", "8080"),
			domain:          env.Must("DOMAIN"),
			telegramToken:   env.Must("TELEGRAM_TOKEN"),
			botName:         env.Must("BOT_NAME"),
			secretKey:       secretKey,
			downloadRetries: downloadRetries,
			dbConn:          db.NewPostgres(),
		}
	})

	return cfg
}
