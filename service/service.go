package service

import "code.aliyun.com/zmdev/wechat_rank/model"

type Service interface {
	model.WecahtService
}

type service struct {
	model.WecahtService
}

func NewService(wSvc model.WecahtService) Service {
	return &service{wSvc}
}
