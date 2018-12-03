package handler

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/service"
	"github.com/gin-gonic/gin"
	"time"
)

type Rank struct {
}

type RankResp struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Period    string    `json:"period"`
	StartDate string    `json:"start_date"`
	EndDate   string    `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RankDetailResp struct {
	Id           int64   `json:"id"`
	VerifyName   string  `json:"verify_name"`
	WxName       string  `json:"wx_name"`
	WxNote       string  `json:"wx_note"`
	WxLogo       string  `json:"wx_logo"`
	WxVip        string  `json:"wx_vip"`
	WxQrcode     string  `json:"wx_qrcode"`
	RankId       int64   `json:"rank_id"`
	WxId         int64   `json:"wx_id"`
	Wci          float64 `json:"wci"`
	TopCount     int64   `json:"top_count"`
	TopReadCount int64   `json:"top_read_count"`
	TopLikeCount int64   `json:"top_like_count"`
	ArticleCount int64   `json:"article_count"`
	ReadCount    int64   `json:"read_count"`
	LikeCount    int64   `json:"like_count"`
	MaxReadCount int64   `json:"max_read_count"`
	MaxLikeCount int64   `json:"max_like_count"`
	AvgReadCount int64   `json:"avg_read_count"`
}

type ArticleResp struct {
	Id           int64      `json:"id"`
	WxId         int64      `json:"wx_id"`
	WxVerifyName string     `json:"wx_verify_name"`
	Top          int64      `json:"top"`
	Title        string     `json:"title"`
	WxName       string     `json:"wx_name"`
	Url          string     `json:"url"`
	ReadCount    int64      `json:"read_count"`
	LikeCount    int64      `json:"like_count"`
	PublishedAt  *time.Time `json:"published_at"`
}

func (r *Rank) RankList(ctx *gin.Context) {
	l := struct {
		Period string `json:"period" form:"period"`
	}{}
	if err := ctx.ShouldBind(&l); err != nil {
		_ = ctx.Error(err)
		return
	}
	ranks, err := service.RankList(ctx, l.Period)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(200, convert2RanksResp(ranks))
}

func (r *Rank) AccountRank(ctx *gin.Context) {
	l := struct {
		RankId     int64 `json:"rank_id" form:"rank_id"`
		CategoryId int64 `json:"category_id" form:"category_id"`
	}{}
	limit, offset := getLimitAndOffset(ctx)
	if err := ctx.ShouldBind(&l); err != nil {
		_ = ctx.Error(err)
		return
	}
	ranks, count, err := service.RankDetail(ctx, l.RankId, l.CategoryId, limit, offset)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(200, map[string]interface{}{
		"count": count,
		"ranks": convert2RankDetailsResp(ranks),
	})
}

func (r *Rank) ArticleRank(ctx *gin.Context) {
	l := struct {
		StartDate  string `json:"start_date" form:"start_date"`
		EndDate    string `json:"end_date" form:"end_date"`
		CategoryId int64  `json:"category_id" form:"category_id"`
	}{}
	limit, offset := getLimitAndOffset(ctx)
	if err := ctx.ShouldBind(&l); err != nil {
		_ = ctx.Error(err)
		return
	}
	articles, count, err := service.ArticleRank(ctx, l.StartDate, l.EndDate, l.CategoryId, offset, limit)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(200, map[string]interface{}{
		"count": count,
		"ranks": convert2ArticlesResp(articles),
	})
}

func convert2RankDetailResp(r *model.RankJoinWechat) *RankDetailResp {
	return &RankDetailResp{
		Id:           r.Id,
		VerifyName:   r.VerifyName,
		WxName:       r.WxName,
		WxNote:       r.WxNote,
		WxLogo:       r.WxLogo,
		WxVip:        r.WxVip,
		WxQrcode:     r.WxQrcode,
		RankId:       r.RankId,
		WxId:         r.WxId,
		Wci:          r.Wci,
		TopCount:     r.TopCount,
		TopReadCount: r.TopReadCount,
		TopLikeCount: r.TopLikeCount,
		ArticleCount: r.ArticleCount,
		ReadCount:    r.ReadCount,
		LikeCount:    r.LikeCount,
		MaxReadCount: r.MaxReadCount,
		MaxLikeCount: r.MaxLikeCount,
		AvgReadCount: r.AvgReadCount,
	}
}

func convert2RankDetailsResp(rs []*model.RankJoinWechat) []*RankDetailResp {
	rankDetailsResp := make([]*RankDetailResp, 0, len(rs))
	for _, r := range rs {
		rankDetailsResp = append(rankDetailsResp, convert2RankDetailResp(r))
	}
	return rankDetailsResp
}

func convert2RankResp(r *model.Rank) *RankResp {
	return &RankResp{
		Id:        r.Id,
		Name:      r.Name,
		Period:    r.Period,
		StartDate: r.StartDate,
		EndDate:   r.EndDate,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func convert2RanksResp(rs []*model.Rank) []*RankResp {
	ranksResp := make([]*RankResp, 0, len(rs))
	for _, r := range rs {
		ranksResp = append(ranksResp, convert2RankResp(r))
	}
	return ranksResp
}

func convert2ArticleResp(a *model.Article) *ArticleResp {
	return &ArticleResp{
		Id:           a.Id,
		WxId:         a.WxId,
		WxVerifyName: a.WxVerifyName,
		Top:          a.Top,
		Title:        a.Title,
		WxName:       a.WxName,
		Url:          a.Url,
		ReadCount:    a.ReadCount,
		LikeCount:    a.LikeCount,
		PublishedAt:  a.PublishedAt,
	}
}

func convert2ArticlesResp(as []*model.Article) []*ArticleResp {
	articlesResp := make([]*ArticleResp, 0, len(as))
	for _, r := range as {
		articlesResp = append(articlesResp, convert2ArticleResp(r))
	}
	return articlesResp
}

func NewRank() *Rank {
	return &Rank{}
}
