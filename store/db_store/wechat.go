package db_store

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"github.com/jinzhu/gorm"
)

type dbWechat struct {
	db *gorm.DB
}

func (w *dbWechat) WechatLoad(wechatName string) (wechat *model.Wechat, err error) {
	wechat = &model.Wechat{}
	err = w.db.First(&wechat, "wx_name = ?", wechatName).Error
	return
}

func (w *dbWechat) WechatCreate(wechat *model.Wechat) error {
	return w.db.FirstOrCreate(&wechat, "wx_name = ?", wechat.WxName).Error
}

func NewDBWechat(db *gorm.DB) model.WechatStore {
	return &dbWechat{db: db}
}
