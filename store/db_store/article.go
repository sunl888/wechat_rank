package db_store

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"github.com/jinzhu/gorm"
)

type dbArticle struct {
	db *gorm.DB
}

func (d *dbArticle) ArticleCreate(article *model.Article) error {
	var c int64
	q := d.db.Model(model.Article{}).Where("article_id = ?", article.ArticleId)
	q.Limit(1).Count(&c)
	//fmt.Printf("[database] wx_name: %s\tcount: %d\n", article.WxName, c)
	if c == 0 {
		err := d.db.Create(&article).Error
		if err != nil {
			return err
		}
	} else {
		err := q.Omit("article_id", "wx_name", "url", "published_at").Update(&article).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func NewDBArticle(db *gorm.DB) model.ArticleStore {
	return &dbArticle{db: db}
}
