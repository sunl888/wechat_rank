package handler

import (
	"code.aliyun.com/zmdev/wechat_rank/errors"
	"code.aliyun.com/zmdev/wechat_rank/model"
	"github.com/gin-gonic/gin"
)

type Wechat struct {
	ws model.WechatService
}

func (w *Wechat) Create(ctx *gin.Context) {
	l := struct {
		Name string `json:"name" form:"name"`
	}{}
	if err := ctx.ShouldBind(&l); err != nil {
		_ = ctx.Error(err)
		return
	}
	err := w.ws.WechatCreate(&model.Wechat{
		Name: l.Name,
	})
	if err != nil {
		_ = ctx.Error(errors.BadRequest("创建失败", nil))
		return
	}
}

func NewWechat() *Wechat {
	return &Wechat{}
}
