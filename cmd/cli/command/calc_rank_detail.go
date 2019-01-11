package command

import (
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/server"
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"sort"
)

func NewCalcRankDetailCommand(svr *server.Server) cli.Command {
	service := svr.Service
	return cli.Command{
		Name:  "calc_rank_detail",
		Usage: "手动计算所有公众号的排名详情, 在排名表已存在的情况下",
		Action: func(c *cli.Context) error {
			var (
				wechats []*model.Wechat
				count   int64
				err     error
			)
			wechats, count, err = service.WechatList(0, 0)
			if err != nil {
				return cli.NewExitError(fmt.Sprintf("创建排名出错: %+v", err.Error()), 1)
			}
			var ranksDetail Data
			ranks, err := service.RankAllList()
			for _, rank := range ranks {
				for i := 0; i < int(count); i++ {
					rd, err := service.Rank(wechats[i], rank)
					if err != nil {
						return cli.NewExitError(err, 2)
					}
					ranksDetail = append(ranksDetail, rd)
				}
				sort.Sort(ranksDetail)
				for i := 0; i < ranksDetail.Len(); i++ {
					// 总排名
					ranksDetail[i].TotalRank = i + 1
					err = service.RankDetailCreate(ranksDetail[i])
					if err != nil {
						return cli.NewExitError(errors.New(fmt.Sprintf("创建排名出错: %+v", err.Error())), 4)
					}
				}
			}

			return nil
		},
	}
}
