package repo

import "gitlab.com/bookapp/api/models"

type (
	BookappService interface {
		// category
		CategoryCreate(c *models.CategoryReq) (*models.CategoryResp, error)
		CategoryGet(id int) (*models.CategoryResp, error)
		CategoryUpdate(c *models.CategoryUpdateReq) (*models.CategoryResp, error)
		CategoryDelete(id int) error
		CategoryList() ([]*models.CategoryResp, error)
		CategoryGetId(id int, limit, page int) (*models.CategoryResponse, error)

		// subcategory
		SubCategoryCreate(c *models.SubCategoryReq) (*models.SubCategoryRes, error)
		SubCategoryGet(id int, limit, page int) (*models.SubCategoryRes, error)
		SubCategoryUpdate(*models.SubCategoryUpdate) (*models.SubCategoryRes, error)
		SubCategoryDelete(id int) error
		SubCategoryGetCategoryID(id int) ([]*models.SubCategoryRes, error)

		// Book
		BookCreate(b *models.BookReq) (*models.BookRes, error)
		BookGet(id int, useriD int) (*models.BookResponse, error)
		BookUpdate(res *models.BookUpdate) (*models.BookRes, error)
		BookGetSearch(key string) (*models.Books, error)
		BookDelete(id int) error
		BookGetList(limit int, page int, key string) (*models.Books, error)
		BookGetSubCaID(id int, limit, page int) (*models.Books, error)
		BookGetCatId(id int, limit, page int) (*models.Books, error)
		BookGetCategoryId(id int) ([]*models.BooksForList, error)
		BookCreateLike(l *models.LikeReq) error
		BookDeleteLike(l *models.LikeReq) error
		BookGetLiked(id int) ([]*models.BooksForList, error)
		BookGetTops(limit, page int) (*models.Books, error)
		BookGetMoreRead(limit, page int) (*models.Books, error)
		BookGetFilter(req *models.BookFilterReq) (*models.Books, error)
		BookGetAudios(limit, page int, search string) (*models.BooksAudios, error)

		//author
		AuthorCreate(a *models.AuthorReq) (*models.AuthorRes, error)
		AuthorGet(id int) (*models.AuthorRes, error)
		AuthorUpdate(a *models.AuthorUpdateReq) (*models.AuthorRes, error)
		AuthorDelete(id int) error
		AuthorGetList(limit int, page int, search string) (*models.Authors, error)

		// comment
		CommentCreate(c *models.CommentReq) (*models.CommentRes, error)
		CommentGet(id int) (*models.CommentRes, error)
		CommentUpdate(c *models.CommentUpdate) (*models.CommentRes, error)
		CommentDelete(id int) error
		CommentGetBookID(id int) ([]*models.CommentResponse, error)

		// user
		UserCreate(u *models.UserReq) (*models.UserRes, error)
		UserGet(id int) (*models.UserRes, error)
		UserUpdate(res *models.UserUpdateReq) (*models.UserRes, error)
		UserDelete(id int) error
		CheckField(req *models.CheckFieldReq) (*models.CheckFieldRes, error)
		UserGetLogin(login string) (*models.UserLoginRes, error)
		AdminGetList() ([]*models.UserResForComment, error)

		ViewCreate(req *models.View) error
		CheckFieldView(req *models.CheckFieldViewReq) (*models.CheckFieldRes, error)

		GetStatistic() (*models.StatisticCount, error)
		GetCategoryBookCount() ([]*models.CategoryBookCount, error)
		GetWeekAddedBook() ([]*models.AddedBooks, error)
	}
)

type InMemoryStorageI interface {
	Set(key, value string) error
	SetWithTTL(key, value string, seconds int) error
	Get(key string) (interface{}, error)
	Exists(key string) (interface{}, error)
	Del(key string) (interface{}, error)
	Keys(key string) (interface{}, error)
}
