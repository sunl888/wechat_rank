package db_store

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"fmt"
	"github.com/jinzhu/gorm"
)

type dbWechat struct {
	db *gorm.DB
}

func (w *dbWechat) WechatSearch(keyword string, limit, offset int) (wechats []*model.WechatAndCategory, count int64, err error) {
	fmt.Println(keyword)
	wechats = make([]*model.WechatAndCategory, limit)
	q := w.db.Table("wechats w").
		Select("w.*,c.title as category_name").
		Joins("left join categories c on w.category_id=c.id").
		Where("wx_name like ? or wx_nickname like ?", "%"+keyword+"%", "%"+keyword+"%")
	q.Count(&count)
	err = q.Offset(offset).Limit(limit).Find(&wechats).Error
	return
}

func (w *dbWechat) WechatList(limit, offset int) (wechats []*model.Wechat, count int64, err error) {
	wechats = make([]*model.Wechat, 8)
	w.db.Model(model.Wechat{}).Count(&count)
	if limit == 0 && offset == 0 {
		err = w.db.Find(&wechats).Error
	} else {
		err = w.db.Offset(offset).Limit(limit).Find(&wechats).Error
	}
	return
}

func (w *dbWechat) WechatListByCategory(cId int64, limit, offset int) (wechats []*model.Wechat, count int64, err error) {
	wechats = make([]*model.Wechat, 10)
	w.db.Model(model.Wechat{}).Where("category_id = ?", cId).Count(&count)
	if limit == 0 && offset == 0 {
		err = w.db.Find(&wechats, "category_id = ?", cId).Error
	} else {
		err = w.db.Offset(offset).Limit(limit).Find(&wechats, "category_id = ?", cId).Error
	}
	return
}

func (w *dbWechat) WechatDelete(id int64) error {
	err := w.db.Delete(model.Article{}, "wx_id = ?", id).Error
	if err != nil {
		return err
	}
	err = w.db.Delete(model.RankDetail{}, "wx_id = ?", id).Error
	if err != nil {
		return err
	}
	err = w.db.Delete(model.Wechat{}, "id = ?", id).Error
	return err
}

func (w *dbWechat) WechatUpdate(wechat *model.Wechat) error {
	err := w.db.Model(model.Wechat{}).
		Where("wx_name = ?", wechat.WxName).
		Omit("wx_name").
		Update(&wechat).
		Error
	return err
}

func (w *dbWechat) WechatLoad(wechatName string) (wechat *model.Wechat, err error) {
	wechat = &model.Wechat{}
	err = w.db.First(&wechat, "wx_name = ?", wechatName).Error
	return
}

func (w *dbWechat) WechatCreate(wechat *model.Wechat) error {
	err := w.db.Create(&wechat).Error
	return err
}

func NewDBWechat(db *gorm.DB) model.WechatStore {
	return &dbWechat{db: db}
}
