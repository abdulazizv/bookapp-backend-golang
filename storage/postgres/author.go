package postgres

import (
	"log"
	"time"

	"gitlab.com/bookapp/api/models"
)

const (
	layout_year = "2006"
	layout_date = "2006-01-02"
)

var (
	birthYear time.Time
	diedYear  time.Time
)

func (r storagePg) AuthorCreate(req *models.AuthorReq) (*models.AuthorRes, error) {
	res := models.AuthorRes{}
	query := `
INSERT INTO 
    author(
		first_name, 
		last_name, 
		middle_name, 
		birth_date, 
		country, 
		avatar_url, 
		about_text, 
		creativity, 
		died_year) 
VALUES 
    ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
RETURNING 
	id, 
	first_name, 
	last_name, 
	middle_name, 
	birth_date, 
	country, 
	avatar_url, 
	about_text, 
	creativity, 
	died_year, 
	created_at, 
	updated_at`

	err := r.db.QueryRow(query,
		req.FirstName,
		req.LastName,
		req.MiddleName,
		req.BirthDay,
		req.Country,
		req.AvatarUrl,
		req.AboutText,
		req.Creativity,
		req.DiedYear).Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.MiddleName,
		&birthYear,
		&res.Country,
		&res.AvatarUrl,
		&res.AboutText,
		&res.Creativity,
		&diedYear,
		&CreatedAt,
		&UpdatedAt)
	if err != nil {
		return nil, err
	}
	res.CreatedAt = CreatedAt.Format(Layout)
	res.UpdatedAt = UpdatedAt.Format(Layout)
	res.BirthDay = birthYear.Format(layout_date)
	res.DiedYear = diedYear.Format(layout_date)
	return &res, nil
}

func (r storagePg) AuthorGet(id int) (*models.AuthorRes, error) {
	res := models.AuthorRes{}
	query := `
SELECT 
	id, 
	first_name, 
	last_name, 
	middle_name, 
	birth_date, 
	country, 
	avatar_url, 
	about_text, 
	creativity, 
	died_year, 
	created_at, 
	updated_at
FROM 
	author WHERE id=$1`
	err := r.db.QueryRow(query, id).Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.MiddleName,
		&birthYear,
		&res.Country,
		&res.AvatarUrl,
		&res.AboutText,
		&res.Creativity,
		&diedYear,
		&CreatedAt,
		&UpdatedAt)
	if err != nil {
		log.Println("error: ", err.Error())
		return nil, err
	}
	res.CreatedAt = CreatedAt.Format(Layout)
	res.UpdatedAt = UpdatedAt.Format(Layout)
	res.BirthDay = birthYear.Format(layout_date)
	res.DiedYear = diedYear.Format(layout_date)
	res.Books, err = r.BookGetByAuthorId(id)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r storagePg) AuthorUpdate(req *models.AuthorUpdateReq) (*models.AuthorRes, error) {
	res := models.AuthorRes{}
	query := `
UPDATE 
	author 
SET 
	first_name=$1, 
	last_name=$2, 
	middle_name=$3, 
	birth_date=$4, 
	country=$5, 
	avatar_url=$6,
	about_text=$7, 
	creativity=$8, 
	died_year=$9, 
	updated_at=now()
WHERE 
	id=$10
RETURNING 
	id, 
	first_name, 
	last_name, 
	middle_name, 
	birth_date, 
	country, 
	avatar_url,
	about_text, 
	creativity, 
	died_year, 
	created_at, 
	updated_at`
	err := r.db.QueryRow(query,
		req.FirstName,
		req.LastName,
		req.MiddleName,
		req.BirthDay,
		req.Country,
		req.AvatarUrl,
		req.AboutText,
		req.Creativity,
		req.DiedYear,
		req.ID).Scan(
		&res.ID,
		&res.FirstName,
		&res.LastName,
		&res.MiddleName,
		&birthYear,
		&res.Country,
		&res.AvatarUrl,
		&res.AboutText,
		&res.Creativity,
		&diedYear,
		&CreatedAt,
		&UpdatedAt)
	if err != nil {
		log.Println("error: ", err.Error())
		return nil, err
	}
	res.CreatedAt = CreatedAt.Format(Layout)
	res.UpdatedAt = UpdatedAt.Format(Layout)
	res.BirthDay = birthYear.Format(layout_date)
	res.DiedYear = diedYear.Format(layout_date)
	return &res, nil
}

func (r storagePg) AuthorDelete(id int) error {
	_, err := r.db.Exec(`DELETE FROM views 
	WHERE book_id IN (SELECT id FROM books WHERE author_id = $1);`, id)
	if err != nil {
		log.Println("error delete subcategory: ", err.Error())
		return err
	}

	_, err = r.db.Exec(`DELETE FROM book_likes 
	WHERE book_id IN (SELECT id FROM books WHERE author_id = $1);`, id)
	if err != nil {
		log.Println("error delete subcategory: ", err.Error())
		return err
	}

	_, err = r.db.Exec(`DELETE FROM comments 
	WHERE book_id IN (SELECT id FROM books WHERE author_id = $1);`, id)
	if err != nil {
		log.Println("error delete subcategory: ", err.Error())
		return err
	}

	_, err = r.db.Exec(`DELETE FROM books WHERE author_id=$1`, id)
	if err != nil {
		log.Println("error: ", err.Error())
		return err
	}
	_, err = r.db.Exec(`DELETE FROM author WHERE id=$1`, id)
	if err != nil {
		log.Println("error: ", err.Error())
		return err
	}

	return nil
}

func (r storagePg) AuthorGetList(limit, page int, search string) (*models.Authors, error) {
	res := models.Authors{}
	offset := (page - 1) * limit
	query := `
SELECT 
	id, 
	first_name, 
	last_name, 
	birth_date, 
	died_year,
	avatar_url,
	COUNT(*) OVER () AS total_count
FROM 
	author WHERE (first_name ILIKE $1 OR last_name ILIKE $1) LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		log.Println("error: ", err.Error())
		return nil, err
	}
	for rows.Next() {
		temp := models.AuthorForList{}
		err = rows.Scan(&temp.ID, &temp.FirstName, &temp.LastName,
			&birthYear, &diedYear, &temp.AvatarUrl, &res.Meta.TotalCount)
		if err != nil {
			log.Println("error: ", err.Error())
			return nil, err
		}
		temp.BirthDay = birthYear.Format(layout_year)
		temp.DiedYear = diedYear.Format(layout_year)
		temp.BookCount, err = r.GetBookCount(temp.ID)
		if err != nil {
			return nil, err
		}
		res.Authors = append(res.Authors, &temp)
	}

	res.Meta.PageCount = (res.Meta.TotalCount + limit - 1) / limit
	res.Meta.CurrentPage = page
	res.Meta.PerPage = limit

	return &res, nil
}

func (r storagePg) GetBookCount(id int) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT count(*) FROM books WHERE author_id=$1`, id).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
