package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/bookapp/api/handler/v1/http"
	"gitlab.com/bookapp/api/models"
	"gitlab.com/bookapp/pkg/etc"
)

// Register User
// @ID 			register user
// @Router		/client/register [POST]
// @Summary		register user
// @Tags        User
// @Description	Here user can register
// @Accept 		json
// @Produce		json
// @Param 		body body models.UserReq true "RegisterUser"
// @Success 	201 {object} http.Response{data=models.UserLoginRes} "User data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) RegisterUser(c *gin.Context) {
	var body models.UserReq
	err := c.ShouldBindJSON(&body)
	if err != nil {
		h.handleResponse(c, http.BadRequest, err)
		return
	}
	if body.RoleId != 3 {
		h.handleResponse(c, http.BadRequest, models.Message{
			Message: "role_id only accept 3",
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
	h.JwtHandler.Role = "user"
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

// Get User
// @ID getuser user
// @Router		/client/{id} [GET]
// @Security 	BearerAuth
// @Summary		get user
// @Tags        User
// @Description	Here user can get
// @Accept 		json
// @Produce		json
// @Param 		id path int  true "GetUser"
// @Success 	200 {object} http.Response{data=models.UserRes} "User data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) GetUser(c *gin.Context) {
	id := c.Param("id")
	uId, err := strconv.Atoi(id)
	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}

	res, err := h.Storage.BookApp().UserGet(uId)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, res)
}

// Update User
// @ID updateuser user
// @Router		/client/update [PUT]
// @Security 	BearerAuth
// @Summary		update user
// @Tags        User
// @Description	Here user can update info
// @Accept 		json
// @Produce		json
// @Param 		body body 	 models.UserUpdateReq  true "UpdateUser"
// @Success 	200 {object} http.Response{data=models.UserRes} "User data"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) UpdateUser(c *gin.Context) {
	var body models.UserUpdateReq

	err := c.ShouldBindJSON(&body)
	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}
	body.Password, err = etc.HashPassword(body.Password)
	if err != nil {
		h.handleResponse(c, http.BadRequest, err.Error())
		return
	}
	res, err := h.Storage.BookApp().UserUpdate(&body)
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}

	h.handleResponse(c, http.OK, res)
}

// Delete User
// @ID deleteuser user
// @Router		/client/{id} [DELETE]
// @Security 	BearerAuth
// @Summary		delete user
// @Tags        User
// @Description	Here user can delete info
// @Accept 		json
// @Produce		json
// @Param 		id  path int true "DeleteUser"
// @Success 	200 {object} http.Response{data=string} "Deleted"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) DeleteUser(c *gin.Context) {
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

	h.handleResponse(c, http.OK, "User info deleted successfully")
}

// Login User
// @ID loginuser user
// @Router		/client/login  [POST]
// @Summary		login user
// @Tags        User
// @Description	Here user can login
// @Accept 		json
// @Produce		json
// @Param 		body body models.LoginReq true "LoginUser"
// @Success 	200 {object} http.Response{data=models.UserLoginRes} "Login"
// @Response 	400 {object} http.Response{data=string} "Bad Request"
// @Failure 	500 {object} http.Response{data=string} "Server Error"
func (h *handlerV1) LoginUser(c *gin.Context) {
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
	h.JwtHandler.Role = "user"
	tokens, err := h.JwtHandler.GenerateAuthJWT()
	if err != nil {
		h.handleResponse(c, http.InternalServerError, err.Error())
		return
	}
	user.AccessToken = tokens[0]

	h.handleResponse(c, http.OK, user)

}
