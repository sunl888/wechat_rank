package service

import "code.aliyun.com/zmdev/wechat_rank/model"

type categoryService struct {
	cs model.CategoryStore
}

func NewCategoryService(cs model.CategoryStore) model.CategoryService {
	return &categoryService{cs: cs}
}
