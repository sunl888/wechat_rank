package handler

import (
	"code.aliyun.com/zmdev/wechat_rank/errors"
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/service"
	"github.com/gin-gonic/gin"
)

type Rank struct {
}

type RankDetailOfTypeListResp struct {
	Id          int64   `json:"id"`         // ID
	StartDate   string  `json:"start_date"` // 开始时间
	EndDate     string  `json:"end_date"`   // 结束时间
	CtegoryRank int     `json:"ctegory_rank"`
	TotalRank   int     `json:"total_rank"`
	Wci         float64 `json:"wci"`
}

type RankDetailListResp struct {
	Id                     int64   `json:"id"`             // ID
	StartDate              string  `json:"start_date"`     // 开始时间
	EndDate                string  `json:"end_date"`       // 结束时间
	TopReadCount           int64   `json:"top_read_count"` // 头条阅读数
	TopReadCountGrowthRate int64   `json:"top_read_count_growth_rate"`
	ReadCount              int64   `json:"read_count"` // 总阅读数
	ReadCountGrowthRate    int64   `json:"read_count_growth_rate"`
	LikeCount              int64   `json:"like_count"` // 点赞数
	LikeCountGrowthRate    int64   `json:"like_count_growth_rate"`
	AvgReadCount           int64   `json:"avg_read_count"` // 平均阅读数
	AvgReadCountGrowthRate int64   `json:"avg_read_count_growth_rate"`
	Wci                    float64 `json:"wci"`
	WciGrowthRate          float64 `json:"wci_growth_rate"`
}

func (r *Rank) RankList(ctx *gin.Context) {
	l := struct {
		Period string `json:"period" form:"period"`
	}{}
	if err := ctx.ShouldBind(&l); err != nil {
		_ = ctx.Error(errors.BindError(err))
		return
	}
	ranks, err := service.RankList(ctx, l.Period)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(200, ranks)
}

func (r *Rank) RankChartWithTypes(ctx *gin.Context) {
	l := struct {
		WxName string `json:"wx_name" form:"wx_name"`
	}{}
	if err := ctx.ShouldBind(&l); err != nil {
		_ = ctx.Error(err)
		return
	}
	var rankIds []int64
	var rankMap map[int64]*model.Rank
	wexin, err := service.WechatLoad(ctx, l.WxName)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	// limit = 5
	ranks, err := service.RankList(ctx, "week")
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	rankMap = make(map[int64]*model.Rank, len(ranks))
	for _, r := range ranks {
		rankIds = append(rankIds, r.Id)
		rankMap[r.Id] = r
	}
	// 指定分类
	rankDetailList, err := service.RankDetailListByRankIds(ctx, rankIds, 0, wexin.CategoryId)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(200, gin.H{
		"data": convert2DetailOfTypeListResp(wexin, rankMap, rankDetailList),
	})
}

func (r *Rank) RankChart(ctx *gin.Context) {
	var (
		rankIds []int64
		rankMap map[int64]*model.Rank
	)
	l := struct {
		WxName string `json:"wx_name" form:"wx_name"`
	}{}
	if err := ctx.ShouldBind(&l); err != nil {
		_ = ctx.Error(err)
		return
	}
	wexin, err := service.WechatLoad(ctx, l.WxName)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ranks, err := service.RankList(ctx, "week")
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	rankMap = make(map[int64]*model.Rank, len(ranks))
	for _, r := range ranks {
		rankIds = append(rankIds, r.Id)
		rankMap[r.Id] = r
	}
	rankDetailList, err := service.RankDetailListByRankIds(ctx, rankIds, wexin.Id, 0)
	ctx.JSON(200, gin.H{
		"data": convert2DetailListResp(rankMap, rankDetailList),
	})
}

func (r *Rank) AccountRank(ctx *gin.Context) {
	l := struct {
		RankId     int64 `json:"rank_id" form:"rank_id"`
		CategoryId int64 `json:"category_id" form:"category_id"`
	}{}
	limit, offset := getLimitAndOffset(ctx)
	if err := ctx.ShouldBind(&l); err != nil {
		_ = ctx.Error(errors.BindError(err))
		return
	}
	ranks, count, err := service.RankDetail(ctx, l.RankId, l.CategoryId, limit, offset)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(200, gin.H{
		"count": count,
		"ranks": ranks,
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
		_ = ctx.Error(errors.BindError(err))
		return
	}
	articles, count, err := service.ArticleRank(ctx, l.StartDate, l.EndDate, l.CategoryId, offset, limit)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	ctx.JSON(200, gin.H{
		"count": count,
		"data":  articles,
	})
}

func convert2DetailListResp(rankMap map[int64]*model.Rank, details []*model.RankDetailAndWechat) []*RankDetailListResp {
	var detailListResp []*RankDetailListResp
	detailListResp = make([]*RankDetailListResp, len(details))
	for k, v := range details {
		detailListResp[k] = &RankDetailListResp{
			Id:           v.Id,
			StartDate:    rankMap[v.RankId].StartDate,
			EndDate:      rankMap[v.RankId].EndDate,
			TopReadCount: v.TopReadCount,
			ReadCount:    v.ReadCount,
			LikeCount:    v.LikeCount,
			AvgReadCount: v.AvgReadCount,
			Wci:          v.Wci,
		}
		if k > 0 {
			detailListResp[k].TopReadCountGrowthRate = detailListResp[k].TopReadCount - detailListResp[k-1].TopReadCount
			detailListResp[k].LikeCountGrowthRate = detailListResp[k].LikeCount - detailListResp[k-1].LikeCount
			detailListResp[k].ReadCountGrowthRate = detailListResp[k].ReadCount - detailListResp[k-1].ReadCount
			detailListResp[k].AvgReadCountGrowthRate = detailListResp[k].AvgReadCount - detailListResp[k-1].AvgReadCount
			detailListResp[k].WciGrowthRate = detailListResp[k].Wci - detailListResp[k-1].Wci
		}
	}
	return detailListResp
}

func convert2DetailOfTypeListResp(weixin *model.Wechat, rankMap map[int64]*model.Rank, details []*model.RankDetailAndWechat) []*RankDetailOfTypeListResp {
	var (
		rankIndex      map[int64]int
		detailListResp []*RankDetailOfTypeListResp
	)
	rankIndex = make(map[int64]int, len(rankMap))
	detailListResp = make([]*RankDetailOfTypeListResp, 0, len(rankMap))
	for _, v := range details {
		if v.CategoryId == weixin.CategoryId {
			rankIndex[v.RankId]++
			if v.WxId == weixin.Id {
				detailListResp = append(detailListResp, &RankDetailOfTypeListResp{
					Id:          v.Id,
					StartDate:   rankMap[v.RankId].StartDate,
					EndDate:     rankMap[v.RankId].EndDate,
					CtegoryRank: rankIndex[v.RankId],
					TotalRank:   v.TotalRank,
					Wci:         v.Wci,
				})
			}
		}
	}
	return detailListResp
}

func NewRank() *Rank {
	return &Rank{}
}
