package v1

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/bookapp/api/handler/v1/http"
)

// Get Book
// @ID 			getcounts
// @Security 	BearerAuth
// @Router		/statistic  [GET]
// @Summary		get parametrs count
// @Tags        Statistic
// @Description	Here get by parametrs count.
// @Accept 		json
// @Produce		json
// @Success 	200 {object} http.Response{data=models.StatisticCount} "parametrs  Count"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) GetStatistic(c *gin.Context) {
	res, err := h.Storage.BookApp().GetStatistic()
	if err != nil {
		h.handleResponse(c, http.OK, err.Error())
		return
	}
	h.handleResponse(c, http.OK, res)
}

// Get Book
// @ID 			getcategorybookcount
// @Security 	BearerAuth
// @Router		/statistic/category/bookcount  [GET]
// @Summary		get category book count
// @Tags        Statistic
// @Description	Here get by get category book count.
// @Accept 		json
// @Produce		json
// @Success 	200 {object} http.Response{data=[]models.CategoryBookCount} "category book  Count"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) GetCategoryBookCount(c *gin.Context) {
	res, err := h.Storage.BookApp().GetCategoryBookCount()
	if err != nil {
		h.handleResponse(c, http.OK, err.Error())
		return
	}
	h.handleResponse(c, http.OK, res)
}


// Get Book
// @ID 			getweekbookcount
// @Security 	BearerAuth
// @Router		/statistic/week/bookcount  [GET]
// @Summary		get week book count
// @Tags        Statistic
// @Description	Here get by get week book count.
// @Accept 		json
// @Produce		json
// @Success 	200 {object} http.Response{data=[]models.AddedBooks} "Week book  count"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) GetAddedWeekBook(c *gin.Context) {
	res, err := h.Storage.BookApp().GetWeekAddedBook()
	if err != nil {
		h.handleResponse(c, http.OK, err.Error())
		return
	}
	h.handleResponse(c, http.OK, res)
}
