package main

import (
	"code.aliyun.com/zmdev/wechat_rank/cmd/cli/command"
	"code.aliyun.com/zmdev/wechat_rank/server"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	svr := server.SetupServer()
	app.Name = "微信排名命令行工具"
	app.Usage = ""
	app.Version = "0.0.1"
	app.Commands = append(app.Commands, command.RegisterCommand(svr)...)
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
