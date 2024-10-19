package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"gitlab.com/bookapp/api/models"
)

func (r storagePg) ViewCreate(req *models.View) error {
	_, err := r.db.Exec(`INSERT INTO views(book_id, user_agent, count) VALUES($1, $2, $3)`, req.BookId, req.UserAgent, req.Count)
	if err != nil {
		return err
	}
	return nil
}

func (r storagePg) CheckFieldView(req *models.CheckFieldViewReq) (*models.CheckFieldRes, error) {
	res := models.CheckFieldRes{Exists: false}
	query := fmt.Sprintf("SELECT 1 FROM views WHERE %s=$1 AND book_id=$2", req.Field)
	var temp = 0
	err := r.db.QueryRow(query, req.Value, req.BookId).Scan(&temp)
	if errors.Is(err, sql.ErrNoRows) {
		return &res, nil
	} else if err != nil {
		log.Println(err.Error())
		return &res, err
	}
	if temp == 1 {
		res.Exists = true
		return &res, nil
	}

	return &res, err
}
