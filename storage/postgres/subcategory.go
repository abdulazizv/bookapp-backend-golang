package postgres

import (
	"log"

	"gitlab.com/bookapp/api/models"
)

func (r storagePg) SubCategoryCreate(req *models.SubCategoryReq) (*models.SubCategoryRes, error) {
	res := models.SubCategoryRes{}
	query := `
INSERT INTO 
	sub_category
	(
		subcategory_name, 
		category_id
	) 
VALUES
	($1,$2) 
RETURNING 
	id, 
	category_id, 
	subcategory_name, 
	created_at, 
	updated_at`

	err := r.db.QueryRow(query, req.SubCategoryName, req.CategoryId).
		Scan(&res.Id, &res.CategoryId, &res.SubCategoryName, &CreatedAt, &UpdatedAt)

	if err != nil {
		log.Fatalln("Error create sub_category", err.Error())
		return &models.SubCategoryRes{}, err
	}
	res.CreatedAt = CreatedAt.Format(Layout)
	res.UpdatedAt = UpdatedAt.Format(Layout)
	return &res, nil
}

func (r storagePg) SubCategoryGet(id int, limit, page int) (*models.SubCategoryRes, error) {
	res := models.SubCategoryRes{}

	query := `
SELECT 
	id, 
	subcategory_name, 
	category_id, 
	created_at, 
	updated_at 
FROM 
	sub_category WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&res.Id, &res.SubCategoryName, &res.CategoryId, &CreatedAt, &UpdatedAt,
	)

	if err != nil {
		log.Println("Error get subcategory:", err.Error())
		return &models.SubCategoryRes{}, err
	}
	res.CreatedAt = CreatedAt.Format(Layout)
	res.UpdatedAt = UpdatedAt.Format(Layout)

	res.Books, err = r.BookGetSubCaID(id, limit, page)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r storagePg) SubCategoryUpdate(req *models.SubCategoryUpdate) (*models.SubCategoryRes, error) {
	res := models.SubCategoryRes{}
	query := `
UPDATE 
	sub_category 
SET 
	subcategory_name=$1, 
	updated_at=now() 
WHERE 
	id=$2 
RETURNING 
	id, 
	subcategory_name, 
	category_id, 
	created_at, 
	updated_at`
	err := r.db.QueryRow(query, req.SubCategoryName, req.Id).Scan(&res.Id, &res.SubCategoryName, &res.CategoryId, &CreatedAt, &UpdatedAt)
	if err != nil {
		log.Println("error update subcategory: ", err.Error())
		return &models.SubCategoryRes{}, err
	}
	res.CreatedAt = CreatedAt.Format(Layout)
	res.UpdatedAt = UpdatedAt.Format(Layout)

	return &res, nil
}

func (r storagePg) SubCategoryDelete(id int) error {
	
	_, err := r.db.Exec(`DELETE FROM views 
	WHERE book_id IN (SELECT id FROM books WHERE subcategory_id = $1);`, id)
	if err != nil {
		log.Println("error delete subcategory: ", err.Error())
		return err
	}

	_, err = r.db.Exec(`DELETE FROM book_likes 
	WHERE book_id IN (SELECT id FROM books WHERE subcategory_id = $1);`, id)
	if err != nil {
		log.Println("error delete subcategory: ", err.Error())
		return err
	}

	_, err = r.db.Exec(`DELETE FROM comments 
	WHERE book_id IN (SELECT id FROM books WHERE subcategory_id = $1);`, id)
	if err != nil {
		log.Println("error delete subcategory: ", err.Error())
		return err
	}

	_, err = r.db.Exec(`DELETE FROM books WHERE subcategory_id = $1;`, id)
	if err != nil {
		log.Println("error delete subcategory: ", err.Error())
		return err
	}
	_, err = r.db.Exec(`DELETE FROM sub_category WHERE id=$1`, id)
	if err != nil {
		log.Println("error delete subcategory: ", err.Error())
		return err
	}
	return nil
}

func (r storagePg) SubCategoryGetCategoryID(id int) ([]*models.SubCategoryRes, error) {
	res := []*models.SubCategoryRes{}
	query := `
SELECT 
	id, 
	category_id, 
	subcategory_name, 
	created_at, 
	updated_at 
FROM 
	sub_category 
WHERE 
	category_id=$1`
	rows, err := r.db.Query(query, id)
	if err != nil {
		log.Println("Error get subcategory by categoryId: ", err.Error())
		return []*models.SubCategoryRes{}, err
	}
	for rows.Next() {
		temp := models.SubCategoryRes{}
		err = rows.Scan(&temp.Id, &temp.CategoryId, &temp.SubCategoryName, &CreatedAt, &UpdatedAt)
		if err != nil {
			log.Println("error: ", err.Error())
			return []*models.SubCategoryRes{}, err
		}
		temp.CreatedAt = CreatedAt.Format(Layout)
		temp.UpdatedAt = UpdatedAt.Format(Layout)
		res = append(res, &temp)
	}
	return res, nil
}
