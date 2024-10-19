package postgres

import "gitlab.com/bookapp/api/models"

func (r storagePg) GetStatistic() (*models.StatisticCount, error) {
	res := models.StatisticCount{}
	err := r.db.QueryRow(`SELECT COUNT(*) FROM author`).Scan(&res.AuthorCount)
	if err != nil {
		return nil, err
	}
	res.Author = "author"

	err = r.db.QueryRow(`SELECT COUNT(*) FROM books`).Scan(&res.BookCount)
	if err != nil {
		return nil, err
	}
	res.Book = "book"

	err = r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&res.UserCount)
	if err != nil {
		return nil, err
	}
	res.User = "user"

	err = r.db.QueryRow(`SELECT COUNT(*) FROM books WHERE book_type='top'`).Scan(&res.TopBookCount)
	if err != nil {
		return nil, err
	}

	res.TopBook = "top_book"
	return &res, nil
}

func (r storagePg) GetCategoryBookCount() ([]*models.CategoryBookCount, error) {
	res := []*models.CategoryBookCount{}
	query := `SELECT c.category_name, COUNT(b.category_id) FROM category c JOIN books b ON b.category_id=c.id GROUP BY c.category_name`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		temp := models.CategoryBookCount{}
		err = rows.Scan(&temp.CategoryName, &temp.BookCount)
		if err != nil {
			return nil, err
		}
		res = append(res, &temp)
	}
	return res, err
}

func (r storagePg) GetWeekAddedBook() ([]*models.AddedBooks, error) {
	res := []*models.AddedBooks{}
	query := `
SELECT
    DATE_TRUNC('week', created_at + INTERVAL '1 day' - EXTRACT(ISODOW FROM created_at) * INTERVAL '1 day') AS week_start,
    COUNT(*) AS book_count
FROM 
	books
WHERE 
	created_at >= DATE '2023-01-01' AND created_at < DATE '2024-01-01'
GROUP BY 1
ORDER BY 1;`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		temp := models.AddedBooks{}
		err = rows.Scan(&CreatedAt, &temp.BookCount)
		if err != nil {
			return nil, err
		}
		temp.WeeklyDate = CreatedAt.Format(Layout)
		res = append(res, &temp)
	}
	return res, err
}
