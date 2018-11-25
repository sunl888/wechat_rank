package main

import (
	"code.aliyun.com/zmdev/wechat_rank/handler"
	"code.aliyun.com/zmdev/wechat_rank/server"
	"log"
	"net/http"
)

func main() {
	//qingboWeixinClient := utils.NewQingboClient("hE8COKMsiobwWPlXdqgnING2hmbchdoA", "1345", "weixin")
	//resp, err :=qingboWeixinClient.Get("users","wx_name=hnnu1958")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(resp)

	svr := server.SetupServer()
	log.Fatal(http.ListenAndServe(":8080", handler.CreateHTTPHandler(svr)))
}
