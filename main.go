package main

import (
	"code.aliyun.com/zmdev/wechat_rank/handler"
	"code.aliyun.com/zmdev/wechat_rank/server"
	"log"
	"net/http"
)

func main() {
	svr := server.SetupServer()
	log.Fatal(http.ListenAndServe(":8080", handler.CreateHTTPHandler(svr)))
}
