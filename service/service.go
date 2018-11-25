package service

import "code.aliyun.com/zmdev/wechat_rank/model"

type Service interface {
	model.WechatService
	model.CategoryService
	model.ArticleService
}

type service struct {
	model.WechatService
	model.CategoryService
	model.ArticleService
}

func NewService(wSvc model.WechatService, cSvc model.CategoryService, aSvc model.ArticleService) Service {
	return &service{wSvc, cSvc, aSvc}
}
