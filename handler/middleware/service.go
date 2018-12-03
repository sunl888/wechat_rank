package middleware

import (
	"code.aliyun.com/zmdev/wechat_rank/service"
	"github.com/gin-gonic/gin"
)

var serviceKey = "service"

func SetService(c *gin.Context, s service.Service) {
	c.Set(serviceKey, s)
}

func ServiceMiddleware(s service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		SetService(c, s)
		c.Next()
	}
}
