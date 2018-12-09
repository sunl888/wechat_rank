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

type Data []*model.RankDetail

func (p Data) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Data) Len() int           { return len(p) }
func (p Data) Less(i, j int) bool { return p[i].Wci > p[j].Wci }

func NewRankCommand(svr *server.Server) cli.Command {
	service := svr.Service
	log := svr.Logger
	return cli.Command{
		Name:  "rank",
		Usage: "计算所有公众号的排名",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "type",
				Value: "week",
				Usage: "按指定类型排名:week,month,year",
			},
		},
		Action: func(c *cli.Context) error {
			var (
				startDate string
				endDate   string
			)
			switch c.String("type") {
			case "week":
				// 上周一到上周日
				// TODO date
				//now := time.Now()
				//if now.Weekday() == time.Monday {
				//	startDate = now.AddDate(0, 0, -7).Format(DATE_FORMAT)
				//	endDate = now.AddDate(0, 0, -1).Format(DATE_FORMAT)
				//} else {
				//	return cli.NewExitError(fmt.Sprintf("日期不正确,type:%s,startDate:%s,endDate:%s\n", c.String("type"), startDate, endDate), 1)
				//}
				// 2018-11-26 2018-12-02
				startDate = "2018-12-03"
				endDate = "2018-12-09"
			case "month":
				year, month, _ := time.Now().Date()
				thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
				startDate = thisMonth.AddDate(0, -1, 0).Format(DATE_FORMAT)
				endDate = thisMonth.AddDate(0, 0, -1).Format(DATE_FORMAT)
			case "year":
				year, _, _ := time.Now().Date()
				t := time.Date(year-1, 1, 1, 0, 0, 0, 0, time.Local)
				startDate = t.Format(DATE_FORMAT)
				endDate = t.AddDate(0, 12, -1).Format(DATE_FORMAT)
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
