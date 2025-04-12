package middleware

import (
	"errors"
	"strings"

	"github.com/BingyanStudio/configIT/internal/controller"
	"github.com/BingyanStudio/configIT/internal/model"
	"github.com/BingyanStudio/configIT/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Auth(c *gin.Context) {
	appName, ok := c.Params.Get("app_name")
	if !ok {
		controller.ErrUnauthorize(c)
		c.Abort()
		return
	}
	app, err := model.GetAppByName(c, appName)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			controller.ErrNotFound(c)
		} else {
			controller.ErrInternal(c, err)
		}
		c.Abort()
		return
	}

	ip := c.ClientIP()

	pass := false

	for _, scope := range app.Scopes {
		switch scope.Scope {
		case model.AccessScopePod:
			s := strings.Split(scope.Value, "/") // namespace/pod/is_fuzzy
			ips, err := utils.GetPodIPs(c, s[0], s[1], s[2] == "true")
			if err != nil {
				controller.ErrInternal(c, err)
				c.Abort()
				return
			}
			pass = utils.Contains(ips, ip)
		case model.AccessScopeNamespace:
			ips, err := utils.GetNamespaceIPs(c, scope.Value)
			if err != nil {
				controller.ErrInternal(c, err)
				c.Abort()
				return
			}
			pass = utils.Contains(ips, ip)
		case model.AccessScopeIP:
			// 0.0.0.0/0,1.1.1.1/1
			cidrs := strings.Split(scope.Value, ",")
			for _, cidr := range cidrs {
				if utils.IPInCIDR(ip, cidr) {
					pass = true
					break
				}
			}
		case model.AccessScopeToken:
			req_token := c.GetHeader("Authorization")
			pass = req_token == scope.Value
		case model.AccessScopePublic:
			if c.Request.Method == "GET" {
				pass = true
				break
			}
		}
		if pass {
			break
		}
	}
	if !pass {
		controller.ErrUnauthorize(c)
		c.Abort()
		return
	} else {
		//TODO: log
	}
}
