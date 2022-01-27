package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/pbreact/config"
	"github.com/johannes-kuhfuss/pbreact/dto"
	"github.com/johannes-kuhfuss/pbreact/service"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

type WebHookHandler struct {
	Cfg          *config.AppConfig
	PbApiService *service.PbApiService
}

func NewWebHookHandler(cfg *config.AppConfig, service service.PbApiService) WebHookHandler {
	return WebHookHandler{
		Cfg:          cfg,
		PbApiService: &service,
	}
}

func (whh *WebHookHandler) PbWhSubscription(c *gin.Context) {
	err := whh.validateAuthKey(c)
	if err != nil {
		logger.Error("Could not handle subscription response", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	id, exists := c.GetQuery("validationToken")
	if !exists {
		msg := "Could not find validation token"
		err := api_error.NewBadRequestError(msg)
		logger.Error(msg, nil)
		c.JSON(err.StatusCode(), err)
	}
	c.String(200, id)
}

func (whh *WebHookHandler) validateAuthKey(c *gin.Context) api_error.ApiErr {
	authKey := c.GetHeader("Authorization")
	if (authKey == "") || (authKey != whh.Cfg.RunTime.CallbackAuthToken) {
		return api_error.NewUnauthenticatedError("Wrong or missing auth key")
	}
	return nil
}

func (whh *WebHookHandler) PbWhEvents(c *gin.Context) {
	var eventData = dto.PbEventNotification{}

	err := whh.validateAuthKey(c)
	if err != nil {
		logger.Error("Could not handle event notification", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	if err := c.ShouldBindJSON(&eventData); err != nil {
		logger.Error("Invalid JSON body in event notification", err)
		apiErr := api_error.NewBadRequestError("Invalid json body")
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	// queue element for processing
	c.JSON(http.StatusNoContent, nil)
}
