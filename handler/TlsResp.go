package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

func TlsData(c *gin.Context) {
	logger.Info("TLS invoked")
	tlsStr := fmt.Sprintf("Hello user! Your config: %+v", c.Request.TLS)
	c.String(http.StatusOK, tlsStr)
}
