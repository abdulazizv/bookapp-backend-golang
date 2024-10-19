package models

type UserReq struct {
	RoleId       int    `json:"role_id" example:"3"`
	FullName     string `json:"full_name"`
	AvatarUrl    string `json:"avatar_url" example:"https://shorturl.at/akpDK"`
	Login        string `json:"login"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
}

type UserRes struct {
	Id        int             `json:"id"`
	FullName  string          `json:"full_name"`
	AvatarUrl string          `json:"avatar_url"`
	Login     string          `json:"login"`
	Books     []*BooksForList `json:"like_books"`
}

type UserLoginRes struct {
	Id          int    `json:"id"`
	FullName    string `json:"full_name"`
	AvatarUrl   string `json:"avatar_url"`
	Login       string `json:"login"`
	RoleId      int    `json:"role_id"`
	Password    string `json:"password"`
	AccessToken string `json:"access_token"`
}

type UserResForComment struct {
	Id        int    `json:"id"`
	FullName  string `json:"full_name"`
	AvatarUrl string `json:"avatar_url"`
}

type UserUpdateReq struct {
	Id        int    `json:"id"`
	FullName  string `json:"full_name"`
	AvatarUrl string `json:"avatar_url"`
	Login     string `json:"login"`
	Password  string `json:"password"`
}

type CheckFieldReq struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type CheckFieldRes struct {
	Exists bool `json:"exists"`
}

type LoginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CheckFieldViewReq struct {
	Field  string `json:"field"`
	Value  string `json:"value"`
	BookId int    `json:"book_id"`
}
