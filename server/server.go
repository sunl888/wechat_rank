package server

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"github.com/go-redis/redis"
	"code.aliyun.com/zmdev/wechat_rank/config"
)

type Server struct {
	Debug       bool
	RedisClient *redis.Client
	DB          *gorm.DB
	Conf        config.Config
	Logger      *zap.Logger
}
