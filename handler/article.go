package handler

import (
	"code.aliyun.com/zmdev/wechat_rank/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Article struct {
}

/*type ArticleResp struct {
	Id          int64      `json:"id"`
	WxId        int64      `json:"wx_id"`
	WxNickname  string     `json:"wx_nickname"`
	Top         int64      `json:"top"`
	Title       string     `json:"title"`
	WxName      string     `json:"wx_name"`
	Url         string     `json:"url"`
	ReadCount   int64      `json:"read_count"`
	LikeCount   int64      `json:"like_count"`
	PublishedAt *time.Time `json:"published_at"`
}*/

func (*Article) List(ctx *gin.Context) {
	l := struct {
		WxId int64 `json:"wx_id" form:"wx_id"`
	}{}
	limit, offset := getLimitAndOffset(ctx)
	if err := ctx.ShouldBind(&l); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	articles, count, err := service.ArticleListWithWx(ctx, l.WxId, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(200, gin.H{
		"count": count,
		"data":  articles,
	})
	return
}

/*func convert2ArticleResp(a *model.ArticleJoinWechat) *ArticleResp {
	return &ArticleResp{
		Id:          a.Id,
		WxId:        a.WxId,
		Top:         a.Top,
		Title:       a.Title,
		WxName:      a.WxName,
		WxNickname:  a.WxNickname,
		Url:         a.Url,
		ReadCount:   a.ReadCount,
		LikeCount:   a.LikeCount,
		PublishedAt: a.PublishedAt,
	}
}

func convert2ArticlesResp(as []*model.ArticleJoinWechat) []*ArticleResp {
	articlesResp := make([]*ArticleResp, 0, len(as))
	for _, r := range as {
		articlesResp = append(articlesResp, convert2ArticleResp(r))
	}
	return articlesResp
}*/

func NewArticle() *Article {
	return &Article{}
}
