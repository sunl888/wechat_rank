package service

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/utils"
)

type articleService struct {
	model.ArticleStore
	client *utils.OfficialAccount
}

func (a *articleService) ArticleSave() error {
	articles := make([]*model.Article, 10)
	return a.ArticleStore.ArticleSave(articles)
}

func NewArticleService(as model.ArticleStore, client *utils.OfficialAccount) model.ArticleService {
	return &articleService{as, client}
}
