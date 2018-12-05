package model

import "time"

type Rank struct {
	Id        int64 `gorm:"primary_index"`                     // ID
	Name      string                                           // 11月26日-02日
	Period    string `gorm:"type:enum('week','month','year')"` // 时间段
	StartDate string                                           // 开始时间
	EndDate   string                                           // 结束时间
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RankDetail struct {
	Id           int64 `gorm:"primary_index"`
	RankId       int64 `gorm:"index:rank_id_wx_id_index"`
	WxId         int64 `gorm:"index:rank_id_wx_id_index"`
	Wci          float64
	TopCount     int64 // 统计周期内的总发布头条数
	TopReadCount int64 // 统计周期内所有发布内容中的头条阅读数总和
	TopLikeCount int64 // 统计周期内所有发布内容中的头条点赞数总和
	ArticleCount int64 // 统计周期内的总发布文章数 (普通文章+头条文章)
	ReadCount    int64 // 统计周期内所有发布内容的阅读数总和
	LikeCount    int64 // 统计周期内所有发布内容的点赞数总和
	MaxReadCount int64 // 统计周期内所有发布内容中的单篇最高阅读数
	MaxLikeCount int64 // 统计周期内所有发布内容中的单篇最高点赞数
	AvgReadCount int64 // 统计周期内所有发布内容的阅读数平均值
}

type RankJoinWechat struct {
	Id int64
	*RankDetail
	*Wechat
}

type RankStore interface {
	RankCreate(rank *Rank) error
	RankDetailCreate(detail *RankDetail) error
	RankList(period string) (ranks []*Rank, err error)
	RankLoad(rankId int64) (rank *Rank, err error)
	RankDetail(rankId, categoryId int64, limit, offset int) (ranks []*RankJoinWechat, count int64, err error)
}

type RankService interface {
	RankStore
	Rank(wechat *Wechat, rank *Rank) error
}
