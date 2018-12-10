package command

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/server"
	"fmt"
	"github.com/urfave/cli"
	"time"
)

const DATE_FORMAT = "2006-01-02"

func NewGetCommand(svr *server.Server) cli.Command {
	service := svr.Service
	return cli.Command{
		Name:  "get",
		Usage: "获取所有公众号的最新文章列表",
		Action: func(c *cli.Context) error {
			var (
				wechats   []*model.Wechat
				count     int64
				err       error
				now       time.Time
				startDate string
				yesterday string
			)
			now = time.Now()
			yesterday = now.AddDate(0, 0, -1).Format(DATE_FORMAT)
			//yesterday = "2018-11-01"
			wechats, count, err = service.WechatList(0, 0)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("创建排名出错: %+v", err.Error()), 1)
			}
			for i := 0; i < int(count); i++ {
				if wechats[i].LastGetArticleAt == "" {
					// seven days ago
					startDate = now.AddDate(0, 0, -7).Format(DATE_FORMAT)
					//startDate = "2018-10-01"
				} else {
					startDate = wechats[i].LastGetArticleAt
				}
				err = service.ArticleGrab(wechats[i], startDate, yesterday)
				if err != nil {
					return cli.NewExitError(err, 2)
				}
				if (i+1)%10 == 0 {
					time.Sleep(1100 * time.Millisecond)
				}
			}
			return nil
		},
	}
}
