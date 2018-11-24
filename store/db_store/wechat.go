package db_store

import (
	"github.com/jinzhu/gorm"
	"code.aliyun.com/zmdev/wechat_rank/model"
)

type dbWechat struct {
	db *gorm.DB
}

func NewDBWechat(db *gorm.DB) model.WechatStore {
	return &dbWechat{db: db}
}
