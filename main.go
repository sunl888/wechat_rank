package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

const (
	LimitRequest = 3
	IntervalTime = time.Second * 3600 * 24
	DateFormat   = "2006-01-02"
)

func main() {
	r := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	fmt.Println(r.Exists("t"))
	r.Set("t", 2, time.Duration(time.Second*1000))
	fmt.Println(r.Exists("t"))
	fmt.Println(r.Get("t"))

	fmt.Println(fmt.Sprintf("拒绝更新,%f小时内最多只能请求%d次该接口", IntervalTime.Hours(), LimitRequest))

	s, _ := time.Parse("2006-01-02", "2019-01-01")
	e, _ := time.Parse("2006-01-02", "2019-01-31")
	year, month, _ := e.Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	endDate := thisMonth.AddDate(0, 0, -1).Format("2006-01-02")
	fmt.Println(endDate, s)
}
