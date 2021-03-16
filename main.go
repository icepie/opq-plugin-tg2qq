package main

import (
	"fmt"
	"opq-plugin-tg2qq/client"
	"opq-plugin-tg2qq/conf"
)

func main() {
	fmt.Println(conf.ProConf)
	go client.TGBotInit()

	client.OPQBotInit()
}
