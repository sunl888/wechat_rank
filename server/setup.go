package server

import (
	"code.aliyun.com/zmdev/wechat_rank/config"
	"code.aliyun.com/zmdev/wechat_rank/model"
	"code.aliyun.com/zmdev/wechat_rank/pkg/hasher"
	"code.aliyun.com/zmdev/wechat_rank/service"
	"code.aliyun.com/zmdev/wechat_rank/store"
	"code.aliyun.com/zmdev/wechat_rank/store/db_store"
	"code.aliyun.com/zmdev/wechat_rank/store/redis_store"
	"code.aliyun.com/zmdev/wechat_rank/utils"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"path"
	"runtime"

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
		&model.Rank{},
		&model.RankDetail{},
		&model.User{},
		&model.Certificate{},
	)
}

func setupRedis(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func SetupServer() *Server {
	s := &Server{}
	s.Debug = os.Getenv("DEBUG") == "true"
	var err error
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前路径失败. ERR:%s", err.Error())
	}
	s.Conf = config.LoadConfig(path.Join(pwd, "config/config.yml"))
	s.ServiceName = s.Conf.ServiceName
	s.RedisClient = setupRedis(s.Conf.Redis.Address + ":" + s.Conf.Redis.Port)
	s.DB = setupGorm(
		/*s.Debug*/ true,
		s.Conf.DB.Driver,
		s.Conf.DB.Host,
		s.Conf.DB.Port,
		s.Conf.DB.DBName,
		s.Conf.DB.User,
		s.Conf.DB.Password,
	)
	if s.Debug {
		s.Logger, err = zap.NewDevelopment()
	} else {
		s.Logger, err = zap.NewProduction()
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
		db_store.NewDBRank(s.DB),
		db_store.NewDBUser(s.DB),
		redis_store.NewRedisTicket(s.RedisClient),
		db_store.NewDBCertificate(s.DB),
	)
}

func setupService(serv *Server) service.Service {
	qingboClient := utils.NewQingboClient(serv.Conf.Qingbo.AppKey, serv.Conf.Qingbo.AppId)
	officialAccount := utils.NewOfficialAccount(qingboClient)
	s := setupStore(serv)
	h := hasher.NewArgon2Hasher(
		[]byte(serv.Conf.AppSalt),
		3,
		32<<10,
		uint8(runtime.NumCPU()),
		32,
	)
	tSvc := service.NewTicketService(s, time.Duration(serv.Conf.Ticket.TTL)*time.Microsecond)
	return service.NewService(
		service.NewWechatService(s, officialAccount),
		service.NewCategoryService(s),
		service.NewArticleService(s, officialAccount, s),
		service.NewRankService(s, s, s),
		tSvc,
		service.NewUserService(s, s, tSvc, h),
		service.NewCertificateService(s),
	)
}
