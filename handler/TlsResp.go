package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TlsData(c *gin.Context) {
	tlsStr := fmt.Sprintf("Hello user! Your config: %+v", c.Request.TLS)
	c.String(http.StatusOK, tlsStr)
}
