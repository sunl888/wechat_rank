package server

import (
	"code.aliyun.com/zmdev/wechat_rank/config"
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/service"
	"code.aliyun.com/zmdev/wechat_rank/store"
	"code.aliyun.com/zmdev/wechat_rank/store/db_store"
	"code.aliyun.com/zmdev/wechat_rank/utils"
	"fmt"
	"github.com/jinzhu/gorm"
	"path"

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
			if debug {
				autoMigrate(db)
			}
			return db
		}
		log.Println(err)
		time.Sleep(2 * time.Second)
	}
	log.Fatalf("数据库链接失败！ error: %+v", err)
	return nil
}

func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.Wechat{},
		&model.Category{},
		&model.Article{},
	)
}

func SetupServer() *Server {
	s := &Server{}
	s.Debug = os.Getenv("DEBUG") == "true"
	var err error
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前路径失败. ERR:%s", err.Error())
	}
	s.Conf = config.LoadConfig(path.Join(pwd, "../../config/config.yml"))
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
	if s.Debug {
		s.Logger, err = zap.NewProduction()
	} else {
		s.Logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatal(err)
	}
	s.Service = setupService(s)
	return s
}

func setupStore(s *Server) store.Store {
	return store.NewStore(
		db_store.NewDBWechat(s.DB),
		db_store.NewDBCategory(s.DB),
		db_store.NewDBArticle(s.DB),
	)
}

func setupService(serv *Server) service.Service {
	qingboClient := utils.NewQingboClient(serv.Conf.Qingbo.AppKey, serv.Conf.Qingbo.AppId)

	officialAccount := utils.NewOfficialAccount(qingboClient)
	s := setupStore(serv)
	return service.NewService(
		service.NewWechatService(s, officialAccount),
		service.NewCategoryService(s),
		service.NewArticleService(s),
	)
}
