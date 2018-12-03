package db_store

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"github.com/jinzhu/gorm"
	"time"
)

type dbArticle struct {
	db *gorm.DB
}

func (d *dbArticle) ArticleRank(startDate, endDate string, categoryId int64, offset, limit int) (articles []*model.Article, count int64, err error) {
	s, err := time.Parse(DATE_FORMAT, startDate)
	if err != nil {
		return nil, 0, err
	}
	e, err := time.Parse(DATE_FORMAT, endDate)
	if err != nil {
		return nil, 0, err
	}
	e = e.Add(time.Duration(time.Second*60*60*24 - 1))
	q := d.db.Table("articles a").
		Select("a.*").
		Joins("left join wechats w on a.wx_id = w.id").
		Where("a.published_at >=? and a.published_at <=?", s, e)
	if categoryId != 0 {
		q = q.Where("w.category_id = ?", categoryId)
	}
	q.Count(&count)
	articles = make([]*model.Article, 0, limit)
	err = q.Order("a.read_count desc,a.like_count desc").Offset(offset).Limit(limit).Find(&articles).Error
	return
}

const DATE_FORMAT = "2006-01-02"

func (d *dbArticle) ArticleList(startDate, endDate string, offset, limit int) (articles []*model.Article, count int64, err error) {

	s, err := time.Parse(DATE_FORMAT, startDate)
	if err != nil {
		return nil, 0, err
	}
	e, err := time.Parse(DATE_FORMAT, endDate)
	if err != nil {
		return nil, 0, err
	}
	e = e.Add(time.Duration(time.Second*60*60*24 - 1))

	q := d.db.Model(model.Article{}).
		Where("published_at >=? and published_at <=?", s, e)
	q.Count(&count)
	articles = make([]*model.Article, 0, limit)
	err = q.Offset(offset).Limit(limit).Find(&articles).Error
	return
}

func (d *dbArticle) ArticleCreate(article *model.Article) error {
	var c int64
	q := d.db.Model(model.Article{}).Where("article_id = ?", article.ArticleId)
	q.Limit(1).Count(&c)
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
