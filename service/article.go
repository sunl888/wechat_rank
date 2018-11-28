package service

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
)

type articleService struct {
	model.ArticleStore
}

func NewArticleService(as model.ArticleStore) model.ArticleService {
	return &articleService{as}
}
