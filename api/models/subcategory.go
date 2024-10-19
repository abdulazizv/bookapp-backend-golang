package models

type SubCategoryReq struct {
	SubCategoryName string `json:"subcategory_name"`
	CategoryId      int    `json:"category_id"`
}

type SubCategoryRes struct {
	Id              int    `json:"id"`
	SubCategoryName string `json:"subcategory_name"`
	CategoryId      int    `json:"category_id"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	*Books
}

type SubCategoryUpdate struct {
	Id              int    `json:"id"`
	SubCategoryName string `json:"subcategory_name"`
}
