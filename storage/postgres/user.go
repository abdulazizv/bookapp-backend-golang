package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"gitlab.com/bookapp/api/models"
)

func (r storagePg) UserCreate(req *models.UserReq) (*models.UserRes, error) {
	fmt.Println(req)
	res := models.UserRes{}
	query := `
INSERT INTO 
	users(
        role_id, 
		full_name, 
        avatar_url, 
        login, 
		password, 
		refresh_token
	)
VALUES 
	($1, $2, $3, $4, $5, $6)
RETURNING 
	id, 
	full_name, 
	avatar_url, 
	login
`
	err := r.db.QueryRow(query,
		req.RoleId,
		req.FullName,
		req.AvatarUrl,
		req.Login,
		req.Password,
		req.RefreshToken).Scan(
		&res.Id,
		&res.FullName,
		&res.AvatarUrl,
		&res.Login)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func (r storagePg) UserGet(id int) (*models.UserRes, error) {
	res := models.UserRes{}
	query := `
SELECT 
	id, 
	full_name, 
	avatar_url, 
	login
FROM 
	users WHERE id=$1
`
	err := r.db.QueryRow(query, id).Scan(
		&res.Id, &res.FullName, &res.AvatarUrl, &res.Login)
	if err != nil {
		return nil, err
	}
	res.Books, err = r.BookGetLiked(id)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func (r storagePg) UserUpdate(req *models.UserUpdateReq) (*models.UserRes, error) {
	res := models.UserRes{}
	query := `
UPDATE 
    users 
SET 
    full_name=$1, 
	avatar_url=$2, 
	login=$3, 
	password=$4 
	WHERE id=$5
RETURNING 
	id, 
	full_name,
	avatar_url, 
	login`
	err := r.db.QueryRow(query, req.FullName, req.AvatarUrl, req.Login, req.Password, req.Id).
		Scan(
			&res.Id, &res.FullName, &res.AvatarUrl, &res.Login)
	if err != nil {
		return nil, err
	}
	res.Books = []*models.BooksForList{}

	return &res, nil
}
func (r storagePg) UserDelete(id int) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return err
	}
	return nil
}
func (r storagePg) CheckField(req *models.CheckFieldReq) (*models.CheckFieldRes, error) {
	res := models.CheckFieldRes{Exists: false}
	query := fmt.Sprintf("SELECT 1 FROM users WHERE %s=$1", req.Field)
	var temp = 0
	err := r.db.QueryRow(query, req.Value).Scan(&temp)
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
func (r storagePg) UserGetLogin(login string) (*models.UserLoginRes, error) {
	res := models.UserLoginRes{}
	query := `SELECT id, full_name, avatar_url, login, password, role_id FROM users WHERE login=$1`
	err := r.db.QueryRow(query, login).Scan(
		&res.Id, &res.FullName, &res.AvatarUrl, &res.Login, &res.Password, &res.RoleId,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r storagePg) AdminGetList() ([]*models.UserResForComment, error) {
	res := []*models.UserResForComment{}
	query := `SELECT id, full_name, avatar_url FROM users WHERE role_id=2`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		temp := models.UserResForComment{}
		err = rows.Scan(&temp.Id, &temp.FullName, &temp.AvatarUrl)
		if err != nil {
			return nil, err
		}
		res = append(res, &temp)
	}
	return res, nil
}
