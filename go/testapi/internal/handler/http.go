package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinHandler struct {
}

func NewGinHandler() GinHandler {
	return GinHandler{}
}

func (h GinHandler) HandleHealthCheck(c *gin.Context) {
	c.Data(http.StatusOK, "text/plain", []byte("OK"))
}

// HandleEcho returns the whole body as received with the same Content-Type header as the one received
func (h GinHandler) HandleEcho(c *gin.Context) {
	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	contentType := c.GetHeader("Content-Type")
	c.Data(http.StatusOK, contentType, payload)
}
