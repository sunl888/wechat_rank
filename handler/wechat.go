package handler

import (
	"code.aliyun.com/zmdev/wechat_rank/errors"
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Wechat struct {
}

func (w *Wechat) ListByCategory(ctx *gin.Context) {
	cId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errors.BadRequest("id 格式不正确", nil))
		return
	}
	limit, offset := getLimitAndOffset(ctx)
	wechats, count, err := service.WechatListByCategory(ctx, cId, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(200, gin.H{
		"count": count,
		"data":  wechats,
	})
}

func (w *Wechat) Show(ctx *gin.Context) {

}

func (w *Wechat) List(ctx *gin.Context) {
	limit, offset := getLimitAndOffset(ctx)
	wechats, count, err := service.WechatList(ctx, limit, offset)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(200, gin.H{
		"count": count,
		"data":  wechats,
	})
}

func (w *Wechat) Delete(ctx *gin.Context) {
	cId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		_ = ctx.Error(errors.BadRequest("id 格式不正确", nil))
		return
	}
	err = service.WechatDelete(ctx, cId)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.Status(204)
}

func (w *Wechat) Create(ctx *gin.Context) {
	l := struct {
		WxName     string `json:"wx_name" form:"wx_name"`
		CategoryId int64  `json:"category_id" form:"category_id"`
	}{}
	if err := ctx.ShouldBind(&l); err != nil {
		_ = ctx.Error(err)
		return
	}
	wechat := model.Wechat{
		WxName:     l.WxName,
		CategoryId: l.CategoryId,
	}
	err := service.WechatCreate(ctx, &wechat)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(200, wechat)
}

func NewWechat() *Wechat {
	return &Wechat{}
}
