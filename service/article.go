package service

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/utils"
	"fmt"
	"github.com/pkg/errors"
	"time"
)

type articleService struct {
	model.ArticleStore
	model.WechatStore
	*utils.OfficialAccount
}

// 抓取文章
// laskWeekStartDate := "2018-11-19"
// laskWeekEndDate := "2018-11-25"
func (aServ *articleService) ArticleGrab(laskWeekStartDate, laskWeekEndDate string) error {
	wechats, err := aServ.WechatStore.WechatList()
	if err != nil {
		return err
	}
	// 获取每个公众号最近的文章列表
	for i, w := range wechats {
		index := 1
		for {
			articleResp, err := aServ.OfficialAccount.GetArticles(w.WxName, laskWeekStartDate, laskWeekEndDate, 50, index)
			if err != nil {
				return errors.New(fmt.Sprintf("获取文章失败: %+v\n", err.Error()))
			}
			// 保存文章
			for _, a := range articleResp {
				var publishedAt time.Time
				publishedAt, _ = time.Parse("2006-01-02 15:04:05", a.CreatedAt)
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
		}
		if (i+1)%10 == 0 {
			// 延时1.2秒
			time.Sleep(1200 * time.Millisecond)
			fmt.Printf("第%d次延时\n", (i+1)/10)
		}
	}
	return nil
}

func NewArticleService(as model.ArticleStore, client *utils.OfficialAccount, wechat model.WechatStore) model.ArticleService {
	return &articleService{as, wechat, client}
}
