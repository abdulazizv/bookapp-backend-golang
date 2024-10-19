package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/bookapp/api/handler/v1/http"
	"gitlab.com/bookapp/api/models"
)

// @Router		/subcategory [POST]
// @Security 	BearerAuth
// @Summary		create subcategory
// @Tags        Subcategory
// @Description	Here category will be created.
// @Accept 		json
// @Produce		json
// @Param       body    body models.SubCategoryReq true "Create"
// @Success		200 	{object}  models.SubCategoryRes
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) CreateSubCategory(c *gin.Context) {
	body := models.SubCategoryReq{}

	err := c.ShouldBindJSON(&body)

	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	res, err := h.Storage.BookApp().SubCategoryCreate(&body)

	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	h.handleResponse(c, http.Created, res)
}

// Get SubCategory ById
// @ID 			getsubcategory
// @Router		/subcategory/{id} [GET]
// @Summary		get subcategory
// @Tags        Subcategory
// @Description	Here get by subcategory.
// @Accept 		json
// @Produce		json
// @Param 		id path int true "id"
// @Param 		limit query int false "LIMIT"
// @Param 		page  query int false "PAGE_NUMBER"
// @Success 	200 {object} http.Response{data=models.SubCategoryRes} "subcategory data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) GetSubCategoryById(c *gin.Context) {
	id := c.Param("id")
	limitStr, pageStr := c.Query("limit"), c.Query("page")
	limit, page := h.CheckLimitOffset(limitStr, pageStr)
	sid, err := strconv.Atoi(id)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	res, err := h.Storage.BookApp().SubCategoryGet(sid, limit, page)

	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	h.handleResponse(c, http.Created, res)
}

// Update  		Subcategory
// @Security 	BearerAuth
// @ID	 		Subcategory
// @Router 		/subcategory [PUT]
// @Summary 	Update Subcategory
// @Description Update Subcategory
// @Tags 		Subcategory
// @Accept 		json
// @Produce 	json
// @Param 		body body models.SubCategoryUpdate true "SubCategoryUpdate"
// @Success 	200 {object} http.Response{data=models.SubCategoryRes} "Subcategory data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) UpdateSubCategory(c *gin.Context) {
	var body models.SubCategoryUpdate

	err := c.ShouldBindJSON(&body)
	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	res, err := h.Storage.BookApp().SubCategoryUpdate(&body)

	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, res)
}

// Delete 	Subcategory
// @Security 	BearerAuth
// @ID 			deletesubcategory
// @Router 		/subcategory/{id} [DELETE]
// @Summary 	Delete subcategory
// @Description Delete subcategory
// @Tags 		Subcategory
// @Accept 		json
// @Produce 	json
// @Param 		id path int true "id"
// @Success 	200 {object} http.Response{data=models.Success} "Category data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) DeleteSubCategory(c *gin.Context) {
	id := c.Param("id")
	cid, err := strconv.Atoi(id)

	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	err = h.Storage.BookApp().SubCategoryDelete(cid)

	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, "Successfully deleted")
}
