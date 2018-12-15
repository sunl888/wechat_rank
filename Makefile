docker:
	docker build -t wechat_rank:latest .
docker_upload: docker
	docker tag wechat_rank:latest registry.cn-hangzhou.aliyuncs.com/wqer1019/wechat_rank:latest
	docker push registry.cn-hangzhou.aliyuncs.com/wqer1019/wechat_rank:latest
run:
	docker run -d --network app_net -v $(pwd)/../wechat_rank-frontend/build/:/app/frontend -p 80:8080  wechat_rank:latest

pull:
	docker pull registry.cn-hangzhou.aliyuncs.com/wqer1019/wechat_rank:latest