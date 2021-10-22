package main

import (
	"license-gen/conf"
	"license-gen/logger"
	"license-gen/server"
)

func main() {

	//基本配置初始化
	conf.Setup()

	//日志初始化
	logger.Setup()

	//启动http服务
	server.Serve()
}
