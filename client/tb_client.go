package client

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"opq-plugin-tg2qq/conf"
	"opq-plugin-tg2qq/util/proxy"

	tb "gopkg.in/tucnak/telebot.v2"
)

func TGBotInit() {

	TGSet := tb.Settings{
		// You can also set custom API URL.
		// If field is empty it equals to "https://api.telegram.org".
		// URL: "http://195.129.111.17:8012",

		Token:  conf.ProConf.TGBot.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	}

	// setting proxy

	if conf.ProConf.TGBot.Proxy.Enable {

		purl, err := url.Parse(conf.ProConf.TGBot.Proxy.Url)
		if err != nil {
			log.Fatal(err)
		}

		if purl.Scheme == "http" {
			log.Println("[TG] 使用http代理")
			httpclient, err := proxy.HttpClient(conf.ProConf.TGBot.Proxy.Url)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			TGSet.Client = httpclient
		} else if purl.Scheme == "sock5" {
			sockclient, err := proxy.Socks5Client(purl.Host)
			log.Println("[TG] 使用sock5代理")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			TGSet.Client = sockclient
		}

	}

	TGBot, err := tb.NewBot(TGSet)
	if err != nil {
		log.Fatal(err)
		return
	} else {
		log.Println(TGBot.Me)
	}

	TGBot.Handle("/hello", func(m *tb.Message) {
		TGBot.Send(m.Sender, "Hello World!")
	})

	TGBot.Start()

}
