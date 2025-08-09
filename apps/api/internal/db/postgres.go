package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Traunin/stickerpack-editor/apps/api/internal/env"
	_ "github.com/lib/pq"
)

const (
	maxRetries              = 5
	retryWait               = time.Second
	insertStickerpackQuery  = `INSERT INTO stickerpacks (user_id, title, is_public) VALUES ($1, $2, $3)`
	publicStickerpacksQuery = `SELECT id, title FROM stickerpacks WHERE is_public = true OFFSET $1 LIMIT $2`
	countPacksQuery         = `SELECT COUNT(*) FROM stickerpacks WHERE is_public = true`
)

type Postgres struct {
	db *sql.DB
}

type PublicPack struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
}

func (s *PublicPack) ScanRow(rows *sql.Rows) error {
	return rows.Scan(&s.Id, &s.Title)
}

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

func (p Postgres) AddStickerpack(id int64, title string, isPublic bool) error {
	_, err := p.db.Exec(insertStickerpackQuery, id, title, isPublic)
	return err
}

func (p Postgres) PublicStickerpacks(page, pageSize int) ([]PublicPack, error) {
	offset := page * pageSize
	rows, err := p.db.Query(publicStickerpacksQuery, offset, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	stickerpacks := make([]PublicPack, 0, pageSize)

	for rows.Next() {
		var sp PublicPack
		err := sp.ScanRow(rows)
		if err != nil {
			log.Printf("failed to scan stickerpack: %v", err)
			continue
		}
		stickerpacks = append(stickerpacks, sp)
	}

	return stickerpacks, nil
}

func (p Postgres) PublicPacksCount() (int, error) {
	var count int
	err := p.db.QueryRow(countPacksQuery).Scan(&count)
	return count, err
}
