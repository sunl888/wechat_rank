package server

import (
	"code.aliyun.com/zmdev/wechat_rank/config"
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/service"
	"code.aliyun.com/zmdev/wechat_rank/store"
	"code.aliyun.com/zmdev/wechat_rank/store/db_store"
	"fmt"
	"github.com/jinzhu/gorm"
	// 引入数据库驱动注册及初始化
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"log"
	"os"
	"time"
)

func setupGorm(debug bool, driverName, dbHost, dbPort, dbName, dbUser, dbPassword string) *gorm.DB {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPassword,
		dbHost+":"+dbPort,
		dbName,
	)
	var (
		db  *gorm.DB
		err error
	)
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(driverName, dataSourceName)
		if err == nil {
			db.LogMode(debug)
			autoMigrate(db)
			return db
		}
		fmt.Println(driverName, dataSourceName)
		log.Println(err)
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("数据库链接失败！error: %+v", err)
	return nil
}

func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.Wechat{},
	)
}

func SetupServer() *Server {
	s := &Server{}
	s.Debug = os.Getenv("DEBUG") == "true"
	s.Conf = config.LoadConfig()
	// s.RedisClient = setupRedis(s.Conf.Redis.Address + ":" + s.Conf.Redis.Port)
	s.DB = setupGorm(
		s.Debug,
		s.Conf.DB.Driver,
		s.Conf.DB.Host,
		s.Conf.DB.Port,
		s.Conf.DB.DBName,
		s.Conf.DB.User,
		s.Conf.DB.Password,
	)
	var err error
	if s.Debug {
		s.Logger, err = zap.NewProduction()
	} else {
		s.Logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func setupStore(s *Server) store.Store {
	return store.NewStore(
		db_store.NewDBWechat(s.DB),
	)
}

func SetupService(serv *Server) service.Service {
	s := setupStore(serv)
	return service.NewService(
		service.NewWechatService(s),
	)
}
