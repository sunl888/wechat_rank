package main

import (
	"code.aliyun.com/zmdev/wechat_rank/pkg/qingbo"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
)

const (
	AppKey = "ATSWhpmUd5c86zOZwGGx1fDM0ECoS0aL"
	AppId  = "1374"
)

func main() {
	ws := hashset.New("rmrbwx")
	client := qingbo.NewQingboClient(AppKey, AppId)
	account := qingbo.NewWxAccount(client)
	for i := 0; i < 1; i++ {
		wxname := ws.Values()[i].(string)
		_, err := account.GetAccount(wxname)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(wxname)
		articles, err := account.GetArticles(wxname, "", "", 0, 0)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(articles[0])
	}
}
