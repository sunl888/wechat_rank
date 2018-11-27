package handler

import (
	"code.aliyun.com/zmdev/wechat_rank/handler/middleware"
	"code.aliyun.com/zmdev/wechat_rank/server"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateHTTPHandler(svr *server.Server) http.Handler {

	if svr.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	wechatHandler := NewWechat()
	categoryHandler := NewCategory()

	router := gin.Default()

	router.Use(middleware.ServiceMiddleware(svr.Service))
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"hello": "world"})
	})

	router.POST("/wechat", wechatHandler.Create)

	router.POST("/category", categoryHandler.Create)
	router.DELETE("/category/:id", categoryHandler.Delete)
	router.PUT("/category/:id", categoryHandler.Update)
	router.GET("/category", categoryHandler.List)
	router.GET("/category/:id", categoryHandler.Show)

	return router
}
