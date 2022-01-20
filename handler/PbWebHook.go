package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/pbreact/config"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

type WebHookHandler struct {
	Cfg *config.AppConfig
}

func (whh *WebHookHandler) PbWhSubscription(c *gin.Context) {
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
	authkey, exists := c.GetQuery("Authorization")
	if !exists {
		logger.Info("No auth key")
		return nil
	} else {
		logger.Info(fmt.Sprintf("Auth key: %v", authkey))
	}
	return nil
}

func (whh *WebHookHandler) PbWhEvents(c *gin.Context) {
	var eventData = make(map[string]interface{})
	err := whh.validateAuthKey(c)
	if err != nil {
		logger.Error("Invalid or missing authkey", err)
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
	c.JSON(201, nil)
}
