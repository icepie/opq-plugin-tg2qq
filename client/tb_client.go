package client

import (
	"fmt"
	"log"
	"os"
	"time"

	"opq-plugin-tg2qq/util/proxy"

	tb "gopkg.in/tucnak/telebot.v2"
)

func TGBotInit() {
	Socks5client, err := proxy.HttpClient("http://127.0.0.1:12333")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	TGBot, err := tb.NewBot(tb.Settings{
		// You can also set custom API URL.
		// If field is empty it equals to "https://api.telegram.org".
		// URL: "http://195.129.111.17:8012",

		Token:  "",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		Client: Socks5client,
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	TGBot.Handle("/hello", func(m *tb.Message) {
		TGBot.Send(m.Sender, "Hello World!")
	})

	TGBot.Start()

}
