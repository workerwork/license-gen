package conf

import (
	"bytes"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
)

// logger 日志配置结构
type logger struct {
	Level  string `mapstructure:"level"`
	Pretty bool   `mapstructure:"pretty"`
	Color  bool   `mapstructure:"color"`
}

// LoggerConf 日志配置
var LoggerConf = &logger{}

// license 配置结构
type license struct {
	Exec string `mapstructure:"exec"`
	Src  string `mapstructure:"src"`
	Bin  string `mapstructure:"bin"`
}

// License 配置
var LicenseConf = &license{}

// .env变量
var URL string

// Setup 生成服务配置
func Setup() {
	// 读取.env配置
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	URL = viper.GetString("URL")
    log.Println(".env parsed success!")

	// 读取配置文件内容
	viper.SetConfigType("YAML")
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Read 'config.yml' fail: %v\n", err)
	}
	// 配置内容解析
	viper.ReadConfig(bytes.NewBuffer(data))
	// 解析配置赋值
	viper.UnmarshalKey("logger", LoggerConf)
	viper.UnmarshalKey("license", LicenseConf)
    log.Println("config.yml parsed success!")
}
