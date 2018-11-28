package command

import (
	"code.aliyun.com/zmdev/wechat_rank/server"
	"github.com/urfave/cli"
)

func RegisterCommand(svr *server.Server) []cli.Command {
	return []cli.Command{
		NewGetCommand(svr),
	}
}
