package main

import (
	"code.aliyun.com/zmdev/wechat_rank/handler"
	"code.aliyun.com/zmdev/wechat_rank/server"
	"log"
	"net/http"
)

func main() {
	//qingboWeixinClient := utils.NewQingboClient("hE8COKMsiobwWPlXdqgnING2hmbchdoA", "1345")
	//official := utils.NewOfficialAccount(qingboWeixinClient)
	//resp, err := official.GetArticles("wwwtongcheng", "", 50, 0)
	//if err != nil {
	//	panic(err)
	//}
	//for _, r := range resp.DataResp {
	//	fmt.Println(r.Top, r.WxName, r.Name, r.Title, r.CreatedAt)
	//}

	svr := server.SetupServer()
	log.Fatal(http.ListenAndServe(":8080", handler.CreateHTTPHandler(svr)))
}
