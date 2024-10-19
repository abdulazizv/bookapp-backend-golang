package postgres

import (
	"database/sql"
	"time"
)

const (
	Layout = "2006-01-02 15:04:05"
)

var (
	CreatedAt, UpdatedAt time.Time
)

type storagePg struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *storagePg {
	return &storagePg{
		db: db,
	}
}
