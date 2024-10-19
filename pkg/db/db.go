package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"gitlab.com/bookapp/config"
)

func ConnectToDb(cfg config.Config) (*sql.DB, error) {
	psqlString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbName,
	)

	conn, err := sql.Open("postgres", psqlString)
	if err != nil {
		return nil, err
	}
	return conn, err
}
