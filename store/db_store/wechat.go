package db_store

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"github.com/jinzhu/gorm"
)

type dbWechat struct {
	db *gorm.DB
}

func (w *dbWechat) WechatUpdate(wechat *model.Wechat) error {
	err := w.db.Model(model.Wechat{}).
		Where("wx_name = ?", wechat.WxName).
		Omit("wx_name").
		Update(&wechat).
		Error
	return err
}

func (w *dbWechat) WechatList() (wechats []*model.Wechat, err error) {
	wechats = make([]*model.Wechat, 8)
	err = w.db.Find(&wechats).Error
	return
}

func (w *dbWechat) WechatLoad(wechatName string) (wechat *model.Wechat, err error) {
	wechat = &model.Wechat{}
	err = w.db.First(&wechat, "wx_name = ?", wechatName).Error
	return
}

func (w *dbWechat) WechatCreate(wechat *model.Wechat) error {
	return w.db.Create(&wechat).Error
}

func NewDBWechat(db *gorm.DB) model.WechatStore {
	return &dbWechat{db: db}
}
