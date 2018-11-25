package model

import "time"

type Article struct {
	Id          int64  `gorm:"primary_key"`
	WxId        string `gorm:"unique_index"`
	Title       string `gorm:"type:varchar(255)"`
	Desc        string `gorm:"type:varchar(255)"`
	Top         int64
	Url         string `gorm:"type:varchar(255)"`
	ReadCount   int64
	LikeCount   int64
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ArticleStore interface {
}

type ArticleService interface {
	ArticleStore
}
