package main

import (
	"fmt"
	"license-gen/conf"
	"license-gen/logger"
	"license-gen/server"
)

// go:generate go env -w GO111MODULE=on
// go:generate go env -w GOPROXY=https://goproxy.cn,direct
// go:generate go mod tidy
// go:generate go mod download

// @title License Generate Program
// @version 0.1.0
// @description This is a License Generate program for WCG/5GC/EPC
// @.env parse the API from LicenseCenter
// @config.yml parse the config of the program
func main() {
	conf.Setup()   //配置初始化
	banner()       //打印banner
	logger.Setup() //日志初始化
	server.Serve() //启动服务
}

func banner() {
	fmt.Printf(`
  **************************************************************************************
  **                              License Generate Program                            **
  **************************************************************************************
  欢迎使用 License Generate Program来自动构建生产License文件
  当前版本:V0.1.0 beta
  本程序是LicenseCenter的后端，配合LicenseCenter自动生产WCG/5GC/EPC License文件
  服务端提供的API:
  %s
  %s
  %s  
  **************************************************************************************
`, conf.URL_GET, conf.URL_POST1, conf.URL_POST2)
}
