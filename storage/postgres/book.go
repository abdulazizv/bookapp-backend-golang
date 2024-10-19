package postgres

import (
	"database/sql"
	"errors"
	"log"
	"strconv"

	"gitlab.com/bookapp/api/models"
)

func (r storagePg) BookCreate(req *models.BookReq) (*models.BookRes, error) {
	res := models.BookRes{}
	query := `
INSERT INTO books(
        title, description, image, download_url, audio_url, category_id, subcategory_id, author_id, book_type) 
VALUES 
	($1, $2, $3, $4, $5, $6, $7, $8, $9) 
RETURNING 
	id, title, description, image, download_url, audio_url, book_type, category_id, subcategory_id, author_id, created_at, updated_at`

	err := r.db.QueryRow(query, req.Title, req.Description, req.Image, req.DownloadUrl, req.AudioUrl, req.CategoryID, req.SubCategoryID, req.AuthorID, req.BookType).
		Scan(
			&res.ID, &res.Title, &res.Description, &res.Image, &res.DownloadUrl, &res.AudioUrl, &res.BookType, &res.CategoryID,
			&res.SubCategoryID, &res.AuthorID, &CreatedAt, &UpdatedAt)
	if err != nil {
		log.Println("error: ", err.Error())
		return &models.BookRes{}, err
	}
	res.CreatedAt = CreatedAt.Format(Layout)
	res.UpdatedAt = UpdatedAt.Format(Layout)

	return &res, err
}

func (r storagePg) BookGet(id int, userId int) (*models.BookResponse, error) {
	res := models.BookResponse{}
	query := `
	SELECT
        b.id, 
		b.title, 
		b.description, 
		b.image, 
		b.download_url, 
		b.audio_url, 
		b.book_type, 
		b.category_id, 
        b.subcategory_id, 
		b.author_id, 
		b.created_at, 
		b.updated_at, 
		a.first_name, 
		a.last_name
	FROM
	    books b 
	JOIN 
		author a ON a.id=b.author_id WHERE b.id=$1
`
	err := r.db.QueryRow(query, id).
		Scan(
			&res.ID, &res.Title,
			&res.Description,
			&res.Image,
			&res.DownloadUrl,
			&res.AudioUrl,
			&res.BookType,
			&res.CategoryID,
			&res.SubCategoryID,
			&res.AuthorID,
			&CreatedAt,
			&UpdatedAt,
			&res.AuthorFirstName,
			&res.AuthorLastName)
	if err != nil {
		log.Println("error: ", err.Error())
		return &models.BookResponse{}, err
	}
	res.CreatedAt = CreatedAt.Format(Layout)
	res.UpdatedAt = UpdatedAt.Format(Layout)

	res.LikeCount, err = r.BookLikeCount(id)
	if err != nil {
		return nil, err
	}
	res.ViewCount, err = r.BookViewCount(id)
	if err != nil {
		return nil, err
	}
	res.Comments, err = r.CommentGetBookID(id)
	if err != nil {
		return nil, err
	}
	res.CommentCount = len(res.Comments)
	res.SimilarBooks, err = r.BookGetCategoryId(res.CategoryID)
	if err != nil {
		return nil, err
	}

	res.IsLike, err = r.CheckBooklike(id, userId)
	if err != nil {
		return nil, err
	}
	return &res, err
}

func (r storagePg) CheckBooklike(bookId, userId int) (bool, error) {
	var temp = 0
	var isLike = false
	query := "SELECT 1 FROM book_likes WHERE book_id=$1 AND user_id=$2"
	err := r.db.QueryRow(query, bookId, userId).Scan(&temp)
	if errors.Is(err, sql.ErrNoRows) {
		return isLike, nil
	} else if err != nil {
		log.Println(err.Error())
		return isLike, err
	}
	if temp == 1 {
		isLike = true
		return isLike, nil
	}

	return isLike, nil
}

func (r storagePg) BookUpdate(req *models.BookUpdate) (*models.BookRes, error) {
	res := models.BookRes{}
	query := `
UPDATE 
	books  
SET 
	title=$1, 
	description=$2, 
	image=$3, 
	download_url=$4, 
	audio_url=$5, 
	book_type=$6, 
	updated_at=now()
WHERE 
	id=$7
RETURNING 
	id,title, 
	description, 
	image, 
	download_url, 
	audio_url, 
	book_type,
	category_id,
	subcategory_id, 
	author_id, 
	created_at, 
	updated_at`

	err := r.db.QueryRow(query,
		req.Title,
		req.Description,
		req.Image,
		req.DownloadUrl,
		req.AudioUrl,
		req.BookType,
		req.ID).
		Scan(
			&res.ID,
			&res.Title,
			&res.Description,
			&res.Image,
			&res.DownloadUrl,
			&res.AudioUrl,
			&res.BookType,
			&res.CategoryID,
			&res.SubCategoryID,
			&res.AuthorID,
			&CreatedAt,
			&UpdatedAt)
	if err != nil {
		log.Println("error: ", err.Error())
		return &models.BookRes{}, err
	}
	res.CreatedAt = CreatedAt.Format(Layout)
	res.UpdatedAt = UpdatedAt.Format(Layout)

	return &res, err
}

func (r storagePg) BookDelete(id int) error {
	_, err := r.db.Exec(`DELETE FROM views WHERE book_id=$1`, id)
	if err != nil {
		log.Println("Error: ", err.Error())
		return err
	}
	_, err = r.db.Exec(`DELETE FROM book_likes WHERE book_id=$1`, id)
	if err != nil {
		log.Println("Error: ", err.Error())
		return err
	}
	_, err = r.db.Exec(`DELETE FROM comments WHERE book_id=$1`, id)
	if err != nil {
		log.Println("Error: ", err.Error())
		return err
	}

	_, err = r.db.Exec(`DELETE FROM books WHERE id=$1`, id)
	if err != nil {
		log.Println("Error: ", err.Error())
		return err
	}
	return nil
}

func (r storagePg) BookGetList(limit, page int, search string) (*models.Books, error) {
	res := models.Books{}
	offset := (page - 1) * limit

	query := `
SELECT
	b.id, 
	b.title, 
	b.image, 
	b.book_type, 
	a.first_name, 
	a.last_name, 
	b.category_id, 
	b.subcategory_id, 
	b.author_id,
	COUNT(*) OVER () AS total_count
FROM 
	books b 
JOIN 
	author a 
ON 
	a.id=b.author_id 
WHERE  
	(b.title ILIKE $1 OR b.description ILIKE $1 OR a.first_name ILIKE $1 OR a.last_name ILIKE $1)
LIMIT $2 OFFSET $3
`
	rows, err := r.db.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		log.Println("error: ", err.Error())
		return &models.Books{}, err
	}
	for rows.Next() {
		tem := models.BooksForList{}
		err = rows.Scan(
			&tem.Id,
			&tem.Title,
			&tem.Image,
			&tem.BookType,
			&tem.AuthorFirstName,
			&tem.AuthorLastName,
			&tem.CategoryId,
			&tem.SubCategoryId,
			&tem.AuthorId,
			&res.Meta.TotalCount)
		if err != nil {
			log.Println("error: ", err.Error())
			return &models.Books{}, err
		}
		tem.LikeCount, err = r.BookLikeCount(tem.Id)
		if err != nil {
			return nil, err
		}
		tem.ViewCount, err = r.BookViewCount(tem.Id)
		if err != nil {
			return nil, err
		}
		res.BookList = append(res.BookList, &tem)
	}

	res.Meta.PageCount = (res.Meta.TotalCount + limit - 1) / limit
	res.Meta.CurrentPage = page
	res.Meta.PerPage = limit
	return &res, nil
}

func (r storagePg) BookGetSearch(search string) (*models.Books, error) {
	res := models.Books{}

	query := `
SELECT 
	b.id, 
	b.title, 
	b.image, 
	b.book_type, 
	a.first_name, 
	a.last_name, 
	b.category_id, 
	b.subcategory_id, 
	b.author_id
FROM 
	books b JOIN author a ON a.id=b.author_id
WHERE 
	title ILIKE $1 OR description ILIKE $2`

	rows, err := r.db.Query(query, "%"+search+"%", "%"+search+"%")
	if err != nil {
		log.Println("error: ", err.Error())
		return &models.Books{}, err
	}
	for rows.Next() {
		tem := models.BooksForList{}
		err = rows.Scan(
			&tem.Id, &tem.Title, &tem.Image, &tem.BookType,
			&tem.AuthorFirstName, &tem.AuthorLastName,
			&tem.CategoryId, &tem.SubCategoryId, &tem.AuthorId)
		if err != nil {
			log.Println("error: ", err.Error())
			return &models.Books{}, err
		}
		tem.LikeCount, err = r.BookLikeCount(tem.Id)
		if err != nil {
			return nil, err
		}
		tem.ViewCount, err = r.BookViewCount(tem.Id)
		if err != nil {
			return nil, err
		}

		res.BookList = append(res.BookList, &tem)
	}
	return &res, nil
}

func (r storagePg) BookGetSubCaID(id int, limit, page int) (*models.Books, error) {
	res := models.Books{}
	offset := (page - 1) * limit
	query := `
	SELECT 
		b.id, 
		b.title, 
		b.image, 
		a.first_name, 
		a.last_name, 
		b.category_id, 
		b.subcategory_id, 
		b.author_id,
		COUNT(*) OVER () AS total_count
	FROM 
		books b 
	JOIN 
		author a 
	ON 
		a.id=b.author_id 
	WHERE 
		b.subcategory_id=$1 OFFSET $2 LIMIT $3`
	rows, err := r.db.Query(query, id, offset, limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		temp := models.BooksForList{}
		err = rows.Scan(
			&temp.Id,
			&temp.Title,
			&temp.Image,
			&temp.AuthorFirstName,
			&temp.AuthorLastName,
			&temp.CategoryId,
			&temp.SubCategoryId,
			&temp.AuthorId,
			&res.Meta.TotalCount)
		if err != nil {
			return nil, err
		}
		temp.LikeCount, err = r.BookLikeCount(temp.Id)
		if err != nil {
			return nil, err
		}
		temp.ViewCount, err = r.BookViewCount(temp.Id)
		if err != nil {
			return nil, err
		}
		res.BookList = append(res.BookList, &temp)
	}

	res.Meta.PageCount = (res.Meta.TotalCount + limit - 1) / limit
	res.Meta.CurrentPage = page
	res.Meta.PerPage = limit

	return &res, nil
}

func (r storagePg) BookGetCatId(id int, limit, page int) (*models.Books, error) {
	offset := (page - 1) * limit
	res := models.Books{}
	query := `
	SELECT 
		b.id, 
		b.title, 
		b.image, 
		a.first_name, 
		a.last_name, 
		b.category_id, 
		b.subcategory_id, 
		b.author_id,
		COUNT(*) OVER () AS total_count
	FROM 
		books b JOIN author a ON a.id=b.author_id 
	WHERE 
		b.category_id=$1 OFFSET $2 LIMIT $3`
	rows, err := r.db.Query(query, id, offset, limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		temp := models.BooksForList{}
		err = rows.Scan(
			&temp.Id,
			&temp.Title,
			&temp.Image,
			&temp.AuthorFirstName,
			&temp.AuthorLastName,
			&temp.CategoryId,
			&temp.SubCategoryId,
			&temp.AuthorId,
			&res.Meta.TotalCount)
		if err != nil {
			return nil, err
		}
		temp.LikeCount, err = r.BookLikeCount(temp.Id)
		if err != nil {
			return nil, err
		}
		temp.ViewCount, err = r.BookViewCount(temp.Id)
		if err != nil {
			return nil, err
		}
		res.BookList = append(res.BookList, &temp)
	}

	res.Meta.PageCount = (res.Meta.TotalCount + limit - 1) / limit
	res.Meta.CurrentPage = page
	res.Meta.PerPage = limit

	return &res, nil
}

func (r storagePg) BookGetCategoryId(CategoryId int) ([]*models.BooksForList, error) {
	res := []*models.BooksForList{}

	query := `
	SELECT 
		b.id, 
		b.title, 
		b.image, 
		a.first_name, 
		a.last_name, 
		b.category_id, 
		b.subcategory_id, 
		b.author_id
	FROM 
		books b JOIN author a ON a.id=b.author_id 
	WHERE 
		b.category_id=$1`
	rows, err := r.db.Query(query, CategoryId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		temp := models.BooksForList{}
		err = rows.Scan(
			&temp.Id,
			&temp.Title,
			&temp.Image,
			&temp.AuthorFirstName,
			&temp.AuthorLastName,
			&temp.CategoryId,
			&temp.SubCategoryId,
			&temp.AuthorId)
		if err != nil {
			return nil, err
		}
		temp.LikeCount, err = r.BookLikeCount(temp.Id)
		if err != nil {
			return nil, err
		}
		temp.ViewCount, err = r.BookViewCount(temp.Id)
		if err != nil {
			return nil, err
		}
		res = append(res, &temp)
	}

	return res, nil
}

func (r storagePg) BookCreateLike(req *models.LikeReq) error {
	_, err := r.db.Exec(`INSERT INTO book_likes(user_id, book_id) VALUES($1, $2)`, req.UserId, req.BookId)
	if err != nil {
		return err
	}
	return nil
}

func (r storagePg) BookDeleteLike(req *models.LikeReq) error {
	_, err := r.db.Exec(`DELETE FROM book_likes WHERE user_id=$1 AND book_id=$2`, req.UserId, req.BookId)
	if err != nil {
		return err
	}
	return nil
}

func (r storagePg) BookLikeCount(bookId int) (int, error) {
	var count = 0
	query := `SELECT COALESCE(COUNT(*), 0) AS like_count FROM book_likes WHERE book_id=$1`
	err := r.db.QueryRow(query, bookId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r storagePg) BookViewCount(bookId int) (int, error) {
	var count = 0
	query := `SELECT COALESCE(COUNT(*), 0) AS view_count FROM views WHERE book_id=$1`
	err := r.db.QueryRow(query, bookId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r storagePg) BookGetLiked(id int) ([]*models.BooksForList, error) {
	res := []*models.BooksForList{}
	query := `
SELECT 
	b.id, 
	b.title, 
	b.image, 
	b.book_type, 
	b.author_id, 
	b.category_id, 
	b.subcategory_id, 
	a.first_name, 
	a.last_name
FROM 
	books b 
JOIN 
	author a ON a.id=b.author_id 
JOIN 
	book_likes bl ON b.id=bl.book_id WHERE bl.user_id=$1`

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		temp := models.BooksForList{}
		err = rows.Scan(&temp.Id, &temp.Title, &temp.Image,
			&temp.BookType, &temp.AuthorId, &temp.CategoryId,
			&temp.SubCategoryId, &temp.AuthorFirstName, &temp.AuthorLastName)
		if err != nil {
			return nil, err
		}
		temp.LikeCount, err = r.BookLikeCount(temp.Id)
		if err != nil {
			return nil, err
		}
		temp.ViewCount, err = r.BookViewCount(temp.Id)
		if err != nil {
			return nil, err
		}
		res = append(res, &temp)
	}
	return res, err
}

func (r storagePg) BookGetMoreRead(limit, page int) (*models.Books, error) {
	res := models.Books{}
	offset := (page - 1) * limit

	query := `
SELECT
	b.id, 
	b.title, 
	b.image, 
	b.book_type, 
	a.first_name, 
	a.last_name, 
	b.category_id, 
	b.subcategory_id, 
	b.author_id,
	COUNT(v.book_id) AS view_count,
	COUNT(*) OVER () AS total_count
FROM 
	books b 
JOIN 
	author a 
ON 
	a.id=b.author_id 
LEFT JOIN 
    views v ON b.id = v.book_id
GROUP BY 
    b.id, 
    a.first_name, 
    a.last_name
ORDER BY 
    view_count DESC
LIMIT 
	$1 OFFSET $2
`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		log.Println("error: ", err.Error())
		return &models.Books{}, err
	}
	for rows.Next() {
		tem := models.BooksForList{}
		err = rows.Scan(
			&tem.Id,
			&tem.Title,
			&tem.Image,
			&tem.BookType,
			&tem.AuthorFirstName,
			&tem.AuthorLastName,
			&tem.CategoryId,
			&tem.SubCategoryId,
			&tem.AuthorId,
			&tem.ViewCount,
			&res.Meta.TotalCount,
		)
		if err != nil {
			log.Println("error: ", err.Error())
			return &models.Books{}, err
		}
		tem.LikeCount, err = r.BookLikeCount(tem.Id)
		if err != nil {
			return nil, err
		}
		tem.ViewCount, err = r.BookViewCount(tem.Id)
		if err != nil {
			return nil, err
		}
		res.BookList = append(res.BookList, &tem)
	}
	res.Meta.PageCount = (res.Meta.TotalCount + limit - 1) / limit
	res.Meta.CurrentPage = page
	res.Meta.PerPage = limit

	return &res, nil
}

func (r storagePg) BookGetTops(limit, page int) (*models.Books, error) {
	res := models.Books{}
	offset := (page - 1) * limit

	query := `
SELECT
	b.id, 
	b.title, 
	b.image, 
	b.book_type, 
	a.first_name, 
	a.last_name, 
	b.category_id, 
	b.subcategory_id, 
	b.author_id,
	COUNT(*) OVER () AS total_count
FROM 
	books b 
JOIN 
	author a 
ON 
	a.id=b.author_id 
WHERE 
	b.book_type=$1 LIMIT $2 OFFSET $3
`
	rows, err := r.db.Query(query, "top", limit, offset)
	if err != nil {
		log.Println("error: ", err.Error())
		return &models.Books{}, err
	}
	for rows.Next() {
		tem := models.BooksForList{}
		err = rows.Scan(
			&tem.Id,
			&tem.Title,
			&tem.Image,
			&tem.BookType,
			&tem.AuthorFirstName,
			&tem.AuthorLastName,
			&tem.CategoryId,
			&tem.SubCategoryId,
			&tem.AuthorId,
			&res.Meta.TotalCount,
		)
		if err != nil {
			log.Println("error: ", err.Error())
			return &models.Books{}, err
		}
		tem.LikeCount, err = r.BookLikeCount(tem.Id)
		if err != nil {
			return nil, err
		}
		tem.ViewCount, err = r.BookViewCount(tem.Id)
		if err != nil {
			return nil, err
		}
		res.BookList = append(res.BookList, &tem)
	}
	res.Meta.PageCount = (res.Meta.TotalCount + limit - 1) / limit
	res.Meta.CurrentPage = page
	res.Meta.PerPage = limit

	return &res, nil
}

func (r storagePg) BookGetAudios(limit, page int, search string) (*models.BooksAudios, error) {
	res := models.BooksAudios{}
	offset := (page - 1) * limit

	query := `
SELECT
	b.id, 
	b.title, 
	b.image, 
	b.audio_url,
	b.book_type, 
	a.first_name, 
	a.last_name, 
	b.category_id, 
	b.subcategory_id, 
	b.author_id,
	COUNT(*) OVER () AS total_count
FROM 
	books b 
JOIN 
	author a 
ON 
	a.id=b.author_id 
WHERE 
	(b.audio_url IS NOT NULL AND b.audio_url<>'')
	AND 
	(b.title ILIKE $1 OR b.description ILIKE $1 OR a.first_name ILIKE $1 OR a.last_name ILIKE $1)
	 LIMIT $2 OFFSET $3
`
	rows, err := r.db.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		log.Println("error: ", err.Error())
		return &models.BooksAudios{}, err
	}
	for rows.Next() {
		tem := models.BooksAudioList{}
		err = rows.Scan(
			&tem.Id,
			&tem.Title,
			&tem.Image,
			&tem.AudioURL,
			&tem.BookType,
			&tem.AuthorFirstName,
			&tem.AuthorLastName,
			&tem.CategoryId,
			&tem.SubCategoryId,
			&tem.AuthorId,
			&res.Meta.TotalCount,
		)
		if err != nil {
			log.Println("error: ", err.Error())
			return &models.BooksAudios{}, err
		}
		tem.LikeCount, err = r.BookLikeCount(tem.Id)
		if err != nil {
			return nil, err
		}
		tem.ViewCount, err = r.BookViewCount(tem.Id)
		if err != nil {
			return nil, err
		}
		res.BookAudios = append(res.BookAudios, &tem)
	}
	res.Meta.PageCount = (res.Meta.TotalCount + limit - 1) / limit
	res.Meta.CurrentPage = page
	res.Meta.PerPage = limit

	return &res, nil
}

func (r storagePg) BookGetByAuthorId(id int) ([]*models.BooksForList, error) {
	res := []*models.BooksForList{}
	query := `
SELECT
	b.id, 
	b.title, 
	b.image, 
	b.book_type, 
	a.first_name, 
	a.last_name, 
	b.category_id, 
	b.subcategory_id, 
	b.author_id
FROM 
	books b JOIN author a ON a.id=b.author_id 
WHERE 
	a.id=$1 
`
	rows, err := r.db.Query(query, id)
	if err != nil {
		log.Println("error: ", err.Error())
		return []*models.BooksForList{}, err
	}
	for rows.Next() {
		tem := models.BooksForList{}
		err = rows.Scan(
			&tem.Id, &tem.Title, &tem.Image, &tem.BookType,
			&tem.AuthorFirstName, &tem.AuthorLastName,
			&tem.CategoryId, &tem.SubCategoryId, &tem.AuthorId)
		if err != nil {
			log.Println("error: ", err.Error())
			return []*models.BooksForList{}, err
		}
		tem.LikeCount, err = r.BookLikeCount(tem.Id)
		if err != nil {
			return nil, err
		}
		tem.ViewCount, err = r.BookViewCount(tem.Id)
		if err != nil {
			return nil, err
		}
		res = append(res, &tem)
	}
	return res, nil
}

func (r storagePg) BookGetFilter(req *models.BookFilterReq) (*models.Books, error) {
	res := models.Books{}
	offset := (req.Page - 1) * req.Limit
	query := `
SELECT
	b.id, 
	b.title, 
	b.image, 
	b.book_type, 
	a.first_name, 
	a.last_name, 
	b.category_id, 
	b.subcategory_id, 
	b.author_id,
	COUNT(*) OVER () AS total_count
FROM 
	books b JOIN author a ON a.id=b.author_id 
WHERE 1=1
`

	queryParams := []interface{}{}
	// Counter for placeholders
	placeholderCount := 1

	if req.CategoryId != 0 {
		query += " AND b.category_id = $" + strconv.Itoa(placeholderCount)
		queryParams = append(queryParams, req.CategoryId)
		placeholderCount++
	}
	if req.SubCategoryId != 0 {
		query += " AND b.subcategory_id = $" + strconv.Itoa(placeholderCount)
		queryParams = append(queryParams, req.SubCategoryId)
		placeholderCount++
	}
	if req.AuthorId != 0 {
		query += " AND b.author_id = $" + strconv.Itoa(placeholderCount)
		queryParams = append(queryParams, req.AuthorId)
		placeholderCount++
	}

	if req.Search != "" {
		searchParam := "%" + req.Search + "%"
		query += " AND (b.title ILIKE $" + strconv.Itoa(placeholderCount) +
			" OR b.description ILIKE $" + strconv.Itoa(placeholderCount) +
			" OR a.first_name ILIKE $" + strconv.Itoa(placeholderCount) +
			" OR a.last_name ILIKE $" + strconv.Itoa(placeholderCount) + ")"
		queryParams = append(queryParams, searchParam)
		placeholderCount++
	}
	query += " ORDER BY b.created_at DESC"

	query += " LIMIT $" + strconv.Itoa(placeholderCount) + " OFFSET $" + strconv.Itoa(placeholderCount+1)

	queryParams = append(queryParams, req.Limit, offset)

	rows, err := r.db.Query(query, queryParams...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		tem := models.BooksForList{}
		err = rows.Scan(
			&tem.Id,
			&tem.Title,
			&tem.Image,
			&tem.BookType,
			&tem.AuthorFirstName,
			&tem.AuthorLastName,
			&tem.CategoryId,
			&tem.SubCategoryId,
			&tem.AuthorId, &res.Meta.TotalCount)
		if err != nil {
			log.Println("error: ", err.Error())
			return &models.Books{}, err
		}
		tem.LikeCount, err = r.BookLikeCount(tem.Id)
		if err != nil {
			return nil, err
		}
		tem.ViewCount, err = r.BookViewCount(tem.Id)
		if err != nil {
			return nil, err
		}
		res.BookList = append(res.BookList, &tem)

	}

	res.Meta.PageCount = (res.Meta.TotalCount + req.Limit - 1) / req.Limit
	res.Meta.CurrentPage = req.Page
	res.Meta.PerPage = req.Limit

	return &res, nil
}
