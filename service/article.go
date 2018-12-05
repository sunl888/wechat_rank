package service

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/utils"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/pkg/errors"
	"time"
)

type articleService struct {
	model.ArticleStore
	model.WechatStore
	*utils.OfficialAccount
}

const DateFormat = "2006-01-02 15:04:05"

// 抓取文章
func (aServ *articleService) ArticleGrab(laskWeekStartDate, laskWeekEndDate string) error {
	wechats, _, err := aServ.WechatStore.WechatList(0, 0)
	if err != nil {
		return err
	}
	// 获取每个公众号最近的文章列表
	for i, w := range wechats {
		ids := hashset.New()
		index := 1
	A:
		for {
			articleResp, err := aServ.OfficialAccount.GetArticles(w.WxName, laskWeekStartDate, laskWeekEndDate, 50, index)
			if err != nil {
				return errors.New(fmt.Sprintf("获取文章失败: %+v\n", err.Error()))
			}
			// 保存文章
			for _, a := range articleResp {
				if ids.Contains(a.Id) {
					// 有些公众号获取到的都是重复文章  不得不这样写...
					// Tip: 清博大数据Api是史上最垃圾的Api
					break A
				} else {
					ids.Add(a.Id)
				}
				var publishedAt time.Time
				publishedAt, _ = time.Parse(DateFormat, a.CreatedAt)
				err := aServ.ArticleStore.ArticleCreate(&model.Article{
					WxId:         w.Id,
					Top:          a.Top,
					Url:          a.Url,
					Title:        a.Title,
					WxName:       a.WxName,
					ArticleId:    a.Id,
					ReadCount:    a.ReadCount,
					LikeCount:    a.LikeCount,
					PublishedAt:  &publishedAt,
					WxVerifyName: w.VerifyName,
					WxCategoryId: w.CategoryId,
				})
				if err != nil {
					return errors.New(fmt.Sprintf("保存文章失败: %+v\n", err.Error()))
				}
			}
			if len(articleResp) < 50 {
				break
			}
			index += 1
			if (index-1)%10 == 0 {
				time.Sleep(1300 * time.Millisecond)
			}
		}
		if (i+1)%10 == 0 {
			time.Sleep(1300 * time.Millisecond)
		}
	}
	return nil
}

func NewArticleService(as model.ArticleStore, client *utils.OfficialAccount, wechat model.WechatStore) model.ArticleService {
	return &articleService{as, wechat, client}
}
