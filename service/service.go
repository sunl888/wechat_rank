package service

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"github.com/gin-gonic/gin"
)

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

func WechatCreate(ctx *gin.Context, wechat *model.Wechat) error {
	//TODO 这里好像只能这样写
	return ctx.Value("service").(Service).WechatCreate(wechat)
}
