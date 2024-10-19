package postgres

import (
	"gitlab.com/bookapp/api/models"
)

func (r storagePg) CommentCreate(req *models.CommentReq) (*models.CommentRes, error) {
	res := models.CommentRes{}
	query := `
	INSERT INTO 
		comments(
			user_id, 
			book_id, 
			text
			) 
	VALUES 
		($1, $2, $3) 
	RETURNING 
		id, 
		user_id, 
		book_id, 
		text, 
		created_at, 
		updated_at`
	err := r.db.QueryRow(query,
		req.UserId,
		req.BookId,
		req.Text).
		Scan(
			&res.Id,
			&res.UserId,
			&res.BookId,
			&res.Text,
			&CreatedAt,
			&UpdatedAt,
		)
	if err != nil {
		return nil, err
	}
	res.CreatedAt = CreatedAt.Format(Layout)
	res.UpdatedAt = UpdatedAt.Format(Layout)
	return &res, nil
}
func (r storagePg) CommentGet(id int) (*models.CommentRes, error) {
	res := models.CommentRes{}
	query := `
	SELECT 
		id, 
		user_id, 
		book_id, 
		text, 
		created_at, 
		updated_at 
	FROM 
		comments 
	WHERE 
		id=$1`
	err := r.db.QueryRow(query, id).Scan(
		&res.Id,
		&res.UserId,
		&res.BookId,
		&res.Text,
		&CreatedAt,
		&UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	res.CreatedAt = CreatedAt.Format(Layout)
	res.UpdatedAt = UpdatedAt.Format(Layout)
	return &res, nil
}
func (r storagePg) CommentUpdate(req *models.CommentUpdate) (*models.CommentRes, error) {
	res := models.CommentRes{}
	query := `
UPDATE 
	comments 
SET 
	text=$1, 
	updated_at=now() 
WHERE 
	id=$2
RETURNING 
	id, 
	user_id, 
	book_id, 
	text, 
	created_at, 
	updated_at`
	err := r.db.QueryRow(query, req.Text, req.Id).Scan(
		&res.Id,
		&res.UserId,
		&res.BookId,
		&res.Text,
		&CreatedAt,
		&UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	res.CreatedAt = CreatedAt.Format(Layout)
	res.UpdatedAt = UpdatedAt.Format(Layout)
	return &res, nil
}
func (r storagePg) CommentDelete(id int) error {
	_, err := r.db.Exec(`DELETE FROM comments WHERE id=$1`, id)
	if err != nil {
		return err
	}
	return nil
}
func (r storagePg) CommentGetBookID(id int) ([]*models.CommentResponse, error) {
	res := []*models.CommentResponse{}
	query := `
SELECT 
	c.id, 
	c.user_id, 
	c.book_id, 
	c.text, 
	c.created_at, 
	c.updated_at, 
	u.id, 
	u.full_name, 
	u.avatar_url
FROM 
	comments c 
JOIN 
	users u 
ON u.id=c.user_id 
	WHERE book_id=$1 
ORDER BY 
	c.created_at ASC`
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		temp := models.CommentResponse{}
		err = rows.Scan(
			&temp.Id,
			&temp.UserId,
			&temp.BookId,
			&temp.Text,
			&CreatedAt,
			&UpdatedAt,
			&temp.User.Id,
			&temp.User.FullName,
			&temp.User.AvatarUrl,
		)
		if err != nil {
			return nil, err
		}
		temp.CreatedAt = CreatedAt.Format(Layout)
		temp.UpdatedAt = UpdatedAt.Format(Layout)
		res = append(res, &temp)
	}
	return res, nil
}
