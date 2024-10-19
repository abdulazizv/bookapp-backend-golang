package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/bookapp/api/handler/v1/http"
	"gitlab.com/bookapp/api/models"
)

// @Router		/category [POST]
// @Security 	BearerAuth
// @Summary		create category
// @Tags        Category
// @Description	Here category will be created.
// @Accept 		json
// @Produce		json
// @Param       body    body models.CategoryReq true "Create"
// @Success		201 	{object}  models.CategoryResp
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) CreateCategory(c *gin.Context) {
	body := models.CategoryReq{}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	res, err := h.Storage.BookApp().CategoryCreate(&body)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.Created, res)
}

// Get Category ById
// @ID getcategory
// @Router		/category/{id} [GET]
// @Summary		get category
// @Tags        Category
// @Description	Here get by category.
// @Accept 		json
// @Produce		json
// @Param 		id path int true "id"
// @Success 	200 {object} http.Response{data=models.CategoryResp} "Category data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) GetByIdCategory(c *gin.Context) {
	id := c.Param("id")
	cid, err := strconv.Atoi(id)

	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	res, err := h.Storage.BookApp().CategoryGet(cid)

	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, res)
}

// Update  		Category
// @Security 	BearerAuth
// @ID 			updatecategory
// @Router 		/category [PUT]
// @Summary 	Update Category
// @Description Update Category
// @Tags 		Category
// @Accept 		json
// @Produce 	json
// @Param 		body 	body models.CategoryUpdateReq true "UpdateCategory"
// @Success 	200 {object} http.Response{data=models.CategoryResp} "Category data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) UpdateCategory(c *gin.Context) {

	var body models.CategoryUpdateReq

	err := c.ShouldBindJSON(&body)
	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	res, err := h.Storage.BookApp().CategoryUpdate(&body)

	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, res)
}

// Delete Category
// @Security 	BearerAuth
// @ID 			deletecategory
// @Router 		/category/{id} [DELETE]
// @Summary 	Delete Category
// @Description Delete Category
// @Tags		Category
// @Accept 		json
// @Produce 	json
// @Param 		id 	path int true "id"
// @Success 	200 {object} http.Response{data=models.Success} "Category data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	cid, err := strconv.Atoi(id)

	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	err = h.Storage.BookApp().CategoryDelete(cid)

	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, "Successfully deleted")
}

// Get Categories
// @ID getcategories
// @Router		/category/list [GET]
// @Summary		get category
// @Tags        Category
// @Description	Here category will be get.
// @Accept 		json
// @Produce		json
// @Success 	200 {object} http.Response{data=[]models.CategoryResp} "Category data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) GetListCategory(c *gin.Context) {
	res, err := h.Storage.BookApp().CategoryList()
	
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, res)
}

// Get Category ById
// @ID 			getcategoryid
// @Router		/category/books [GET]
// @Summary		get category
// @Tags        Category
// @Description	Get all books by categoryId.
// @Accept 		json
// @Produce		json
// @Param 		id query int true "id"
// @Param 		limit query int false "LIMIT"
// @Param 		page query int false "PAGE_NUMBER"
// @Success 	200 {object} http.Response{data=models.CategoryResponse} "Category data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) GetCategoryId(c *gin.Context) {
	id := c.Query("id")
	var limit, page int
	limitStr, pageStr := c.Query("limit"), c.Query("page")
	cid, err := strconv.Atoi(id)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			h.handleResponse(c, http.BadRequest, err.Error())
			return
		}
	} else {
		limit = 10
	}

	if pageStr != "" {
		page, err = strconv.Atoi(limitStr)
		if err != nil {
			h.handleResponse(c, http.BadRequest, err.Error())
			return
		}
	} else {
		page = 1
	}
	res, err := h.Storage.BookApp().CategoryGetId(cid, limit, page)

	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, res)
}
