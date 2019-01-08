package command

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/server"
	"fmt"
	"github.com/urfave/cli"
	"time"
)

const (
	MaxDays  = 90 // 最多间隔 90 天
	MaxHours = MaxDays * 24
)

func NewManualGetCommand(svr *server.Server) cli.Command {
	service := svr.Service
	return cli.Command{
		Name:  "manual_get",
		Usage: "手动获取所有公众号的文章列表",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "start",
				Value: "2019-01-01",
				Usage: "开始日期",
			},
			cli.StringFlag{
				Name:  "end",
				Value: "2019-01-07",
				Usage: "结束日期",
			},
		},
		Action: func(c *cli.Context) error {
			var (
				wechats   []*model.Wechat
				count     int64
				err       error
				startDate string
				endDate   string
			)

			s, err := time.Parse(DATE_FORMAT, c.String("start"))
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("开始时间格式有误: %+v", err.Error()), 1)
			}
			e, err := time.Parse(DATE_FORMAT, c.String("end"))
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("结束时间格式有误: %+v", err.Error()), 1)
			}
			if s.Sub(e).Hours() > MaxHours {
				return cli.NewExitError(fmt.Sprintf("开始日期与结束日期超过%d天.", MaxDays), 1)
			}
			startDate = s.Format(DATE_FORMAT)
			endDate = e.Format(DATE_FORMAT)
			wechats, count, err = service.WechatList(0, 0)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("创建排名出错: %+v", err.Error()), 1)
			}
			for i := 0; i < int(count); i++ {
				err = service.ArticleGrab(wechats[i], startDate, endDate)
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
