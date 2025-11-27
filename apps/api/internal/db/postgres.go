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
	maxRetries      = 5
	retryWait       = time.Second
	insertPackQuery = `
	INSERT INTO stickerpacks (user_id, name, title, is_public, thumbnail_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, title, name, thumbnail_id`
	publicPacksQuery = `
	SELECT id, title, name, thumbnail_id FROM stickerpacks
	WHERE is_public = true
	ORDER BY id DESC
	OFFSET $1 LIMIT $2`
	countPublicPacksQuery = `
	SELECT COUNT(*) FROM stickerpacks
	WHERE is_public = true`
	userPacksQuery = `
	SELECT id, title, name, thumbnail_id FROM stickerpacks
	WHERE user_id = $1
	ORDER BY id DESC
	OFFSET $2 LIMIT $3`
	countUserPacksQuery = `
	SELECT COUNT(*) FROM stickerpacks
	WHERE user_id = $1`
	nameExistsQuery = `SELECT EXISTS(SELECT 1 FROM stickerpacks WHERE name=$1)`
	packOwnedQuery  = `
	SELECT EXISTS (SELECT 1 FROM stickerpacks WHERE name=$1 AND user_id=$2)`
	packPublicQuery  = `
	SELECT is_public FROM stickerpacks WHERE name=$1`
	deletePackQuery = `
	DELETE FROM stickerpacks WHERE name=$1 AND user_id=$2 RETURNING id`
)

type Postgres struct {
	db *sql.DB
}

type PackResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Name        string `json:"name"`
	ThumbnailID string `json:"thumbnail_id"`
}

type StoredPack struct {
	ID          int64
	UserID      int64
	Name        string
	Title       string
	IsPublic    bool
	ThumbnailID string
}

type Option func(*StoredPack)

func WithID(id int64) Option {
	return func(sp *StoredPack) {
		sp.ID = id
	}
}

func WithUserID(userID int64) Option {
	return func(sp *StoredPack) {
		sp.UserID = userID
	}
}

func WithName(name string) Option {
	return func(sp *StoredPack) {
		sp.Name = name
	}
}

func WithTitle(title string) Option {
	return func(sp *StoredPack) {
		sp.Title = title
	}
}

func WithPublic(isPublic bool) Option {
	return func(sp *StoredPack) {
		sp.IsPublic = isPublic
	}
}

func WithThumbnail(thumbnailID string) Option {
	return func(sp *StoredPack) {
		sp.ThumbnailID = thumbnailID
	}
}

func NewStoredPack(opts ...Option) *StoredPack {
	sp := &StoredPack{}
	for _, opt := range opts {
		opt(sp)
	}
	return sp
}

func (s *PackResponse) ScanRow(rows *sql.Rows) error {
	return rows.Scan(&s.ID, &s.Title, &s.Name, &s.ThumbnailID)
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

func (p *Postgres) AddStickerpack(pack *StoredPack) (*PackResponse, error) {
	var resp PackResponse
	err := p.db.QueryRow(
		insertPackQuery,
		pack.UserID,
		pack.Name,
		pack.Title,
		pack.IsPublic,
		pack.ThumbnailID,
	).Scan(&resp.ID, &resp.Title, &resp.Name, &resp.ThumbnailID)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (p Postgres) PublicStickerpacks(page, pageSize int) ([]PackResponse, error) {
	offset := page * pageSize
	rows, err := p.db.Query(publicPacksQuery, offset, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	stickerpacks := make([]PackResponse, 0, pageSize)

	for rows.Next() {
		var sp PackResponse
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
	err := p.db.QueryRow(countPublicPacksQuery).Scan(&count)
	return count, err
}

func (p Postgres) UserPacks(userID int64, page, pageSize int) ([]PackResponse, error) {
	offset := page * pageSize
	rows, err := p.db.Query(userPacksQuery, userID, offset, pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	stickerpacks := make([]PackResponse, 0, pageSize)

	for rows.Next() {
		var sp PackResponse
		err := sp.ScanRow(rows)
		if err != nil {
			log.Printf("failed to scan stickerpack: %v", err)
			continue
		}
		stickerpacks = append(stickerpacks, sp)
	}

	return stickerpacks, nil
}

func (p Postgres) UserPacksCount(userID int64) (int, error) {
	var count int
	err := p.db.QueryRow(countUserPacksQuery, userID).Scan(&count)
	return count, err
}

func (p *Postgres) NameExists(name string) (bool, error) {
	var exists bool
	err := p.db.QueryRow(nameExistsQuery, name).Scan(&exists)
	return exists, err
}

func (p *Postgres) IsPackOwner(name string, userID int64) (bool, error) {
	var owned bool
	err := p.db.QueryRow(packOwnedQuery, name, userID).Scan(&owned)
	return owned, err
}

func (p *Postgres) IsPackPublic(name string) (bool, error) {
	var public bool
	err := p.db.QueryRow(packPublicQuery, name).Scan(&public)
	return public, err
}


func (p *Postgres) DeletePack(name string, userID int64) error {
	var deletedID int64
	err := p.db.QueryRow(deletePackQuery, name, userID).Scan(&deletedID)
	return err
}
