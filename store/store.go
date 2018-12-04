package store

import "code.aliyun.com/zmdev/wechat_rank/model"

type Store interface {
	model.WechatStore
	model.CategoryStore
	model.ArticleStore
	model.RankStore
	model.UserStore
	model.TicketStore
	model.CertificateStore
}

type store struct {
	model.WechatStore
	model.CategoryStore
	model.ArticleStore
	model.RankStore
	model.UserStore
	model.TicketStore
	model.CertificateStore
}

func NewStore(
	ws model.WechatStore,
	cs model.CategoryStore,
	as model.ArticleStore,
	rs model.RankStore,
	us model.UserStore,
	ts model.TicketStore,
	ccs model.CertificateStore) Store {
	return &store{ws, cs, as, rs, us, ts, ccs}
}
