package service

import (
	"code.aliyun.com/zmdev/wechat_rank/errors"
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/pkg/qingbo"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"log"
	"time"
)

type articleService struct {
	model.ArticleStore
	model.WechatStore
	*qingbo.WxAccount
}

const perPage = 50         // 每页显示多少条
const MaxRequestCount = 50 // 每个公众号最多抓取50次

// 抓取文章
// 每个月最多发布 普通账号 8*31=248  超级账号 20*8*31=4960篇文章
func (aServ *articleService) ArticleGrab(wechat *model.Wechat, startDate, endDate string) error {
	// 防止重复获取文章
	articleIds := hashset.New()
	page := 1
	requestCount := 1
	sDate, _ := time.Parse("2006-01-02", startDate)
	eDate, _ := time.Parse("2006-01-02", endDate)
	//
	if sDate.Sub(eDate).Hours() < 24 {
		log.Println("开始日期和结束日期一致")
		//return errors.BadRequest("开始日期和结束日期一致", nil)
		return nil
	}
	for {
		if requestCount > MaxRequestCount {
			return errors.QingboError("超过最大请求次数",
				fmt.Sprintf("%s 公众号超过最大请求次数",
					wechat.WxName), 400, 400)
		}
		articles, err := aServ.WxAccount.GetArticles(wechat.WxName, sDate.String(), eDate.String(), perPage, page)
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
			publishedAt, _ := time.Parse("2006-01-02 15:04:05", article.CreatedAt)
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
		requestCount++
	}
	// 设置该公众号最近一次获取文章的时间
	wechat.LastGetArticleAt = eDate.String()
	err := aServ.WechatStore.WechatUpdate(wechat)
	if err != nil {
		return err
	}
	return nil
}

func NewArticleService(as model.ArticleStore, client *qingbo.WxAccount, wechat model.WechatStore) model.ArticleService {
	return &articleService{as, wechat, client}
}
