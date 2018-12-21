package db_store

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"github.com/jinzhu/gorm"
)

type dbCategory struct {
	db *gorm.DB
}

func (c *dbCategory) CategoryUpdate(category *model.Category) error {
	return c.db.Omit("created_at").Save(category).Error
}

func (c *dbCategory) CategoryCreate(category *model.Category) error {
	return c.db.Create(&category).Error
}

func (c *dbCategory) CategoryLoad(cId int64) (category *model.Category, err error) {
	category = &model.Category{}
	err = c.db.Where("id = (?)", cId).First(&category).Error
	return
}

func (c *dbCategory) CategoryList(onlyShowPrivate bool) (categories []*model.Category, err error) {
	categories = make([]*model.Category, 0, 5)
	q := c.db.Model(model.Category{})
	if onlyShowPrivate == false {
		q = q.Where("is_private = ?", false)
	}
	err = q.Find(&categories).Error
	return
}

func (c *dbCategory) CategoryDelete(cId int64) (err error) {
	return c.db.Delete(model.Category{}, "id = ?", cId).Error
}

func NewDBCategory(db *gorm.DB) model.CategoryStore {
	return &dbCategory{db}
}
