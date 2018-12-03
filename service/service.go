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
	model.RankService
}

type service struct {
	model.WechatService
	model.CategoryService
	model.ArticleService
	model.RankService
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

func RankList(ctx *gin.Context, period string) (ranks []*model.Rank, err error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.RankList(period)
	}
	return nil, ServiceError
}

func RankLoad(ctx *gin.Context, rankId int64) (rank *model.Rank, err error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.RankLoad(rankId)
	}
	return nil, ServiceError
}

func RankDetail(ctx *gin.Context, rankId, categoryId int64, limit, offset int) (ranks []*model.RankJoinWechat, count int64, err error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.RankDetail(rankId, categoryId, limit, offset)
	}
	return nil, 0, ServiceError
}

func ArticleList(ctx *gin.Context, startDate, endDate string, limit, offset int) (articles []*model.Article, err error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.ArticleList(startDate, endDate, offset, limit)
	}
	return nil, ServiceError
}
func ArticleRank(ctx *gin.Context, startDate, endDate string, categoryId int64, offset, limit int) (articles []*model.Article, count int64, err error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.ArticleRank(startDate, endDate, categoryId, offset, limit)
	}
	return nil, 0, ServiceError
}
func NewService(wSvc model.WechatService, cSvc model.CategoryService, aSvc model.ArticleService, rSvc model.RankService) Service {
	return &service{wSvc, cSvc, aSvc, rSvc}
}
