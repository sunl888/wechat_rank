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

func WechatCreate(ctx *gin.Context, wechat *model.Wechat) error {
	return ctx.Value("service").(Service).WechatCreate(wechat)
}

func CategoryCreate(ctx *gin.Context, category *model.Category) error {
	return ctx.Value("service").(Service).CategoryCreate(category)
}
func CategoryList(ctx *gin.Context) ([]*model.Category, error) {
	return ctx.Value("service").(Service).CategoryList()
}
func CategoryLoad(ctx *gin.Context, categoryId int64) (*model.Category, error) {
	return ctx.Value("service").(Service).CategoryLoad(categoryId)
}
func CategoryDelete(ctx *gin.Context, categoryId int64) error {
	return ctx.Value("service").(Service).CategoryDelete(categoryId)
}
func CategoryUpdate(ctx *gin.Context, category *model.Category) error {
	return ctx.Value("service").(Service).CategoryUpdate(category)
}

func NewService(wSvc model.WechatService, cSvc model.CategoryService, aSvc model.ArticleService) Service {
	return &service{wSvc, cSvc, aSvc}
}
