package model

import "time"

// 公众号
type Wechat struct {
	Id         int64 `gorm:"primary_key"`
	VerifyName string
	WxName     string `gorm:"unique_index"`
	WxNote     string `gorm:"type:varchar(255)"`
	WxLogo     string `gorm:"type:varchar(255)"`
	WxVip      string
	WxQrcode   string
	CategoryId int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type WechatStore interface {
	WechatLoad(wechatName string) (wechat *Wechat, err error)
	WechatList() (wechats []*Wechat, err error)
	WechatCreate(wechat *Wechat) error
	WechatUpdate(wechat *Wechat) error
}

type WechatService interface {
	WechatStore
}
