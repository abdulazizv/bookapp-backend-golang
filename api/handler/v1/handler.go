package v1

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/bookapp/api/handler/v1/http"
	"gitlab.com/bookapp/api/tokens"
	"gitlab.com/bookapp/config"
	"gitlab.com/bookapp/pkg/logger"
	"gitlab.com/bookapp/storage"
	"gitlab.com/bookapp/storage/repo"
)

type handlerV1 struct {
	Cfg        *config.Config
	Storage    storage.StorageI
	Log        logger.Logger
	JwtHandler tokens.JWTHandler
	Redis      repo.InMemoryStorageI
}

type HandlerV1Option struct {
	Cfg        *config.Config
	Storage    storage.StorageI
	Log        logger.Logger
	JwtHandler tokens.JWTHandler
	Redis      repo.InMemoryStorageI
}

func New(optoins *HandlerV1Option) *handlerV1 {
	return &handlerV1{
		Cfg:        optoins.Cfg,
		Storage:    optoins.Storage,
		Log:        optoins.Log,
		JwtHandler: optoins.JwtHandler,
		Redis:      optoins.Redis,
	}
}

func (h *handlerV1) handleResponse(c *gin.Context, status http.Status, data interface{}) {
	switch code := status.Code; {
	case code < 300:
		h.Log.Info(
			"response",
			logger.Int("code", status.Code),
			logger.String("status", status.Status),
			logger.Any("description", status.Description),
			// logger.Any("data", data),
		)
	case code < 400:
		h.Log.Info(
			"response",
			logger.Int("code", status.Code),
			logger.String("status", status.Status),
			logger.Any("description", status.Description),
			logger.Any("data", data),
		)
	default:
		h.Log.Info(
			"response",
			logger.Int("code", status.Code),
			logger.String("status", status.Status),
			logger.Any("description", status.Description),
			logger.Any("data", data),
		)
	}

	c.JSON(status.Code, http.Response{
		Status:      status.Status,
		Description: status.Description,
		Data:        data,
	})
}
