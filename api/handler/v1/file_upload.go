package v1

import (
	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/bookapp/api/handler/v1/http"
	"gitlab.com/bookapp/api/models"
)

// @Router		/store/upload [POST]
// @Summary		Upload file
// @Tags        MediaFile
// @Description	Here admin can be logged in.
// @Accept 		multipart/form-data
// @Produce		json
// @Param       file    formData file true "File"
// @Success		200 	{object}  models.FileResponse
// @Failure     default {object}  models.DefaultResponse
func (h *handlerV1) Upload(c *gin.Context) {
	file := &models.File{}
	err := c.ShouldBind(&file)
	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	file.File.Filename = uuid.New().String() + filepath.Ext(file.File.Filename)

	err = c.SaveUploadedFile(file.File, "./store/"+file.File.Filename)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, fmt.Sprintf("Upload file: SaveUploadedFile(): %v", err.Error()))
		return
	}

	h.handleResponse(c, http.OK, models.FileResponse{
		Url: h.Cfg.MinioEnpoint + "store/" + file.File.Filename,
	})
}

func (h *handlerV1) GetFile(c *gin.Context) {
	filename := c.Param("filename")

	filePath := "./store/" + filename

	c.File(filePath)
}
