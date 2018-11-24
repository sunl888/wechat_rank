package service

import "code.aliyun.com/zmdev/wechat_rank/model"

type wechatService struct {
	ws model.WechatStore
}

func NewWechatService(ws model.WechatStore) model.WecahtService {
	return &wechatService{ws: ws}
}
