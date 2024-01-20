package main

import (
	"ABots/common"
	"ABots/core"
)

func main() {

	// 初始化日志工具
	common.InitLogger()
	// 读取配置
	conf := common.LoadDefaultConf()
	// 主线程
	abot := core.InitAbot(conf)
	abot.Run()
}
