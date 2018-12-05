package service

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"math"
	"time"
)

type rankService struct {
	model.RankStore
	model.ArticleStore
	model.WechatStore
}

const DATE_FORMAT = "2006-01-02"

func (r *rankService) RankCreate(rank *model.Rank) error {
	err := r.RankStore.RankCreate(rank)
	if err != nil {
		return err
	}
	var (
		ranks map[string]*model.RankDetail
	)
	ranks = make(map[string]*model.RankDetail, 10)
	wechats, count, err := r.WechatStore.WechatList(0, 0)
	for i := 0; i < int(count); i++ {
		articles, err := r.ArticleStore.ArticleListByWxId(rank.StartDate, rank.EndDate, wechats[i].Id)
		if err != nil {
			return err
		}
		// 基础信息
		if ranks[wechats[i].WxName] == nil {
			ranks[wechats[i].WxName] = new(model.RankDetail)
			ranks[wechats[i].WxName].WxId = wechats[i].Id
		}
		for _, article := range articles {
			// 总文章数
			ranks[article.WxName].ArticleCount++
			// 单篇文章最高阅读数
			if article.ReadCount > ranks[article.WxName].MaxReadCount {
				ranks[article.WxName].MaxReadCount = article.ReadCount
			}
			// 单篇文章最高点赞数
			if article.LikeCount > ranks[article.WxName].MaxLikeCount {
				ranks[article.WxName].MaxLikeCount = article.LikeCount
			}
			// 所有文章阅读数
			ranks[article.WxName].ReadCount += article.ReadCount
			// 所有文章点赞数
			ranks[article.WxName].LikeCount += article.LikeCount
			if article.Top == 1 {
				// 头条文章数量
				ranks[article.WxName].TopCount++
				// 头条文章阅读数
				ranks[article.WxName].TopReadCount += article.ReadCount
				// 头条文章点赞数
				ranks[article.WxName].TopLikeCount += article.LikeCount
			}
		}
	}
	// 计算每个公众号的周期内平均阅读数
	t1, _ := time.ParseInLocation(DATE_FORMAT, rank.EndDate, time.Local)
	t2, _ := time.ParseInLocation(DATE_FORMAT, rank.StartDate, time.Local)
	days := int64(math.Abs(t1.Sub(t2).Hours() / 24))
	if days == 0 {
		days = 1
	}
	for _, rankDetail := range ranks {
		rankDetail.RankId = rank.Id
		if rankDetail.ArticleCount > 0 {
			// 计算平均阅读量
			rankDetail.AvgReadCount = int64((rankDetail.ReadCount) / days)
			// 计算得分
			rankDetail.Wci = calcScore(rankDetail, days)
		}
		err := r.RankStore.RankDetailCreate(rankDetail)
		if err != nil {
			return err
		}
	}
	return nil
}

//R为评估时间段内所有文章（n）的阅读总数；
//Z为评估时间段内所有文章（n）的点赞总数；
//d为评估时间段所含天数（一般周取7天，月度取30天，年度取365天，其他自定义时间段以真实天数计算）；
//n为评估时间段内账号所发文章数；
//Rt和Zt为评估时间段内账号所发头条的总阅读数和总点赞数；
//Rmax和Zmax为评估时间段内账号所发文章的最高阅读数和最高点赞数。
func calcScore(rank *model.RankDetail, days int64) float64 {
	// 整体传播力 o=85%*ln(R/d+1) + 15%*ln(10*Z/d+1)
	o := 0.85*math.Log(float64(rank.ReadCount/days+1)) + 0.15*math.Log(float64(10*rank.LikeCount/days+1))
	// 篇均传播力 a=85%*ln(R/n+1) + 15%*ln(10*Z/n+1)
	a := 0.85*math.Log(float64(rank.ReadCount/rank.ArticleCount+1)) + 0.15*math.Log(float64(10*rank.LikeCount/rank.ArticleCount+1))
	// 头条传播力 h=85%*ln(Rt/d+1) + 15%*ln(10*Zt/d+1)
	h := 0.85*math.Log(float64(rank.TopReadCount/days+1)) + 0.15*math.Log(float64(10*rank.TopLikeCount/days+1))
	// 峰值传播力 p=85%*ln(Rmax+1) + 15*ln(10*Zman+1)
	p := 0.85*math.Log(float64(rank.MaxReadCount+1)) + 0.15*math.Log(float64(10*rank.MaxLikeCount+1))
	// wci=(30%*o + 30%*a + 30%*h + 10%*p)^2*10
	wci := math.Pow(0.3*o+0.3*a+0.3*h+0.1*p, 2) * 10
	return wci
}

func NewRankService(rs model.RankStore, as model.ArticleStore, ws model.WechatStore) model.RankService {
	return &rankService{rs, as, ws}
}
