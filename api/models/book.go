package models

//book

type BookReq struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	Image         string `json:"image"`
	DownloadUrl   string `json:"download_url"`
	AudioUrl      string `json:"audio_url"`
	BookType      string `json:"book_type" example:"standard"`
	CategoryID    int    `json:"category_id"`
	SubCategoryID int    `json:"sub_category_id"`
	AuthorID      int    `json:"author_id"`
}

type BookRes struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Image         string `json:"image"`
	DownloadUrl   string `json:"download_url"`
	AudioUrl      string `json:"audio_url"`
	BookType      string `json:"book_type"`
	CategoryID    int    `json:"category_id"`
	SubCategoryID int    `json:"sub_category_id"`
	AuthorID      int    `json:"author_id"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type BookUpdate struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	DownloadUrl string `json:"download_url"`
	AudioUrl    string `json:"audio_url"`
	BookType    string `json:"book_type"`
}

type BookResponse struct {
	ID              int                `json:"id"`
	IsLike          bool               `json:"is_like"`
	Title           string             `json:"title"`
	Description     string             `json:"description"`
	Image           string             `json:"image"`
	DownloadUrl     string             `json:"download_url"`
	AudioUrl        string             `json:"audio_url"`
	BookType        string             `json:"book_type"`
	CategoryID      int                `json:"category_id"`
	SubCategoryID   int                `json:"sub_category_id"`
	AuthorID        int                `json:"author_id"`
	CreatedAt       string             `json:"created_at"`
	UpdatedAt       string             `json:"updated_at"`
	AuthorFirstName string             `json:"author_first_name"`
	AuthorLastName  string             `json:"author_last_name"`
	ViewCount       int                `json:"view_count"`
	LikeCount       int                `json:"like_count"`
	CommentCount    int                `json:"comment_count"`
	Comments        []*CommentResponse `json:"comments"`
	SimilarBooks    []*BooksForList    `json:"similar_books"`
}

type BookListRes struct {
	Books []BookRes `json:"books"`
}

type BookFilterReq struct {
	CategoryId    int
	SubCategoryId int
	AuthorId      int
	Limit         int
	Page          int
	Search        string
}

type BooksForList struct {
	Id              int    `json:"id"`
	Title           string `json:"title"`
	Image           string `json:"image"`
	BookType        string `json:"book_type"`
	AuthorFirstName string `json:"author_first_name"`
	AuthorLastName  string `json:"author_last_name"`
	LikeCount       int    `json:"like_count"`
	ViewCount       int    `json:"view_count"`
	CategoryId      int    `json:"category_id"`
	SubCategoryId   int    `json:"sub_category_id"`
	AuthorId        int    `json:"author_id"`
}

type BooksAudioList struct {
	Id              int    `json:"id"`
	Title           string `json:"title"`
	Image           string `json:"image"`
	AudioURL        string `json:"audio_url"`
	BookType        string `json:"book_type"`
	AuthorFirstName string `json:"author_first_name"`
	AuthorLastName  string `json:"author_last_name"`
	LikeCount       int    `json:"like_count"`
	ViewCount       int    `json:"view_count"`
	CategoryId      int    `json:"category_id"`
	SubCategoryId   int    `json:"sub_category_id"`
	AuthorId        int    `json:"author_id"`
}

type Meta struct {
	TotalCount  int `json:"total_count"`
	PageCount   int `json:"page_count"`
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
}

type Books struct {
	BookList []*BooksForList `json:"books"`
	Meta     Meta            `json:"_meta"`
}

type BooksAudios struct {
	BookAudios []*BooksAudioList `json:"books"`
	Meta       Meta              `json:"_meta"`
}
type BookPaginationReq struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

type LikeReq struct {
	UserId int `json:"user_id"`
	BookId int `json:"book_id"`
}

// author
type AuthorReq struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	BirthDay   string `json:"birth_day" example:"1960-02-15"`
	DiedYear   string `json:"died_year" example:"2000-03-24"`
	Country    string `json:"country"`
	AvatarUrl  string `json:"avatar_url"`
	AboutText  string `json:"about_text"`
	Creativity string `json:"creativity"`
}

type AuthorUpdateReq struct {
	ID         int    `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	BirthDay   string `json:"birth_day" example:"1960-02-15"`
	DiedYear   string `json:"died_year" example:"2000-03-24"`
	Country    string `json:"country"`
	AvatarUrl  string `json:"avatar_url"`
	AboutText  string `json:"about_text"`
	Creativity string `json:"creativity"`
}

type AuthorForList struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDay  string `json:"birth_day"`
	DiedYear  string `json:"died_year"`
	AvatarUrl string `json:"avatar_url"`
	BookCount int    `json:"book_count"`
}

type Authors struct {
	Authors []*AuthorForList `json:"authors"`
	Meta    Meta             `json:"_meta"`
}

type AuthorRes struct {
	ID         int             `json:"id"`
	FirstName  string          `json:"first_name"`
	LastName   string          `json:"last_name"`
	MiddleName string          `json:"middle_name"`
	BirthDay   string          `json:"birth_day" example:"1960-02-15"`
	DiedYear   string          `json:"died_year" example:"2000-03-24"`
	Country    string          `json:"country"`
	AvatarUrl  string          `json:"avatar_url"`
	AboutText  string          `json:"about_text"`
	Creativity string          `json:"creativity"`
	CreatedAt  string          `json:"created_at"`
	UpdatedAt  string          `json:"updated_at"`
	Books      []*BooksForList `json:"books"`
}

type BookNotify struct {
	Title    string `json:"title"`
	ImageUrl string `json:"image_url"`
	Body     string `json:"body"`
}

type StatisticCount struct {
	Book         string `json:"book"`
	BookCount    int    `json:"book_count"`
	Author       string `json:"author"`
	AuthorCount  int    `json:"author_count"`
	User         string `json:"user"`
	UserCount    int    `json:"user_count"`
	TopBook      string `json:"top_book"`
	TopBookCount int    `json:"top_book_count"`
}

type CategoryBookCount struct {
	CategoryName string `json:"category_name"`
	BookCount    int    `json:"book_count"`
}

type AddedBooks struct {
	WeeklyDate string `json:"week_date"`
	BookCount  int    `json:"book_count"`
}
