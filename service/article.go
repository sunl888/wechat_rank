package service

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/utils"
	"github.com/emirpasic/gods/sets/hashset"
	"time"
)

type articleService struct {
	model.ArticleStore
	model.WechatStore
	*utils.OfficialAccount
}

const DateFormat = "2006-01-02 15:04:05"
const perPage = 50

// 抓取文章
func (aServ *articleService) ArticleGrab(wechat *model.Wechat, laskWeekStartDate, laskWeekEndDate string) error {
	// 获取每个公众号最近的文章列表
	articleIds := hashset.New()
	page := 1
	for {
		articles, err := aServ.OfficialAccount.GetArticles(wechat.WxName, laskWeekStartDate, laskWeekEndDate, perPage, page)
		if err != nil {
			return err
		}
		// 保存文章
		for _, article := range articles {
			// 有些公众号获取到的都是重复文章  不得不这样写...
			// Tip: 清博大数据Api是史上最垃圾的Api
			if articleIds.Contains(article.Id) {
				return nil
			} else {
				articleIds.Add(article.Id)
			}
			publishedAt, _ := time.Parse(DateFormat, article.CreatedAt)
			err := aServ.ArticleStore.ArticleCreate(&model.Article{
				WxId:        wechat.Id,
				Top:         article.Top,
				Url:         article.Url,
				Desc:        article.Digest,
				Title:       article.Title,
				ArticleId:   article.Id,
				ReadCount:   article.ReadCount,
				LikeCount:   article.LikeCount,
				PublishedAt: &publishedAt,
			})
			if err != nil {
				return err
			}
		}
		if len(articles) < perPage {
			break
		}
		if page%10 == 0 {
			time.Sleep(1100 * time.Millisecond)
		}
		page++
	}
	return nil
}

func NewArticleService(as model.ArticleStore, client *utils.OfficialAccount, wechat model.WechatStore) model.ArticleService {
	return &articleService{as, wechat, client}
}
