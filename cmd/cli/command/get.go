package command

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/server"
	"code.aliyun.com/zmdev/wechat_rank/utils"
	"fmt"
	"github.com/urfave/cli"
	"time"
)

func NewGetCommand(svr *server.Server) cli.Command {
	service := svr.Service
	return cli.Command{
		Name:  "get",
		Usage: "获取所有公众号的最新文章列表",
		Action: func(c *cli.Context) {
			qingboClient := utils.NewQingboClient(svr.Conf.Qingbo.AppKey, svr.Conf.Qingbo.AppId)
			officialClient := utils.NewOfficialAccount(qingboClient)
			wechats, err := service.WechatList()
			if err != nil {
				return
			}
			ch := make(chan int)
			// 获取每个公众号最近的文章列表
			for i, w := range wechats {
				go run(svr, officialClient, w, ch, i)
				if (i+1)%10 == 0 {
					// 延时1.5秒
					time.Sleep(1500 * time.Millisecond)
					fmt.Printf("第%d次延时\n", (i+1)/10)
				}
			}
			// 等待所有的goruntine执行完成后退出程序
			for {
				if <-ch == len(wechats)-1 {
					break;
				}
			}
		},
	}
}

func run(svr *server.Server, c *utils.OfficialAccount, w *model.Wechat, count chan int, index int) {
	service := svr.Service
	log := svr.Logger

	articlesResp, err := c.GetArticles(w.WxName, "", 50, 1)
	if err != nil {
		log.Error(fmt.Sprintf("获取文章失败: %+v\n", err.Error()))
		return
	}
	// 保存文章
	for _, a := range articlesResp.DataResp {
		err := service.ArticleCreate(&model.Article{
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
		if err != nil {
			log.Error(fmt.Sprintf("保存文章失败: %+v\n", err.Error()))
			return
		}
		time.Sleep(time.Millisecond * 100)
	}

	count <- index
}
