package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/pbreact/config"
	"github.com/johannes-kuhfuss/pbreact/service"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

type WebHookHandler struct {
	Cfg          *config.AppConfig
	PbApiService *service.PbApiService
}

func (whh *WebHookHandler) PbWhSubscription(c *gin.Context) {
	err := whh.validateAuthKey(c)
	if err != nil {
		logger.Error("Could not handle subscription response", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	id, _ := c.GetQuery("validationToken")
	c.String(200, id)
}

func (whh *WebHookHandler) validateAuthKey(c *gin.Context) api_error.ApiErr {
	authKey := c.GetHeader("Authorization")
	if authKey != whh.Cfg.RunTime.CallbackAuthToken {
		return api_error.NewUnauthenticatedError("wrong or missing auth key")
	}
	return nil
}

func (whh *WebHookHandler) PbWhEvents(c *gin.Context) {
	var eventData = make(map[string]interface{})
	err := whh.validateAuthKey(c)
	if err != nil {
		logger.Error("Could not handle event notification", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	if err := c.ShouldBindJSON(&eventData); err != nil {
		logger.Error("invalid JSON body in event notification", err)
		apiErr := api_error.NewBadRequestError("invalid json body")
		c.JSON(apiErr.StatusCode(), apiErr)
		return
	}
	log := fmt.Sprintf("Event data: %#v", eventData)
	logger.Info(log)
	c.JSON(http.StatusNoContent, nil)
}

func (whh *WebHookHandler) Register(c *gin.Context) {
	err := whh.PbApiService.RegisterForNotifications()
	if err != nil {
		logger.Error("Could not register for event notification", err)
		c.JSON(err.StatusCode(), err)
		return
	}
	c.JSON(http.StatusOK, nil)
}
