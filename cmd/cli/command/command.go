package command

import (
	"code.aliyun.com/zmdev/wechat_rank/server"
	"github.com/urfave/cli"
)

func RegisterCommand(svr *server.Server) []cli.Command {
	return []cli.Command{
		NewGetCommand(svr),            // 获取文章
		NewRankCommand(svr),           // 自动排名
		NewManualGetCommand(svr),      // 手动获取文章
		NewGetHistoryCommand(svr),     // 历史文章获取(测试使用)
		NewManualRankCommand(svr),     // 手动排名
		NewHistoryRankCommand(svr),    // 历史文章获取后调用的排名(测试使用)
		NewCalcRankDetailCommand(svr), // 根据已存在的Rank 创建RankDetail(内部使用)
	}
}
