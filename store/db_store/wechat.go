package db_store

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"github.com/jinzhu/gorm"
)

type dbWechat struct {
	db *gorm.DB
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
	// todo 删除关联表的相关记录
	err := w.db.Delete(model.Wechat{}, "id = ?", id).Error
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
	return w.db.Create(&wechat).Error
}

func NewDBWechat(db *gorm.DB) model.WechatStore {
	return &dbWechat{db: db}
}
