package main

import (
	"opq-plugin-tg2qq/client"
	"opq-plugin-tg2qq/conf"

	"github.com/beego/beego/v2/core/logs"
	//"github.com/astaxie/beego/logs"
	//"go.uber.org/zap"
)

func main() {

	// 初始化全局配置
	conf.INIT()

//an official log.Logger
l := logs.GetLogger()
l.Println("this is a message of http")
//an official log.Logger with prefix ORM
logs.GetLogger("ORM").Println("this is a message of orm")

logs.Debug("my book is bought in the year of ", 2016)
logs.Info("this %s cat is %v years old", "yellow", 3)
logs.Warn("json is a type of kv like", map[string]int{"key": 2016})
logs.Error(1024, "is a very", "good game")
logs.Critical("oh,crash")


	f := &logs.PatternLogFormatter{
		Pattern:    "%F:%n|%w%t>> %m",
		WhenFormat: "2006-01-02",
	}
	logs.RegisterFormatter("pattern", f)

	_ = logs.SetLogger("console",`{"formatter": "pattern"}`)

	logs.Info("hello, world")


	go client.TGBotInit()

	client.OPQBotInit()
}
