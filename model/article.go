package model

import "time"

type Article struct {
	Id          int64 `gorm:"primary_key"`   // ID
	Url         string                       // 微信地址
	Name        string                       // 微信昵称
	WxName      string                       // 微信账号
	Top         int64                        // 文章位置
	Title       string                       // 标题
	Author      string                       // 作者
	Picurl      string                       // 图片链接
	Digest      string                       // 描述
	ArticleId   string `gorm:"unique_index"` // 微信文章ID
	ReadCount   int64                        // 阅读数
	LikeCount   int64                        // 点赞数
	PublishedAt string                       // 发布时间
	OriginalUrl string                       // 原始微信链接
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ArticleStore interface {
	ArticleCreate(article *Article) error
}

type ArticleService interface {
	ArticleStore
}
