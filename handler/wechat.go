package handler

import (
	"code.aliyun.com/zmdev/wechat_rank/errors"
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/service"
	"github.com/gin-gonic/gin"
)

type Wechat struct{}

func (w *Wechat) Create(ctx *gin.Context) {
	l := struct {
		WxName     string `json:"wx_name" form:"wx_name"`
		CategoryId int64  `json:"category_id" form:"category_id"`
	}{}
	if err := ctx.ShouldBind(&l); err != nil {
		_ = ctx.Error(err)
		return
	}
	err := service.WechatCreate(ctx, &model.Wechat{
		WxName:     l.WxName,
		CategoryId: l.CategoryId,
	})
	if err != nil {
		_ = ctx.Error(errors.BadRequest("创建失败", nil))
		return
	}
}

func NewWechat() *Wechat {
	return &Wechat{}
}
