package models

type CommentReq struct {
	UserId int    `json:"user_id"`
	BookId int    `json:"book_id"`
	Text   string `json:"text"`
}

type CommentRes struct {
	Id        int    `json:"id"`
	UserId    int    `json:"user_id"`
	BookId    int    `json:"book_id"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CommentUpdate struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
}
type CommentResponse struct {
	Id        int               `json:"id"`
	UserId    int               `json:"user_id"`
	BookId    int               `json:"book_id"`
	Text      string            `json:"text"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
	User      UserResForComment `json:"user"`
}
