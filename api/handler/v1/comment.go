package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/bookapp/api/handler/v1/http"
	"gitlab.com/bookapp/api/models"
)

// @Router		/comment/post [POST]
// @Summary		create comment
// @Security 	BearerAuth
// @Tags        Comment
// @Description	Here comments will be created.
// @Accept 		json
// @Produce		json
// @Param       body    body models.CommentReq true "Create"
// @Success		201 	{object}  models.CommentRes
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) CreateComment(c *gin.Context) {
	body := models.CommentReq{}

	err := c.ShouldBindJSON(&body)

	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}
	res, err := h.Storage.BookApp().CommentCreate(&body)

	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.Created, res)
}

// @Router		/comment/put [PUT]
// @Summary		update comment
// @Security 	BearerAuth
// @Tags        Comment
// @Description	Here comment will be updated.
// @Accept 		json
// @Produce		json
// @Param       body    body models.CommentUpdate true "Update"
// @Success		200 	{object}  models.CommentRes
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) UpdateComment(c *gin.Context) {
	body := models.CommentUpdate{}

	err := c.ShouldBindJSON(&body)

	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}
	res, err := h.Storage.BookApp().CommentUpdate(&body)

	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, res)
}

// @Router		/comment/{id} [DELETE]
// @Summary		delete comment
// @Security 	BearerAuth
// @Tags        Comment
// @Description	Here comment will be delete.
// @Accept 		json
// @Produce		json
// @Param      	id 		path int  true "Delete"
// @Success		200 	{object}  http.Response{data=string} "Deleted"
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) DeleteComment(c *gin.Context) {
	id := c.Param("id")
	cId, err := strconv.Atoi(id)
	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	err = h.Storage.BookApp().CommentDelete(cId)

	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, "Comment deleted successfully")
}
