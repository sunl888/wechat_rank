## 微信排名
对公众号进行多维度计算排名

go run cmd/server/main.go

## docker 
  * ~~docker run -d --network app_net -p 8080:8080 --env-file ".env" wechat_rank:latest~~
  * docker run -d --network app_net -p 8080:8080  wechat_rank:latest