package server

import (
	"code.aliyun.com/zmdev/wechat_rank/config"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type Server struct {
	Debug       bool
	RedisClient *redis.Client
	DB          *gorm.DB
	Conf        config.Config
	Logger      *zap.Logger
}
