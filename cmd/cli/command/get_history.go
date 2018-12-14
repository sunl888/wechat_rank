package command

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/server"
	"fmt"
	"github.com/urfave/cli"
	"time"
)

func NewGetHistoryCommand(svr *server.Server) cli.Command {
	service := svr.Service
	return cli.Command{
		Name:  "get_history",
		Usage: "获取所有公众号的历史文章列表",
		Action: func(c *cli.Context) error {
			var (
				wechats []*model.Wechat
				count   int64
				err     error
				now     time.Time
			)
			now = time.Now()
			start, _ := time.Parse(DATE_FORMAT, "2018-10-01")
			end := start.AddDate(0, 0, 5)
			wechats, count, err = service.WechatList(0, 0)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("创建排名出错: %+v", err.Error()), 1)
			}
			// 5天一次
			for {
				for i := 0; i < int(count); i++ {
					err = service.ArticleGrab(wechats[i], start.Format(DATE_FORMAT), end.Format(DATE_FORMAT))
					if err != nil {
						return cli.NewExitError(err, 2)
					}
					if (i+1)%10 == 0 {
						time.Sleep(1100 * time.Millisecond)
					}
				}
				start = end
				end = start.AddDate(0, 0, 5)
				if now.Sub(end).Hours() < 24 {
					break
				}
			}
			return nil
		},
	}
}
