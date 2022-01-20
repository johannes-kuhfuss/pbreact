package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/pbreact/config"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

type WebHookHandler struct {
	Cfg *config.AppConfig
}

func (whh *WebHookHandler) PbWebHook(c *gin.Context) {
	err := whh.validateAuthKey(c)
	if err != nil {
		logger.Error("Invalid or missing authkey", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	id, _ := c.GetQuery("validationToken")
	c.String(200, id)
}

func (whh *WebHookHandler) validateAuthKey(c *gin.Context) api_error.ApiErr {
	authkey, exists := c.GetQuery("authkey")
	if !exists || authkey != whh.Cfg.PbAuthKey {
		//return api_error.NewUnauthorizedError("Could not verify auth key in request")
		return nil
	}
	return nil
}
