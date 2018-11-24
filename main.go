package main

import (
	"code.aliyun.com/zmdev/wechat_rank/server"
	"code.aliyun.com/zmdev/wechat_rank/handler"
	"net/http"
	"log"
)

func main() {
	svr := server.SetupServer()
	log.Fatal(http.ListenAndServe(":8080", handler.CreateHTTPHandler(svr)))
}
