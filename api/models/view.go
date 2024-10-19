package models

type View struct {
	BookId    int    `json:"book_id"`
	UserAgent string `json:"user_agent"`
	Count     int    `json:"count"`
}
