package service

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"errors"
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

var ServiceError = errors.New("service error")

func NewService(wSvc model.WechatService, cSvc model.CategoryService, aSvc model.ArticleService, rSvc model.RankService, tSvc model.TicketService, uSvc model.UserService, ccSvc model.CertificateService) Service {
	return &service{wSvc, cSvc, aSvc, rSvc, uSvc, tSvc, ccSvc}
}
