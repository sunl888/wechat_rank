package db_store

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"github.com/jinzhu/gorm"
)

type dbUser struct {
	db *gorm.DB
}

func (u *dbUser) UserExist(id int64) (bool, error) {
	var count uint8
	err := u.db.Model(&model.User{}).Where(model.User{Id: id}).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (u *dbUser) UserIsNotExistErr(err error) bool {
	return model.UserIsNotExistErr(err)
}

func (u *dbUser) UserLoad(id int64) (user *model.User, err error) {
	if id <= 0 {
		return nil, model.ErrUserNotExist
	}
	user = &model.User{}
	err = u.db.Where(model.User{Id: id}).First(user).Error
	if gorm.IsRecordNotFoundError(err) {
		err = model.ErrUserNotExist
	}
	return
}

func (u *dbUser) UserUpdate(user *model.User) error {
	if user.Id <= 0 {
		return model.ErrUserNotExist
	}
	return u.db.Omit("created_at").Save(user).Error
}

func (u *dbUser) UserCreate(user *model.User) error {
	return u.db.Create(user).Error
}

func NewDBUser(db *gorm.DB) model.UserStore {
	return &dbUser{db: db}
}
