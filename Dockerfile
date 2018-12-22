FROM golang:1.11.1-alpine3.7 as builder

ENV SRC=/go/src/code.aliyun.com/zmdev/wechat_rank
COPY . $SRC

RUN go build -v -o $SRC/app $SRC/cmd/server/main.go && \
    go build -v -o $SRC/cli $SRC/cmd/cli/cli.go

FROM alpine:3.7

ENV TZ=Asia/Shanghai
ENV HOME=/app

COPY --from=builder /go/src/code.aliyun.com/zmdev/wechat_rank/app /app/server
COPY --from=builder /go/src/code.aliyun.com/zmdev/wechat_rank/cli /app/cli
COPY --from=builder /go/src/code.aliyun.com/zmdev/wechat_rank/config/config.yml /app/config/config.yml
COPY --from=builder /go/src/code.aliyun.com/zmdev/wechat_rank/crontabs /var/spool/cron/crontabs/root
COPY --from=builder /go/src/code.aliyun.com/zmdev/wechat_rank/.env.docker.example /app/.env

RUN apk update && \
    apk --no-cache add tzdata && \
    chmod +x $HOME/server $HOME/cli

WORKDIR $HOME

EXPOSE 8080

CMD crond && ./server