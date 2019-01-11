package model

import "time"

type Rank struct {
	Id        int64     `gorm:"primary_index" json:"id"`                        // ID
	Name      string    `json:"name"`                                           // 11月26日-02日
	Period    string    `gorm:"type:enum('week','month','year')" json:"period"` // 时间段
	StartDate string    `json:"start_date"`                                     // 开始时间
	EndDate   string    `json:"end_date"`                                       // 结束时间
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RankDetail struct {
	Id           int64   `gorm:"primary_index" json:"id"`                  //id
	RankId       int64   `gorm:"index:rank_id_wx_id_index" json:"rank_id"` // 排行榜id
	WxId         int64   `gorm:"index:rank_id_wx_id_index" json:"wx_id"`   // 微信id
	Wci          float64 `json:"wci"`                                      // 得分
	TopCount     int64   `json:"top_count"`                                // 统计周期内的总发布头条数
	TopReadCount int64   `json:"top_read_count"`                           // 统计周期内所有发布内容中的头条阅读数总和
	TopLikeCount int64   `json:"top_like_count"`                           // 统计周期内所有发布内容中的头条点赞数总和
	ArticleCount int64   `json:"article_count"`                            // 统计周期内的总发布文章数 (普通文章+头条文章)
	ReadCount    int64   `json:"read_count"`                               // 统计周期内所有发布内容的阅读数总和
	LikeCount    int64   `json:"like_count"`                               // 统计周期内所有发布内容的点赞数总和
	MaxReadCount int64   `json:"max_read_count"`                           // 统计周期内所有发布内容中的单篇最高阅读数
	MaxLikeCount int64   `json:"max_like_count"`                           // 统计周期内所有发布内容中的单篇最高点赞数
	AvgReadCount int64   `json:"avg_read_count"`                           // 统计周期内所有发布内容的阅读数平均值
	AvgLikeCount int64   `json:"avg_like_count"`                           // 统计周期内所有发布内容的阅读数平均值
	LikeRate     float64 `json:"like_rate"`                                // 点赞率
	TotalRank    int     `json:"total_rank"`                               // 总排名
}

type RankDetailAndWechat struct {
	Id         int64  `json:"id"`
	WxName     string `json:"wx_name"`
	WxNickname string `json:"wx_nickname"`
	CategoryId int64  `json:"category_id"`
	WxLogo     string `json:"wx_logo"`
	WxVip      string `json:"wx_vip"`
	WxQrcode   string `json:"wx_qrcode"`
	*RankDetail
}

type RankStore interface {
	RankCreate(rank *Rank) error
	RankAllList() (ranks []*Rank, err error)
	RankDetailCreate(detail *RankDetail) error
	RankList(period string) (ranks []*Rank, err error)
	RankLoad(rankId int64) (rank *Rank, err error)
	RankDetail(rankId, categoryId int64, limit, offset int) (ranks []*RankDetailAndWechat, count int64, err error)
	RankDetailListByRankIds(rankIds []int64, wxId, categoryId int64) (ranks []*RankDetailAndWechat, err error)
}

type RankService interface {
	RankStore
	Rank(wechat *Wechat, rank *Rank) (*RankDetail, error)
}
