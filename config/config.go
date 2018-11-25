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
type Config struct {
	EnvVarPrefix string         `json:"env-var-prefix"`
	DB           DatabaseConfig `json:"database"`
	Redis        RedisConfig    `json:"redis"`
	Qingbo       QingboConfig   `json:"qingbo"`
}

func LoadConfig() Config {
	var c Config
	fileSource := file.NewSource(file.WithPath("config/config.yml"))
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
