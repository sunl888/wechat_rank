package handler

import (
	"code.aliyun.com/zmdev/wechat_rank/server"
	"github.com/gin-gonic/gin"
	"net/http"
	"code.aliyun.com/zmdev/wechat_rank/handler/middleware"
)

func CreateHTTPHandler(svr *server.Server) http.Handler {
	svc := server.SetupService(svr)

	if svr.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.Use(middleware.ServiceMiddleware(svc))
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"hello": "world"})
	})

	return router
}
