package store

import "code.aliyun.com/zmdev/wechat_rank/model"

type Store interface {
	model.WechatStore
}

type store struct {
	model.WechatStore
}

func NewStore(ws model.WechatStore) Store {
	return &store{ws}
}
