package service

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"errors"
	"github.com/gin-gonic/gin"
	"time"
)

type Service interface {
	model.WechatService
	model.CategoryService
	model.ArticleService
	model.RankService
	model.UserService
	model.TicketService
	model.CertificateService
}

type service struct {
	model.WechatService
	model.CategoryService
	model.ArticleService
	model.RankService
	model.UserService
	model.TicketService
	model.CertificateService
}

//ServiceError
var ServiceError = errors.New("service error")

func WechatCreate(ctx *gin.Context, wechat *model.Wechat) error {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.WechatCreate(wechat)
	}
	return ServiceError
}
func WechatList(ctx *gin.Context, limit, offset int) (wechats []*model.Wechat, count int64, err error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.WechatList(limit, offset)
	}
	return nil, 0, ServiceError
}

func WechatDelete(ctx *gin.Context, id int64) error {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.WechatDelete(id)
	}
	return ServiceError
}
func WechatListByCategory(ctx *gin.Context, cId int64, limit, offset int) (wechats []*model.Wechat, count int64, err error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.WechatListByCategory(cId, limit, offset)
	}
	return nil, 0, ServiceError
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

func TicketIsValid(ctx *gin.Context, ticketId string) (isValid bool, userId int64, err error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.TicketIsValid(ticketId)
	}
	return false, 0, ServiceError
}

func TicketGen(ctx *gin.Context, userId int64) (*model.Ticket, error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.TicketGen(userId)
	}
	return nil, ServiceError
}

func TicketDestroy(ctx *gin.Context, ticketId string) error {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.TicketDestroy(ticketId)
	}
	return ServiceError
}

func TicketTTL(ctx *gin.Context) time.Duration {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.TicketTTL()
	}
	return 0
}

func UserLogin(ctx *gin.Context, account, password string) (*model.Ticket, error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.UserLogin(account, password)
	}
	return nil, ServiceError
}

func UserRegister(ctx *gin.Context, account string, certificateType model.CertificateType, password string) (userId int64, err error) {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.UserRegister(account, certificateType, password)
	}
	return 0, ServiceError
}

func UserUpdatePassword(ctx *gin.Context, userId int64, newPassword string) error {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.UserUpdatePassword(userId, newPassword)
	}
	return nil
}

func CertificateUpdate(ctx *gin.Context, oldAccount, newAccount string, certificateType model.CertificateType) error {
	if service, ok := ctx.Value("service").(Service); ok {
		return service.CertificateUpdate(oldAccount, newAccount, certificateType)
	}
	return nil
}

func NewService(
	wSvc model.WechatService,
	cSvc model.CategoryService,
	aSvc model.ArticleService,
	rSvc model.RankService,
	tSvc model.TicketService,
	uSvc model.UserService,
	ccSvc model.CertificateService) Service {
	return &service{wSvc, cSvc, aSvc, rSvc, uSvc, tSvc, ccSvc}
}
