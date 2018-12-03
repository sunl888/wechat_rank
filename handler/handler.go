package handler

import (
	"code.aliyun.com/zmdev/wechat_rank/handler/middleware"
	"code.aliyun.com/zmdev/wechat_rank/server"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreateHTTPHandler(svr *server.Server) http.Handler {

	if svr.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	wechatHandler := NewWechat()
	categoryHandler := NewCategory()
	rankHandler := NewRank()

	router := gin.Default()
	router.Use(middleware.ServiceMiddleware(svr.Service))

	// 添加公众号
	router.POST("/wechat", wechatHandler.Create)
	// 创建分类
	router.POST("/category", categoryHandler.Create)
	router.DELETE("/category/:id", categoryHandler.Delete)
	router.PUT("/category/:id", categoryHandler.Update)
	router.GET("/category", categoryHandler.List)
	router.GET("/category/:id", categoryHandler.Show)
	// 获取日期
	router.GET("/rank/date", rankHandler.RankList)
	// 公众号排名
	router.GET("/rank/account", rankHandler.AccountRank)
	// 文章排名(日期区间随意)
	router.GET("/rank/article", rankHandler.ArticleRank)
	return router
}

func getLimitAndOffset(c *gin.Context) (limit, offset int) {
	var err error
	limit, err = strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}
	if limit > 200 {
		limit = 200
	}
	offset, err = strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 0
	}
	return limit, offset
}
