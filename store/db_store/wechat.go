package db_store

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"github.com/jinzhu/gorm"
)

type dbWechat struct {
	db *gorm.DB
}

func (w *dbWechat) WechatCreate(wechat *model.Wechat) error {
	return w.db.Create(&wechat).Error
}

func NewDBWechat(db *gorm.DB) model.WechatStore {
	return &dbWechat{db: db}
}
