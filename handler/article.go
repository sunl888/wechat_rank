package handler

import (
	"code.aliyun.com/zmdev/wechat_rank/service"
	"errors"
	"github.com/gin-gonic/gin"
)

type Article struct {
}

func (*Article) List(ctx *gin.Context) {
	l := struct {
		WxName string `json:"wx_name" form:"wx_name"`
		Order  string `json:"order" form:"order"`
	}{}
	limit, offset := getLimitAndOffset(ctx)
	if err := ctx.ShouldBind(&l); err != nil {
		_ = ctx.Error(err)
		return
	}
	allowOrders := map[string]string{
		"latest": "published_at desc",
		"read":   "read_count desc",
		"like":   "like_count desc",
	}
	if _, ok := allowOrders[l.Order]; !ok {
		_ = ctx.Error(errors.New("不允许de排序字段值"))
		return
	}
	wexin, err := service.WechatLoad(ctx, l.WxName)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	articles, count, err := service.ArticleListWithWx(ctx, wexin.Id, allowOrders[l.Order], limit, offset)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(200, gin.H{
		"count": count,
		"data":  articles,
	})
	return
}

func NewArticle() *Article {
	return &Article{}
}
