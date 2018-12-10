package config

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/env"
	"github.com/micro/go-config/source/file"
	"log"
)

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string `json:"db-name"`
}

type RedisConfig struct {
	Address string
	Port    string
}

type QingboConfig struct {
	AppId  string
	AppKey string
}
type TicketConfig struct {
	Driver string `json:"driver"` // ticket 使用的驱动 只支持 redis 和 database
	TTL    int64  `json:"ttl"`    // ticket 的过期时间 （秒）
}
type Config struct {
	EnvVarPrefix string         `json:"env-var-prefix"`
	ServiceName  string         `json:"service-name"`
	AppSalt      string         `json:"app_salt"`
	DB           DatabaseConfig `json:"database"`
	Redis        RedisConfig    `json:"redis"`
	Qingbo       QingboConfig   `json:"qingbo"`
	Ticket       TicketConfig   `json:"ticket"`
}

func LoadConfig(path string) Config {
	var c Config
	fileSource := file.NewSource(file.WithPath(path))
	checkErr(config.Load(fileSource))
	// env 的配置会覆盖文件中的配置
	envSource := env.NewSource(env.WithStrippedPrefix(config.Get("env-var-prefix").String("RANK")))
	checkErr(config.Load(envSource))
	checkErr(config.Scan(&c))
	return c
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
