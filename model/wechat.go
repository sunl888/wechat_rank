package model

import "time"

type Wechat struct {
	Id         int64 `gorm:"primary_key"`
	Title      string
	Name       string `gorm:"unique_index"`
	Desc       string `gorm:"type:varchar(255)"`
	Logo       string `gorm:"type:varchar(255)"`
	Vip        string
	CategoryId int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type WechatStore interface {
	WechatCreate(wechat *Wechat) error
}

type WechatService interface {
	WechatStore
}
