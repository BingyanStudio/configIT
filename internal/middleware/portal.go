package middleware

import (
	"strings"

	"github.com/BingyanStudio/configIT/internal/auth"
	"github.com/BingyanStudio/configIT/internal/controller"
	"github.com/gin-gonic/gin"
)

func AdminAuth(c *gin.Context) {
	token := c.GetHeader("Authorization")
	token, _ = strings.CutPrefix(token, "Bearer ")
	if token == "" {
		controller.ErrUnauthorize(c)
		c.Abort()
		return
	}

	claims, err := auth.ParseToken(token)
	if err != nil {
		controller.ErrUnauthorize(c)
		c.Abort()
		return
	}
	c.Set("claims", claims)
}
