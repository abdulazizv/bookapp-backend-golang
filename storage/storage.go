package storage

import (
	"database/sql"

	"gitlab.com/bookapp/storage/postgres"
	"gitlab.com/bookapp/storage/repo"
)

type StorageI interface {
	BookApp() repo.BookappService
}

type StoragePg struct {
	Db       *sql.DB
	bookrepo repo.BookappService
}

func NewStoragePg(db *sql.DB) *StoragePg {
	return &StoragePg{
		Db:       db,
		bookrepo: postgres.NewStorage(db),
	}
}

func (s StoragePg) BookApp() repo.BookappService {
	return s.bookrepo
}
