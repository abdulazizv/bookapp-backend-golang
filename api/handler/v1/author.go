package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/bookapp/api/handler/v1/http"
	"gitlab.com/bookapp/api/models"
)

// Create Author
// @ID createauthor
// @Router		/author [POST]
// @Security 	BearerAuth
// @Summary		create author
// @Tags        Author
// @Description	Here author can update info
// @Accept 		json
// @Produce		json
// @Param 		body body 	 models.AuthorReq  true "AuthorCreate"
// @Success 	201 {object} http.Response{data=models.AuthorRes} "Author data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) CreateAuthor(c *gin.Context) {
	var body models.AuthorReq

	err := c.ShouldBindJSON(&body)
	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	res, err := h.Storage.BookApp().AuthorCreate(&body)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.Created, res)
}

// GET Author
// @ID getauthor
// @Router		/author/{id} [GET]
// @Summary		GET author
// @Tags        Author
// @Description	Here author can get info
// @Accept 		json
// @Produce		json
// @Param 		id	path  int true "GetAuthor"
// @Success 	200 {object} http.Response{data=models.AuthorRes} "Author data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) GetAuthor(c *gin.Context) {
	id := c.Param("id")
	aId, err := strconv.Atoi(id)

	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	res, err := h.Storage.BookApp().AuthorGet(aId)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, res)
}

// GET Author
// @ID getauthors
// @Router		/author/list [GET]
// @Summary		GET author
// @Tags        Author
// @Description	Here author can get info
// @Accept 		json
// @Produce		json
// @Param 		limit query int false "LIMIT"
// @Param 		page query int false "PAGE_NUMBER"
// @Param 		search query string false "SEARCH"
// @Success 	200 {object} http.Response{data=models.Authors} "Author data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) GetAuthorList(c *gin.Context) {
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
	search := c.Query("search")
	res, err := h.Storage.BookApp().AuthorGetList(limitI, pageI, search)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, res)
}

// GET Author
// @ID updateauthor
// @Security 	BearerAuth
// @Router		/author [PUT]
// @Summary		put author
// @Tags        Author
// @Description	Here author can update info
// @Accept 		json
// @Produce		json
// @Param 		body  body models.AuthorUpdateReq  true "UpdateAuthor"
// @Success 	200 {object} http.Response{data=models.AuthorRes} "Author data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) UpdateAuthor(c *gin.Context) {
	var body models.AuthorUpdateReq
	err := c.ShouldBindJSON(&body)
	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	res, err := h.Storage.BookApp().AuthorUpdate(&body)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, res)
}

// delete Author
// @ID deleteauthor
// @Security 	BearerAuth
// @Router		/author/{id} [DELETE]
// @Summary		DELETE author
// @Tags        Author
// @Description	Here author can DELETE info
// @Accept 		json
// @Produce		json
// @Param 		id	path 	 int true "DELETEAuthor"
// @Success 	200 {object} http.Response{data=string} "Author delete"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) DeleteAuthor(c *gin.Context) {
	id := c.Param("id")
	aId, err := strconv.Atoi(id)

	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	err = h.Storage.BookApp().AuthorDelete(aId)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, "Author deleted")
}
