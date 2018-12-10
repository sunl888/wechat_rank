package model

import "time"

// 公众号
type Wechat struct {
	Id               int64     `gorm:"primary_key" json:"id"`
	VerifyName       string    `json:"verify_name"`
	WxName           string    `gorm:"unique_index" json:"wx_name"`
	WxNote           string    `gorm:"type:varchar(255)" json:"wx_note"`
	WxNickname       string    `json:"wx_nickname"`
	WxLogo           string    `gorm:"type:varchar(255)" json:"wx_logo"`
	WxVip            string    `json:"wx_vip"`
	WxQrcode         string    `json:"wx_qrcode"`
	CategoryId       int64     `json:"category_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	LastGetArticleAt string    `json:"last_get_article_at"` // 最近一次获取文章时间
}

type WechatStore interface {
	WechatLoad(wechatName string) (wechat *Wechat, err error)
	WechatList(limit, offset int) (wechats []*Wechat, count int64, err error)
	WechatCreate(wechat *Wechat) error
	WechatUpdate(wechat *Wechat) error
	WechatListByCategory(cId int64, limit, offset int) (wechats []*Wechat, count int64, err error)
	WechatDelete(id int64) error
}

type WechatService interface {
	WechatStore
}
