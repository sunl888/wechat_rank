package model

import "time"

type Article struct {
	Id           int64 `gorm:"primary_key"` // ID
	WxId         int64                      // 公众号id
	WxVerifyName string
	WxCategoryId int64
	Top          int64                        // 文章位置
	Title        string                       // 标题
	WxName       string                       // 微信账号
	Url          string                       // 文章url
	ArticleId    string `gorm:"unique_index"` // 微信文章ID
	ReadCount    int64                        // 阅读数
	LikeCount    int64                        // 点赞数
	PublishedAt  *time.Time                   // 发布时间
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ArticleStore interface {
	ArticleCreate(article *Article) error
	ArticleList(startDate, endDate string, offset, limit int) ([]*Article, int64, error)
	ArticleRank(startDate, endDate string, categoryId int64, offset, limit int) ([]*Article, int64, error)
}

type ArticleService interface {
	ArticleStore
	ArticleGrab(laskWeekStartDate, laskWeekEndDate string) error
}
