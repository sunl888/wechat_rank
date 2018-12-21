package service

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
)

type categoryService struct {
	model.CategoryStore
}

func (cs *categoryService) CategoryCreate(category *model.Category) error {
	err := cs.CategoryStore.CategoryCreate(category)
	if err != nil {
		return err
	}
	return err
}

func (cs *categoryService) CategoryUpdate(category *model.Category) error {
	_, loadErr := cs.CategoryStore.CategoryLoad(category.Id)
	if loadErr != nil {
		return loadErr
	}
	err := cs.CategoryStore.CategoryUpdate(category)
	if err != nil {
		return err
	}
	return nil
}

func (cs *categoryService) CategoryDelete(cId int64) (err error) {
	_, loadErr := cs.CategoryStore.CategoryLoad(cId)
	if loadErr != nil {
		return loadErr
	}
	err = cs.CategoryStore.CategoryDelete(cId)
	if err != nil {
		return err
	}
	return nil
}

func (cs *categoryService) CategoryList(isLogin bool) (c []*model.Category, err error) {
	onlyShowPrivate := isLogin
	c, err = cs.CategoryStore.CategoryList(onlyShowPrivate)
	if err != nil {
		return nil, err
	}
	return
}
func NewCategoryService(cs model.CategoryStore) model.CategoryService {
	return &categoryService{cs}
}
