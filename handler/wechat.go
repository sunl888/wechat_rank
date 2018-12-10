package handler

import (
	"code.aliyun.com/zmdev/wechat_rank/errors"
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/service"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

type Wechat struct {
}

type WechatResp struct {
	Id           int64     `json:"id"`
	VerifyName   string    `json:"verify_name"`
	WxName       string    `json:"wx_name"`
	CategoryId   int64     `json:"category_id"`
	CategoryName string    `json:"category_name"`
	WxNote       string    `json:"wx_note"`
	WxNickname   string    `json:"wx_nickname"`
	WxLogo       string    `json:"wx_logo"`
	WxVip        string    `json:"wx_vip"`
	WxQrcode     string    `json:"wx_qrcode"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (w *Wechat) Search(ctx *gin.Context) {
	var (
		sType   string
		wechats []*model.WechatAndCategory
		count   int64
		err     error
	)
	sType = strings.Trim(ctx.Param("type"), " ")
	l := struct {
		Keyword    string `json:"keyword" form:"keyword"`
		Order      string `json:"order" form:"order"`
		CategoryId int64  `json:"category_id" form:"category_id"`
	}{}
	if err = ctx.ShouldBind(&l); err != nil {
		_ = ctx.Error(err)
		return
	}
	limit, offset := getLimitAndOffset(ctx)
	switch sType {
	case "article":
		allowOrders := map[string]string{
			"latest": "published_at desc",
			"read":   "read_count desc",
			"like":   "like_count desc",
		}
		if l.Order == "" {
			l.Order = allowOrders["latest"]
		} else if _, ok := allowOrders[l.Order]; !ok {
			_ = ctx.Error(errors.BadRequest("不允许de排序字段值", nil))
			return
		}
		l.Order = allowOrders[l.Order]
		articles, count, err := service.ArticleSearch(ctx, l.Keyword, l.Order, l.CategoryId, offset, limit)
		if err != nil {
			_ = ctx.Error(err)
			return
		}
		ctx.JSON(200, gin.H{
			"count": count,
			"data":  articles,
		})
	case "wechat":
		wechats, count, err = service.WechatSearch(ctx, l.Keyword, limit, offset)
		if err != nil {
			_ = ctx.Error(err)
			return
		}
		ctx.JSON(200, gin.H{
			"count": count,
			"data":  convert2WechatsResp(wechats),
		})
	default:
		_ = ctx.Error(errors.BadRequest("不允许的搜索类型", nil))
		return
	}
}

func (w *Wechat) ListByCategory(ctx *gin.Context) {
	cId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		_ = ctx.Error(errors.BadRequest("id 格式不正确", nil))
		return
	}
	limit, offset := getLimitAndOffset(ctx)
	wechats, count, err := service.WechatListByCategory(ctx, cId, limit, offset)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(200, gin.H{
		"count": count,
		"data":  wechats,
	})
}

func (w *Wechat) Show(ctx *gin.Context) {
	l := struct {
		WxName string `json:"wx_name" form:"wx_name"`
	}{}
	if err := ctx.ShouldBind(&l); err != nil {
		_ = ctx.Error(err)
		return
	}
	wechat, err := service.WechatLoad(ctx, l.WxName)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	category, err := service.CategoryLoad(ctx, wechat.CategoryId)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(200, convert2WechatResp(wechat, category))
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

func convert2WechatResp(wechat *model.Wechat, category *model.Category) *WechatResp {
	return &WechatResp{
		Id:           wechat.Id,
		VerifyName:   wechat.VerifyName,
		WxName:       wechat.WxName,
		WxNote:       wechat.WxNote,
		WxNickname:   wechat.WxNickname,
		WxLogo:       wechat.WxLogo,
		WxVip:        wechat.WxVip,
		WxQrcode:     wechat.WxQrcode,
		CategoryId:   category.Id,
		CategoryName: category.Title,
		CreatedAt:    wechat.CreatedAt,
		UpdatedAt:    wechat.UpdatedAt,
	}
}

func convert2WechatsResp(wechats []*model.WechatAndCategory) [] *WechatResp {
	wechatList := make([]*WechatResp, 0, len(wechats))
	for _, wechat := range wechats {
		wechatList = append(wechatList, &WechatResp{
			Id:           wechat.Id,
			VerifyName:   wechat.VerifyName,
			WxName:       wechat.WxName,
			WxNote:       wechat.WxNote,
			WxNickname:   wechat.WxNickname,
			WxLogo:       wechat.WxLogo,
			WxVip:        wechat.WxVip,
			WxQrcode:     wechat.WxQrcode,
			CategoryId:   wechat.CategoryId,
			CategoryName: wechat.CategoryName,
			CreatedAt:    wechat.CreatedAt,
			UpdatedAt:    wechat.UpdatedAt,
		})
	}
	return wechatList
}

func NewWechat() *Wechat {
	return &Wechat{}
}
