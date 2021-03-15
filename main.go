package main

import (
	"opq-plugin-tg2qq/client"
	"opq-plugin-tg2qq/conf"

	"opq-plugin-tg2qq/util/log"
	//"github.com/astaxie/beego/logs"
	//"go.uber.org/zap"
)

func main() {

	// 初始化全局配置
	conf.Init()

	// 初始化日志记录
	log.LogInit()

	go client.TGBotInit()

	client.OPQBotInit()
}
