package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/bookapp/api/handler/v1/http"
	"gitlab.com/bookapp/api/models"
	"gitlab.com/bookapp/pkg/etc"
)

// Register Admin
// @ID 			 addadmin
// @Security 	 BearerAuth
// @Router		/admin [POST]
// @Summary		add admin
// @Tags        Admin
// @Description	Here admin will be added, should enter 2 for role_id
// @Accept 		json
// @Produce		json
// @Param 		body body models.UserReq true "AddAdmin"
// @Success 	201 {object} http.Response{data=models.UserLoginRes} "Admin data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) AddAdmin(c *gin.Context) {
	var body models.UserReq
	err := c.ShouldBindJSON(&body)
	if err != nil {
		h.handleResponse(c, http.BadRequest, err)
		return
	}
	if body.RoleId != 2 {
		h.handleResponse(c, http.BadRequest, models.Message{
			Message: "role_id only accept 2",
		})
		return
	}
	login, err := h.Storage.BookApp().CheckField(&models.CheckFieldReq{
		Field: "login",
		Value: body.Login,
	})
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err)
		return
	}
	if login.Exists {
		h.handleResponse(c, http.Conflict, models.Message{
			Message: "Login already exists, Please enter another login",
		})
		return
	}
	h.JwtHandler.Sub = body.RoleId
	h.JwtHandler.Aud = []string{"bookapp"}
	h.JwtHandler.SigninKey = h.Cfg.SigningKey
	h.JwtHandler.Log = h.Log
	h.JwtHandler.Role = "admin"
	tokens, err := h.JwtHandler.GenerateAuthJWT()
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	body.RefreshToken = tokens[1]
	body.Password, err = etc.HashPassword(body.Password)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	res, err := h.Storage.BookApp().UserCreate(&body)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	result := models.UserLoginRes{
		Id:          res.Id,
		FullName:    res.FullName,
		AvatarUrl:   res.AvatarUrl,
		Login:       res.Login,
		AccessToken: tokens[0],
	}
	h.handleResponse(c, http.Created, result)
}

// Get Admins
// @ID 			 getadmins
// @Security 	 BearerAuth
// @Router		/admin [GET]
// @Summary		Get all admins
// @Tags        Admin
// @Description	Here admin will be get
// @Accept 		json
// @Produce		json
// @Success 	200 {object} http.Response{data=models.UserResForComment} "Admin data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) GetAllAdmin(c *gin.Context) {
	res, err := h.Storage.BookApp().AdminGetList()
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	h.handleResponse(c, http.OK, res)
}

// Login admin
// @ID loginadmin admin
// @Router		/admin/login  [POST]
// @Summary		login admin
// @Tags        Admin
// @Description	Here admin can login
// @Accept 		json
// @Produce		json
// @Param 		body body models.LoginReq true "LoginAdmin"
// @Success 	200 {object} http.Response{data=models.UserLoginRes} "Login"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) LoginAdmin(c *gin.Context) {
	var body models.LoginReq
	err := c.ShouldBindJSON(&body)
	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}
	loginExists, err := h.Storage.BookApp().CheckField(&models.CheckFieldReq{
		Field: "login",
		Value: body.Login,
	})
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	if !(loginExists.Exists) {
		h.handleResponse(c, http.InternalServerError, models.Message{
			Message: "Incorrect login",
		})
		return
	}
	user, err := h.Storage.BookApp().UserGetLogin(body.Login)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	if !(etc.CheckPasswordHash(body.Password, user.Password)) {
		h.handleResponse(c, http.BadRequest, models.Message{
			Message: "Incorrect password",
		})
		return
	}
	if user.RoleId == 2 {
		h.JwtHandler.Role = "admin"
	}else if user.RoleId == 1 {
		h.JwtHandler.Role = "superadmin"
	}
	h.JwtHandler.Sub = user.RoleId
	h.JwtHandler.Aud = []string{"bookapp"}
	h.JwtHandler.SigninKey = h.Cfg.SigningKey
	h.JwtHandler.Log = h.Log
	tokens, err := h.JwtHandler.GenerateAuthJWT()
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	user.AccessToken = tokens[0]

	h.handleResponse(c, http.OK, user)
}

// Delete Admin
// @ID deleteadmin
// @Router		/admin/{id} [DELETE]
// @Security 	BearerAuth
// @Summary		delete admin
// @Tags        Admin
// @Description	Here admin can delete info
// @Accept 		json
// @Produce		json
// @Param 		id  path int true "DeleteAdmin"
// @Success 	200 {object} http.Response{data=string} "Deleted"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) DeleteAdmin(c *gin.Context) {
	id := c.Param("id")
	uId, err := strconv.Atoi(id)
	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	err = h.Storage.BookApp().UserDelete(uId)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, "Admin info deleted successfully")
}

// Login superadmin
// @ID loginsuperadmin
// @Router		/superadmin/login  [POST]
// @Summary		login superadmin
// @Tags        Superadmin
// @Description	Here superadmin can login
// @Accept 		json
// @Produce		json
// @Param 		body body models.LoginReq true "LoginSuperAdmin"
// @Success 	200 {object} http.Response{data=models.UserLoginRes} "SuperadminLogin"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) LoginSuperAdmin(c *gin.Context) {
	var body models.LoginReq
	err := c.ShouldBindJSON(&body)
	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}
	loginExists, err := h.Storage.BookApp().CheckField(&models.CheckFieldReq{
		Field: "login",
		Value: body.Login,
	})
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	if !(loginExists.Exists) {
		h.handleResponse(c, http.InternalServerError, models.Message{
			Message: "Incorrect login",
		})
		return
	}
	user, err := h.Storage.BookApp().UserGetLogin(body.Login)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	if !(etc.CheckPasswordHash(body.Password, user.Password)) {
		h.handleResponse(c, http.BadRequest, models.Message{
			Message: "Incorrect password",
		})
		return
	}

	h.JwtHandler.Sub = user.RoleId
	h.JwtHandler.Aud = []string{"bookapp"}
	h.JwtHandler.SigninKey = h.Cfg.SigningKey
	h.JwtHandler.Log = h.Log
	h.JwtHandler.Role = "superadmin"
	tokens, err := h.JwtHandler.GenerateAuthJWT()
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	user.AccessToken = tokens[0]

	h.handleResponse(c, http.OK, user)
}

// Register Superadmin
// @ID 			 addsuperadmin
// @Router		/superadmin [POST]
// @Summary		add superadmin
// @Tags        Superadmin
// @Description	Here superadmin will be added, should enter 1 for role_id
// @Accept 		json
// @Produce		json
// @Param 		body body models.UserReq true "AddSuperAdmin"
// @Success 	201 {object} http.Response{data=models.UserLoginRes} "Admin data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) AddSuperAdmin(c *gin.Context) {
	var body models.UserReq
	err := c.ShouldBindJSON(&body)
	if err != nil {
		h.handleResponse(c, http.BadRequest, err)
		return
	}
	if body.RoleId != 1 {
		h.handleResponse(c, http.BadRequest, models.Message{
			Message: "role_id only accept 1",
		})
		return
	}
	login, err := h.Storage.BookApp().CheckField(&models.CheckFieldReq{
		Field: "login",
		Value: body.Login,
	})
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err)
		return
	}
	if login.Exists {
		h.handleResponse(c, http.Conflict, models.Message{
			Message: "Login already exists, Please enter another login",
		})
		return
	}
	h.JwtHandler.Sub = body.RoleId
	h.JwtHandler.Aud = []string{"bookapp"}
	h.JwtHandler.SigninKey = h.Cfg.SigningKey
	h.JwtHandler.Log = h.Log
	h.JwtHandler.Role = "superadmin"
	tokens, err := h.JwtHandler.GenerateAuthJWT()
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	body.RefreshToken = tokens[1]
	body.Password, err = etc.HashPassword(body.Password)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	res, err := h.Storage.BookApp().UserCreate(&body)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	result := models.UserLoginRes{
		Id:          res.Id,
		FullName:    res.FullName,
		AvatarUrl:   res.AvatarUrl,
		Login:       res.Login,
		AccessToken: tokens[0],
	}
	h.handleResponse(c, http.Created, result)
}
