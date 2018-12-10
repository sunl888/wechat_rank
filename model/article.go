package model

import "time"

type Article struct {
	Id          int64      `gorm:"primary_key" json:"id"`          // ID
	WxId        int64      `json:"wx_id"`                          // 公众号id
	Desc        string     `json:"desc"`                           // 描述(摘要)
	Top         int64      `json:"top"`                            // 文章位置
	Title       string     `json:"title"`                          // 标题
	Url         string     `json:"url"`                            // 文章url
	ArticleId   string     `gorm:"unique_index" json:"article_id"` // 微信文章ID
	ReadCount   int64      `json:"read_count"`                     // 阅读数
	LikeCount   int64      `json:"like_count"`                     // 点赞数
	PublishedAt *time.Time `json:"published_at"`                   // 发布时间
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ArticleJoinWechat struct {
	*Article
	WxNickname string `json:"wx_nickname"`
	WxName     string `json:"wx_name"`
}

type ArticleStore interface {
	ArticleCreate(article *Article) error
	// 显示指定日期区间所有文章
	ArticleList(startDate, endDate string, offset, limit int) ([]*Article, error)
	// 通过微信id查找文章列表 需要传日期区间
	ArticleListByWxId(startDate, endDate string, wxId int64) ([]*Article, error)
	// 通过微信id查找文章列表  不需要日期
	ArticleListWithWx(wxId int64, order string, offset, limit int) (articles []*ArticleJoinWechat, count int64, err error)
	ArticleRank(startDate, endDate string, categoryId int64, offset, limit int) ([]*ArticleJoinWechat, int64, error)
	ArticleSearch(keyword string, order string, categoryId int64, offset, limit int) (articles []*ArticleJoinWechat, count int64, err error)
}

type ArticleService interface {
	ArticleStore
	ArticleGrab(wechat *Wechat, laskWeekStartDate, laskWeekEndDate string) error
}
