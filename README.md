## 微信排名
对公众号进行多维度计算排名

## docker 
  * docker network create app_net (第一次)
  * docker run -d --network app_net -v $(pwd)/../wechat_rank-frontend/build/:/app/frontend -p 80:8080  wechat_rank:latest
  * docker run -it --rm --network app_net -v $(pwd)/../wechat_rank-frontend/build/:/app/frontend -p 80:8080  wechat_rank:latest /bin/sh
## 启动
go run cmd/server/main.go
## cli
go run cmd/cli/cli.go
 - get 
 - rank --type=week,month,year