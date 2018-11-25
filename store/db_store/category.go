package db_store

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"github.com/jinzhu/gorm"
)

type dbCategory struct {
	db *gorm.DB
}

func NewDBCategory(db *gorm.DB) model.CategoryStore {
	return &dbCategory{db: db}
}
