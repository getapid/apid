package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	username, password = "john.doe", "Pa55word"
	token              = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
)

func (h GinHandler) HandleLogin(c *gin.Context) {
	user, pass, ok := c.Request.BasicAuth()
	if !ok || user != username || pass != password {
		c.Status(http.StatusUnauthorized)
		return
	}

	jsonToken := struct {
		Token string `json:"token"`
	}{token}

	c.JSON(http.StatusOK, jsonToken)
}

func (h GinHandler) AuthMiddleware(c *gin.Context) {
	headerVal := c.GetHeader("Authorization")
	if headerVal != "Bearer "+token {
		response := struct {
			Reason string `json:"reason"`
		}{"provided token isn't what was issued"}

		c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
}
