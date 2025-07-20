package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/env"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

const (
	maxRetries = 5
	retryWait  = time.Second
	insertStickerpackQuery = `INSERT INTO stickerpacks (user_id, title, public) VALUES ($1, $2, $3)`
)

func NewPostgres() *Postgres {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		env.Must("DB_HOST"),
		env.Must("DB_PORT"),
		env.Must("DB_USER"),
		env.Must("DB_PASSWORD"),
		env.Must("DB_NAME"),
	)

	for try := 1; try <= maxRetries; try++ {
		log.Printf("postgres connection attempt %d of %d", try, maxRetries)

		db, err := sql.Open("postgres", connStr)
		if err == nil && db.Ping() == nil {
			log.Println("postgres connection established")
			return &Postgres{db}
		}

		log.Printf("connection failed: %v", err)
		time.Sleep(retryWait)
	}

	log.Fatalf("failed to connect to postgres after %d attempts", maxRetries)
	return nil
}

func (p Postgres) AddStickerpack(id string, title string, isPublic bool) (error) {
	_, err := p.db.Exec(insertStickerpackQuery, id, title, isPublic)
	return err
}
