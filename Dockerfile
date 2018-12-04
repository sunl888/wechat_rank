FROM golang:1.11.1-alpine3.7 as builder
COPY . /go/src/code.aliyun.com/zmdev/wechat_rank

RUN go build -v -o /go/src/code.aliyun.com/zmdev/wechat_rank/server_app /go/src/code.aliyun.com/zmdev/wechat_rank/cmd/server/main.go
RUN go build -v -o /go/src/code.aliyun.com/zmdev/wechat_rank/cli_app /go/src/code.aliyun.com/zmdev/wechat_rank/cmd/cli/cli.go

FROM alpine:3.7
RUN apk update && apk --no-cache add tzdata
ENV TZ=Asia/Shanghai
ENV HOME=/app
COPY --from=builder /go/src/code.aliyun.com/zmdev/wechat_rank/server_app /app/server
COPY --from=builder /go/src/code.aliyun.com/zmdev/wechat_rank/cli_app /app/cli
COPY --from=builder /go/src/code.aliyun.com/zmdev/wechat_rank/config/config.yml /app/config/config.yml
COPY --from=builder /go/src/code.aliyun.com/zmdev/wechat_rank/crontabs /var/spool/cron/crontabs/root
COPY --from=builder /go/src/code.aliyun.com/zmdev/wechat_rank/.env.docker.example /app/.env

EXPOSE 8080

WORKDIR /app
RUN chmod +x /app/server /app/cli
CMD crond && ./server
