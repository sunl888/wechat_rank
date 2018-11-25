package store

import "code.aliyun.com/zmdev/wechat_rank/model"

type Store interface {
	model.WechatStore
	model.CategoryStore
	model.ArticleStore
}

type store struct {
	model.WechatStore
	model.CategoryStore
	model.ArticleStore
}

func NewStore(ws model.WechatStore, cs model.CategoryStore, as model.ArticleStore) Store {
	return &store{ws, cs, as}
}
