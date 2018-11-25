package db_store

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"github.com/jinzhu/gorm"
)

type dbArticle struct {
	db *gorm.DB
}

func NewDBArticle(db *gorm.DB) model.ArticleStore {
	return &dbArticle{db: db}
}
