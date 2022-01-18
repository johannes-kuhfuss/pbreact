package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

func Ping(c *gin.Context) {
	logger.Info("Ping invoked")
	c.String(http.StatusOK, "Pong")
}
