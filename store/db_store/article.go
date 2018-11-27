package db_store

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"github.com/jinzhu/gorm"
)

type dbArticle struct {
	db *gorm.DB
}

// 保存文章列表
func (d *dbArticle) ArticleSave([]*model.Article) error {
	panic("implement me")
}

func NewDBArticle(db *gorm.DB) model.ArticleStore {
	return &dbArticle{db: db}
}
