package middleware

import (
	"code.aliyun.com/zmdev/wechat_rank/service"
	"github.com/gin-gonic/gin"
)

var serviceKey = "service"

func SetService(c *gin.Context, s service.Service) {
	c.Set(serviceKey, s)
}

// TODO
/*func GetService(c *gin.Context) (s service.Service, exists bool) {
	s, exists = c.Get(serviceKey)
	return
}*/

func ServiceMiddleware(s service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		SetService(c, s)
		c.Next()
	}
}
