package service

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"errors"
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

//ServiceError
var ServiceError = errors.New("service error")

func WechatCreate(ctx *gin.Context, wechat *model.Wechat) error {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.WechatCreate(wechat)
	}
	return ServiceError
}
func WechatList(ctx *gin.Context) (wechats []*model.Wechat, err error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.WechatList()
	}
	return nil, ServiceError
}

func CategoryCreate(ctx *gin.Context, category *model.Category) error {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.CategoryCreate(category)
	}
	return ServiceError
}
func CategoryList(ctx *gin.Context) ([]*model.Category, error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.CategoryList()
	}
	return nil, ServiceError
}
func CategoryLoad(ctx *gin.Context, categoryId int64) (*model.Category, error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.CategoryLoad(categoryId)
	}
	return nil, ServiceError
}
func CategoryDelete(ctx *gin.Context, categoryId int64) error {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.CategoryDelete(categoryId)
	}
	return ServiceError
}
func CategoryUpdate(ctx *gin.Context, category *model.Category) error {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.CategoryUpdate(category)
	}
	return ServiceError
}

func NewService(wSvc model.WechatService, cSvc model.CategoryService, aSvc model.ArticleService) Service {
	return &service{wSvc, cSvc, aSvc}
}
