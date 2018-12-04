package command

import (
	"code.aliyun.com/zmdev/wechat_rank/server"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"time"
)

const DATE_FORMAT = "2006-01-02"

func NewGetCommand(svr *server.Server) cli.Command {
	service := svr.Service
	log := svr.Logger

	return cli.Command{
		Name:  "get",
		Usage: "获取所有公众号的最新文章列表",
		Action: func(c *cli.Context) {
			now := time.Now()
			if now.Weekday() != time.Monday {
				log.Error("日期有误.")
				return
			}
			laskWeekStartDate := now.AddDate(0, 0, -7).Format(DATE_FORMAT)
			laskWeekEndDate := now.AddDate(0, 0, -1).Format(DATE_FORMAT)
			// TODO date
			//laskWeekStartDate := "2018-11-25"
			//laskWeekEndDate := "2018-12-01"
			err := service.ArticleGrab(laskWeekStartDate, laskWeekEndDate)
			if err != nil {
				log.Error("文章抓取失败.", zap.String("detail", err.Error()))
				return
			}
		},
	}
}
