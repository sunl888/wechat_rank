package store

import "code.aliyun.com/zmdev/wechat_rank/model"

type Store interface {
	model.WechatStore
	model.CategoryStore
	model.ArticleStore
	model.RankStore
}

type store struct {
	model.WechatStore
	model.CategoryStore
	model.ArticleStore
	model.RankStore
}

func NewStore(ws model.WechatStore, cs model.CategoryStore, as model.ArticleStore, rs model.RankStore) Store {
	return &store{ws, cs, as, rs}
}
