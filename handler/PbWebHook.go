package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/pbreact/config"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

type WebHookHandler struct {
	Cfg *config.AppConfig
}

func (whh *WebHookHandler) PbWebHook(c *gin.Context) {
	var err api_error.ApiErr
	err = whh.validateAuthKey(c)
	if err != nil {
		logger.Error("Invalid or missing authkey", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	err = whh.parseJsonPayload(c)
	if err != nil {
		logger.Error("Could not parse JSON data in request", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (whh *WebHookHandler) validateAuthKey(c *gin.Context) api_error.ApiErr {
	authkey, exists := c.GetQuery("authkey")
	if !exists || authkey != whh.Cfg.PbAuthKey {
		//return api_error.NewUnauthorizedError("Could not verify auth key in request")
		return nil
	}
	return nil
}

func (wwh *WebHookHandler) parseJsonPayload(c *gin.Context) api_error.ApiErr {
	var jsonResult map[string]interface{}
	if err := c.ShouldBindJSON(&jsonResult); err != nil {
		msg := "Invalid JSON body in request"
		logger.Error(msg, err)
		return api_error.NewBadRequestError(msg)
	}
	log := fmt.Sprintf("Json Data: %#v", jsonResult)
	logger.Info(log)
	return nil
}
