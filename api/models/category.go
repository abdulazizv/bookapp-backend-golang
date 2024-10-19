package models

type CategoryReq struct {
	Name   string `json:"category_name"`
	Status bool   `json:"status"`
}

type CategoryResp struct {
	Id            int               `json:"id"`
	Name          string            `json:"category_name"`
	Status        bool              `json:"status"`
	CreatedAt     string            `json:"created_at"`
	UpdatedAt     string            `json:"updated_at"`
	SubCategories []*SubCategoryRes `json:"subcategories"`
}

type CategoryResponse struct {
	Id        int    `json:"id"`
	Name      string `json:"category_name"`
	Status    bool   `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	*Books
}

type Success struct {
	Message string `json:"message"`
}

type ListCategory struct {
	Categories []CategoryResp `json:"categories"`
}

type CategoryUpdateReq struct {
	Id           int    `json:"id"`
	CategoryName string `json:"category_name"`
}
