package v1

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/bookapp/api/handler/v1/http"
	"gitlab.com/bookapp/api/models"
	"gitlab.com/bookapp/pkg/util"
)

// @Router		/book [POST]
// @Summary		create books
// @Security BearerAuth
// @Tags        Book
// @Description	Here books will be created.
// @Accept 		json
// @Produce		json
// @Param       body    body models.BookReq true "Create"
// @Success		201 	{object}  models.BookRes
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) CreateBook(c *gin.Context) {
	body := models.BookReq{}

	err := c.ShouldBindJSON(&body)

	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	res, err := h.Storage.BookApp().BookCreate(&body)

	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.Created, res)

	notify := models.BookNotify{
		Title:    "🎉📚 Hey, Bizda yangilik!!!",
		ImageUrl: body.Image,
		Body:     fmt.Sprintf("Bookapp dasturiga %s nomli yangi kitob qo'shildi", body.Title),
	}
	err = util.SendNotify(&notify)
	if err != nil {
		log.Println("Error while sending notify: ", err.Error())
	}
}

// Get Book ById
// @ID 			getbook
// @Router		/book/{id} [GET]
// @Summary		get books
// @Tags        Book
// @Description	Here get by books.
// @Accept 		json
// @Produce		json
// @Param 		id path int true "id"
// @Param 		user_id query int false "user_id"
// @Success 	200 {object} http.Response{data=models.BookResponse} "Book data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) GetByIdBook(c *gin.Context) {
	var userId int
	id := c.Param("id")
	userid := c.Query("user_id")
	bid, err := strconv.Atoi(id)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	if userid != "" {
		userId, err = strconv.Atoi(userid)
		if err != nil {
			h.handleResponse(c, http.BadRequest, err.Error())
			return
		}
	}
	UserAgent := c.Request.UserAgent()

	res, err := h.Storage.BookApp().BookGet(bid, userId)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	newView := &models.View{
		BookId:    bid,
		UserAgent: UserAgent,
		Count:     1,
	}

	err = h.CreateViews(newView)
	if err != nil {
		log.Println("Error: h.CreateViews(&newView): ", err.Error())
	}

	h.handleResponse(c, http.OK, res)
}

func (h *handlerV1) CreateViews(req *models.View) error {
	checkAgent, err := h.Storage.BookApp().CheckFieldView(&models.CheckFieldViewReq{
		Field:  "user_agent",
		Value:  req.UserAgent,
		BookId: req.BookId,
	})
	if err != nil {
		return err
	}
	if !(checkAgent.Exists) {
		err = h.Storage.BookApp().ViewCreate(req)
		if err != nil {
			log.Println("Error create view: ", err.Error())
		}
	}

	return nil
}

// Update 		Book
// @ID 			updatebook
// @Security 	BearerAuth
// @Router 		/book [PUT]
// @Summary 	Update Book
// @Description Update Book
// @Tags 		Book
// @Accept 		json
// @Produce 	json
// @Param 		body body models.BookUpdate true "UpdateBook"
// @Success 	200 {object} http.Response{data=models.BookRes} "Books data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) UpdateBook(c *gin.Context) {
	var body models.BookUpdate

	err := c.ShouldBindJSON(&body)

	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	res, err := h.Storage.BookApp().BookUpdate(&body)

	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, res)
}

// Delete Books
// @ID deletebook
// @Security BearerAuth
// @Router /book/{id} [DELETE]
// @Summary Delete Category
// @Description Delete Category
// @Tags Book
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} http.Response{data=models.Success} "Category data"
// @Response 400 {object} http.Response{data=string} "Bad Request"
// @Failure 500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) DeleteBook(c *gin.Context) {
	id := c.Param("id")
	cid, err := strconv.Atoi(id)

	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	err = h.Storage.BookApp().BookDelete(cid)

	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, "Successfully deleted")
}

// GET Books
// @ID 			getbooklist
// @Router 		/book/search [GET]
// @Summary 	GET Books
// @Description GET Books
// @Tags 		Book
// @Accept 		json
// @Produce 	json
// @Param 		limit  query int false "LIMIT"
// @Param 		page   query int false "PAGE_NUMBER"
// @Param 		search query string false "SEARCH_KEY"
// @Success 	200 {object} http.Response{data=models.Books} "Books data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) GetListBooks(c *gin.Context) {
	limit, page := c.Query("limit"), c.Query("page")
	search := c.Query("search")
	limitI, pageI := h.CheckLimitOffset(limit, page)

	res, err := h.Storage.BookApp().BookGetList(limitI, pageI, search)

	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, res)
}

// GET Books
// @ID 			getbooks
// @Router 		/book/filter [GET]
// @Summary 	GET Books
// @Description GET Books
// @Tags 		Book
// @Accept 		json
// @Produce 	json
// @Param 		limit  query int false "LIMIT"
// @Param 		page   query int false "PAGE_NUMBER"
// @Param 		search query string false "SEARCH_KEY"
// @Param 		category_id query int false "CATEGORY_ID"
// @Param 		subcategory_id query int false "SUBCATEGORY_ID"
// @Param       author_id  query int false "AUTHOR_ID"
// @Success 	200 {object} http.Response{data=models.Books} "Books data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) GetBooksFilter(c *gin.Context) {
	limitStr, pageStr := c.Query("limit"), c.Query("page")
	search := c.Query("search")
	categoryIDStr := c.Query("category_id")
	subcategoryIDStr := c.Query("subcategory_id")
	authorIDStr := c.Query("author_id")

	var limit, page, categoryID, subcategoryID, authorID int
	var err error

	limit, page = h.CheckLimitOffset(limitStr, pageStr)

	if categoryIDStr != "" {
		categoryID, err = strconv.Atoi(categoryIDStr)
		if err != nil {
			h.handleResponse(c, http.InternalServerError, err.Error())
			return
		}
	}

	if subcategoryIDStr != "" {
		subcategoryID, err = strconv.Atoi(subcategoryIDStr)
		if err != nil {
			h.handleResponse(c, http.InternalServerError, err.Error())
			return
		}
	}

	if authorIDStr != "" {
		authorID, err = strconv.Atoi(authorIDStr)
		if err != nil {
			h.handleResponse(c, http.InternalServerError, err.Error())
			return
		}
	}

	res, err := h.Storage.BookApp().BookGetFilter(&models.BookFilterReq{
		Search:        search,
		CategoryId:    categoryID,
		SubCategoryId: subcategoryID,
		AuthorId:      authorID,
		Limit:         limit,
		Page:          page,
	})

	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	h.handleResponse(c, http.OK, res)

}

// GET Books
// @ID 			getbooksearch
// @Router 		/book/top [GET]
// @Summary 	GET TOP Books
// @Description GET TOP books
// @Tags 		Book
// @Accept 		json
// @Produce 	json
// @Param 		limit query string false "LIMIT"
// @Param 		page  query string false "PAGE_NUMBER"
// @Success 	200 {object} http.Response{data=models.Books} "Books data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) GetBookTop(c *gin.Context) {
	limit, page := c.Query("limit"), c.Query("page")
	if limit == "" {
		limit = "10"
	}
	if page == "" {
		page = "1"
	}
	limitI, err := strconv.Atoi(limit)
	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}
	pageI, err := strconv.Atoi(page)
	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}
	res, err := h.Storage.BookApp().BookGetTops(limitI, pageI)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, res)
}

// GET Books
// @ID 			getbookmostread
// @Router 		/book/mostread [GET]
// @Summary 	GET most read books
// @Description GET list of most read books
// @Tags 		Book
// @Accept 		json
// @Produce 	json
// @Param 		limit query string false "LIMIT"
// @Param 		page  query string false "PAGE_NUMBER"
// @Success 	200 {object} http.Response{data=models.Books} "Books data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) GetBookReadALot(c *gin.Context) {
	limit, page := c.Query("limit"), c.Query("page")
	limitI, pageI := h.CheckLimitOffset(limit, page)
	res, err := h.Storage.BookApp().BookGetMoreRead(limitI, pageI)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, res)
}

// GET Books audios
// @ID 			getbookaudios
// @Router 		/book/audios [GET]
// @Summary 	GET audios
// @Description GET audios of books
// @Tags 		Book
// @Accept 		json
// @Produce 	json
// @Param 		limit query string false "LIMIT"
// @Param 		page  query string false "PAGE_NUMBER"
// @Param 		search query string false "SEARCH"
// @Success 	200 {object} http.Response{data=models.BooksAudios} "Books data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) GetBookAudios(c *gin.Context) {
	limit, page := c.Query("limit"), c.Query("page")
	search := c.Query("search")
	limitI, pageI := h.CheckLimitOffset(limit, page)
	res, err := h.Storage.BookApp().BookGetAudios(limitI, pageI, search)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, res)
}

// @Router		/book/like [POST]
// @Summary		create book like
// @Security 	BearerAuth
// @Tags        Book
// @Description	Here books will be created like.
// @Accept 		json
// @Produce		json
// @Param       body    body models.LikeReq true "CreateLike"
// @Success		201 	{object}  http.Response{data=string} "CreatedLike"
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) CreateBookLike(c *gin.Context) {
	body := models.LikeReq{}

	err := c.ShouldBindJSON(&body)

	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	err = h.Storage.BookApp().BookCreateLike(&body)

	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.Created, models.Message{
		Message: "Liked",
	})
}

// @Router		/book/like [DELETE]
// @Summary		delete book like
// @Security 	BearerAuth
// @Tags        Book
// @Description	Here books will be delete like.
// @Accept 		json
// @Produce		json
// @Param       body    body models.LikeReq true "DeleteLike"
// @Success		200 	{object}  http.Response{data=string} "DeletedLike"
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) DeleteBookLike(c *gin.Context) {
	body := models.LikeReq{}

	err := c.ShouldBindJSON(&body)

	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	err = h.Storage.BookApp().BookDeleteLike(&body)

	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.Created, models.Message{
		Message: "Deleted like",
	})
}

func (h *handlerV1) CheckLimitOffset(limitStr, pageStr string) (int, int) {
	var limit, page int
	var err error
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			log.Println("error parse limit: ", err.Error())
		}
	} else {
		limit = 10
	}
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			log.Println("error parse page: ", err.Error())
		}
	} else {
		page = 1
	}

	return limit, page
}
