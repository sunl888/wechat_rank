package command

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/server"
	"code.aliyun.com/zmdev/wechat_rank/utils"
	"fmt"
	"github.com/urfave/cli"
)

func NewGetCommand(svr *server.Server) cli.Command {
	service := svr.Service
	log := svr.Logger

	return cli.Command{
		Name:  "get",
		Usage: "获取所有公众号的最新文章列表",
		Action: func(c *cli.Context) {
			qingboClient := utils.NewQingboClient(svr.Conf.Qingbo.AppKey, svr.Conf.Qingbo.AppId)
			officialAccount := utils.NewOfficialAccount(qingboClient)
			wechats, err := service.WechatList()
			if err != nil {
				return
			}
			// 获取每个公众号最近的文章列表
			for _, w := range wechats {
				articlesResp, err := officialAccount.GetArticles(w.WxName, "", 50, 1)
				if err != nil {
					log.Error(fmt.Sprintf("获取文章失败: %+v\n", err.Error()))
					continue
				}
				// 保存文章
				articles := make([]*model.Article, len(articlesResp.DataResp))
				articles = convert2ArticlesModel(articlesResp)
				err = service.ArticleSave(articles)
				if err != nil {
					log.Error(fmt.Sprintf("保存文章失败: %+v\n", err.Error()))
					continue
				}
				// 保存最近一次获取文章的时间
				//err = service.WechatUpdate(&model.Wechat{
				//	WxName:       w.WxName,
				//	LastGrabTime: time.Now().Format("2006-01-02"),
				//})
				//if err != nil {
				//	log.Error(fmt.Sprintf("更新时间失败: %+v\n", err.Error()))
				//	continue
				//}
			}
		},
	}
}

func convert2ArticlesModel(articlesResp *utils.ArticleResponse) (articles []*model.Article) {
	for _, a := range articlesResp.DataResp {
		articles = append(articles, &model.Article{
			Url:         a.Url,
			Name:        a.Name,
			Title:       a.Title,
			Top:         a.Top,
			Author:      a.Author,
			Picurl:      a.Picurl,
			Digest:      a.Digest,
			WxName:      a.WxName,
			ArticleId:   a.Id,
			ReadCount:   a.ReadCount,
			LikeCount:   a.LikeCount,
			PublishedAt: a.CreatedAt,
			OriginalUrl: a.OriginalUrl,
		})
	}
	return
}
