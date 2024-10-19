package postgres

import (
	"log"

	"gitlab.com/bookapp/api/models"
)

func (r storagePg) CategoryCreate(req *models.CategoryReq) (*models.CategoryResp, error) {
	res := models.CategoryResp{}
	query := `INSERT INTO category(category_name, status) VALUES($1, $2) RETURNING id, category_name, status, created_at, updated_at`
	err := r.db.QueryRow(query, req.Name, req.Status).
		Scan(&res.Id, &res.Name, &res.Status, &CreatedAt, &UpdatedAt)
	if err != nil {
		log.Fatalln("Error create category: ", err.Error())
		return &models.CategoryResp{}, err
	}
	res.CreatedAt = CreatedAt.Format(Layout)
	res.UpdatedAt = UpdatedAt.Format(Layout)

	return &res, nil
}

func (r storagePg) CategoryGet(id int) (*models.CategoryResp, error) {
	res := models.CategoryResp{}

	query := `SELECT id, category_name, status, created_at, updated_at FROM category WHERE id=$1 AND status=true`
	err := r.db.QueryRow(query, id).Scan(
		&res.Id, &res.Name, &res.Status, &CreatedAt, &UpdatedAt,
	)
	if err != nil {
		log.Println("Error get category: ", err.Error())
		return &models.CategoryResp{}, err
	}
	res.CreatedAt = CreatedAt.Format(Layout)
	res.UpdatedAt = UpdatedAt.Format(Layout)

	res.SubCategories, err = r.SubCategoryGetCategoryID(id)
	if err != nil {
		log.Println("error: ", err.Error())
		return &models.CategoryResp{}, err
	}

	return &res, nil
}

func (r storagePg) CategoryUpdate(req *models.CategoryUpdateReq) (*models.CategoryResp, error) {
	res := models.CategoryResp{}
	query := `UPDATE category SET category_name=$1, updated_at=now() WHERE id=$2 
		RETURNING id, category_name, status, created_at, updated_at`
	err := r.db.QueryRow(query, req.CategoryName, req.Id).Scan(
		&res.Id, &res.Name, &res.Status, &CreatedAt, &UpdatedAt,
	)
	if err != nil {
		log.Println("error udpate category: ", err.Error())
		return &models.CategoryResp{}, err
	}
	res.CreatedAt = CreatedAt.Format(Layout)
	res.UpdatedAt = UpdatedAt.Format(Layout)
	return &res, nil
}

func (r storagePg) CategoryDelete(id int) error {

	_, err := r.db.Exec(`DELETE FROM views 
	WHERE book_id IN (SELECT id FROM books WHERE category_id = $1);`, id)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`DELETE FROM book_likes 
	WHERE book_id IN (SELECT id FROM books WHERE category_id = $1);`, id)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`DELETE FROM comments 
	WHERE book_id IN (SELECT id FROM books WHERE category_id = $1);`, id)
	if err != nil {
		log.Println("error delete subcategory: ", err.Error())
		return err
	}
	_, err = r.db.Exec(`DELETE FROM books WHERE category_id=$1`, id)
	if err != nil {
		log.Println("Error delete books: ", err.Error())
		return err
	}

	_, err = r.db.Exec(`DELETE FROM sub_category WHERE category_id=$1`, id)

	if err != nil {
		log.Println("Error delete sub_category: ", err.Error())
		return err
	}

	_, err = r.db.Exec(`DELETE FROM category WHERE id=$1`, id)
	if err != nil {
		log.Println("Error delete category: ", err.Error())
		return err
	}

	return nil
}

func (r storagePg) CategoryList() ([]*models.CategoryResp, error) {
	res := []*models.CategoryResp{}
	query := `SELECT id, category_name, status, created_at, updated_at FROM category WHERE status=true`
	rows, err := r.db.Query(query)
	if err != nil {
		log.Println("Error: ", err.Error())
		return []*models.CategoryResp{}, err
	}
	for rows.Next() {
		temp := models.CategoryResp{}
		err = rows.Scan(&temp.Id, &temp.Name, &temp.Status, &CreatedAt, &UpdatedAt)
		if err != nil {
			log.Println("error: ", err.Error())
			return []*models.CategoryResp{}, err
		}
		temp.SubCategories, err = r.SubCategoryGetCategoryID(temp.Id)
		if err != nil {
			return nil, err
		}
		temp.CreatedAt = CreatedAt.Format(Layout)
		temp.UpdatedAt = UpdatedAt.Format(Layout)
		res = append(res, &temp)
	}
	return res, nil
}

func (r storagePg) CategoryGetId(id int, limit, page int) (*models.CategoryResponse, error) {
	res := models.CategoryResponse{}
	query := `
SELECT 
	id, 
	category_name, 
	status, 
	created_at, 
	updated_at 
FROM 
	category 
WHERE 
	id=$1 AND status=true`
	err := r.db.QueryRow(query, id).Scan(&res.Id, &res.Name, &res.Status, &CreatedAt, &UpdatedAt)
	if err != nil {
		return nil, err
	}
	res.CreatedAt = CreatedAt.Format(Layout)
	res.UpdatedAt = UpdatedAt.Format(Layout)
	res.Books, err = r.BookGetCatId(id, limit, page)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
