package command

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/server"
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"sort"
	"time"
)

func NewManualRankCommand(svr *server.Server) cli.Command {
	service := svr.Service
	log := svr.Logger
	return cli.Command{
		Name:  "manual_rank",
		Usage: "手动计算所有公众号的排名",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "type",
				Value: "week",
				Usage: "按指定类型排名:week,month,year",
			},
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
				startDate string
				endDate   string
			)
			s, err := time.Parse(DATE_FORMAT, c.String("start"))
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("开始日期格式有误: %+v", err.Error()), 1)
			}
			e, err := time.Parse(DATE_FORMAT, c.String("end"))
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("结束日期格式有误: %+v", err.Error()), 1)
			}
			switch c.String("type") {
			case "week":
				if s.Weekday() != time.Monday {
					return cli.NewExitError("开始日期不是星期一", 1)
				}
				if e.Weekday() != time.Sunday {
					return cli.NewExitError("结束日期不是星期天", 1)
				}
				startDate = s.Format(DATE_FORMAT)
				endDate = e.Format(DATE_FORMAT)
			case "month":
				if s.Day() != 1 {
					return cli.NewExitError("开始日期不是一号", 1)
				}
				// todo  这里的结束日期需要验证
				if e.Day() != 31 {
					return cli.NewExitError("结束日期不是星期天", 1)
				}
				startDate = s.Format(DATE_FORMAT)
				endDate = e.Format(DATE_FORMAT)
			case "year":
				// TODO 省略
			default:
				return cli.NewExitError("类型错误", 2)
			}
			wechats, count, err := service.WechatList(0, 0)
			if err != nil {
				log.Error(fmt.Sprintf("创建排名出错: %+v", err.Error()))
				return cli.NewExitError(err, 3)
			}
			var ranks Data
			for i := 0; i < int(count); i++ {
				rankDetail, err := service.Rank(wechats[i], &model.Rank{
					Period:    c.String("type"),
					StartDate: startDate,
					EndDate:   endDate,
				})
				if err != nil {
					log.Error(fmt.Sprintf("创建排名出错: %+v", err.Error()))
					return cli.NewExitError(err, 4)
				}
				ranks = append(ranks, rankDetail)
			}
			sort.Sort(ranks)
			for i := 0; i < ranks.Len(); i++ {
				// 总排名
				ranks[i].TotalRank = i + 1
				err = service.RankDetailCreate(ranks[i])
				if err != nil {
					return cli.NewExitError(errors.New(fmt.Sprintf("创建排名出错: %+v", err.Error())), 4)
				}
			}
			return nil
		},
	}
}
