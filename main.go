package main

import (
	"opq-plugin-tg2qq/client"
)

func main() {

	go client.TGBotInit()

	client.OPQBotInit()
}
