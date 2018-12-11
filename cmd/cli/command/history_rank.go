package command

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/server"
	"errors"
	"fmt"
	"github.com/gpmgo/gopm/modules/log"
	"github.com/urfave/cli"
	"sort"
	"time"
)

func NewHistoryRankCommand(svr *server.Server) cli.Command {
	service := svr.Service
	return cli.Command{
		Name:  "history_rank",
		Usage: "获取所有公众号的历史文章列表",
		Action: func(c *cli.Context) error {
			var (
				wechats []*model.Wechat
				count   int64
				err     error
				now     time.Time
			)
			now = time.Now()
			// monday
			start, _ := time.Parse(DATE_FORMAT, "2018-10-01")
			end := start.AddDate(0, 0, 6)
			wechats, count, err = service.WechatList(0, 0)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("创建排名出错: %+v", err.Error()), 1)
			}
			// 星期一到星期日
			for {
				var ranks Data
				for i := 0; i < int(count); i++ {
					rankDetail, err := service.Rank(wechats[i], &model.Rank{
						Period:    "week",
						StartDate: start.Format(DATE_FORMAT),
						EndDate:   end.Format(DATE_FORMAT),
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
				start = end.AddDate(0, 0, 1)
				end = start.AddDate(0, 0, 6)
				if now.Sub(end).Hours() < 24 {
					break
				}
			}

			start, _ = time.Parse(DATE_FORMAT, "2018-10-01")
			end = start.AddDate(0, 1, -1)
			// 月初到月末
			for {
				var ranks Data
				for i := 0; i < int(count); i++ {
					rankDetail, err := service.Rank(wechats[i], &model.Rank{
						Period:    "month",
						StartDate: start.Format(DATE_FORMAT),
						EndDate:   end.Format(DATE_FORMAT),
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
				start = end.AddDate(0, 0, 1)
				end = start.AddDate(0, 1, -1)
				if now.Sub(end).Hours() < 24 {
					break
				}
			}
			return nil
		},
	}
}
