package main

import (
	"code.aliyun.com/zmdev/wechat_rank/pkg/qingbo"
	"fmt"
)

const (
	AppKey = "ATSWhpmUd5c86zOZwGGx1fDM0ECoS0aL"
	AppId  = "1374"
)

func main() {
	client := qingbo.NewQingboClient(AppKey, AppId)
	account := qingbo.NewWxAccount(client)
	group := qingbo.NewWxGroup(client)

	wxname := "rmrbwx"
	weixin, err := account.GetAccount(wxname)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(weixin)
	articles, err := account.GetArticles(wxname, "", "", 0, 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(articles[0])
	wx, err := group.AddWx2Group("107562", "105622")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(wx.Data)
}
