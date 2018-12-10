package handler

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/service"
	"errors"
	"github.com/gin-gonic/gin"
	"time"
)

type Article struct {
}

const SHORTDATE = "2006-01-02"

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

func (*Article) Glab(c *gin.Context) {
	l := struct {
		WxName string `json:"wx_name" form:"wx_name"`
	}{}
	if err := c.ShouldBind(&l); err != nil {
		_ = c.Error(err)
		return
	}
	var (
		wechat    *model.Wechat
		wechats   []*model.Wechat
		err       error
		now       time.Time
		startDate string
		yesterday string
	)
	now = time.Now()
	yesterday = now.AddDate(0, 0, -1).Format(SHORTDATE)
	if l.WxName == "" {
		wechats, _, err = service.WechatList(c, 0, 0)
		if err != nil {
			_ = c.Error(err)
			return
		}
	} else {
		wechat, err = service.WechatLoad(c, l.WxName)
		if err != nil {
			_ = c.Error(err)
			return
		}
		wechats = append(wechats, wechat)
	}
	for i := 0; i < len(wechats); i++ {
		if wechats[i].LastGetArticleAt == "" {
			// seven days ago
			startDate = now.AddDate(0, 0, -7).Format(SHORTDATE)
			//startDate = "2018-10-01"
		} else {
			startDate = wechats[i].LastGetArticleAt
		}
		err = service.ArticleGrab(c, wechats[i], startDate, yesterday)
		if err != nil {
			if err != nil {
				_ = c.Error(err)
				return
			}
		}
		if (i+1)%10 == 0 {
			time.Sleep(1100 * time.Millisecond)
		}
	}
	return
}
func NewArticle() *Article {
	return &Article{}
}
