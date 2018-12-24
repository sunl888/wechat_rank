package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

const (
	LimitRequest = 3
	IntervalTime = time.Second * 3600 * 24
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
}
