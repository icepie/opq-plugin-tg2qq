package main

import (
	"opq-plugin-tg2qq/client"
	"opq-plugin-tg2qq/conf"
)

func main() {

	// 初始化全局配置
	conf.INIT()

	go client.TGBotInit()

	client.OPQBotInit()
}
