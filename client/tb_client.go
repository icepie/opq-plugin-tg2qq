package client

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/astaxie/beego/logs"

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

	logs.Info("[TGBot] Threads starting...")

	// setting proxy
	if conf.ProConf.TGBot.Proxy.Enable {

		purl, err := url.Parse(conf.ProConf.TGBot.Proxy.Url)
		if err != nil {
			logs.Error("[TGBot] Proxy:", "Can not parse the proxy url")
		}

		if purl.Scheme == "http" {
			logs.Info("[TGBot] Proxy: http")
			httpclient, err := proxy.HttpClient(conf.ProConf.TGBot.Proxy.Url)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			TGSet.Client = httpclient
		} else if purl.Scheme == "sock5" {
			sockclient, err := proxy.Socks5Client(purl.Host)
			logs.Info("[TGBot] Proxy: %s", "sock5")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			TGSet.Client = sockclient
		}

	}

	TGBot, err := tb.NewBot(TGSet)
	if err != nil {
		logs.Emergency("[TGBot] Connet %s", err.Error())
	} else {
		logs.Info("[TGBot] Online: %+v", *TGBot.Me)
	}

	TGBot.Handle("/hello", func(m *tb.Message) {
		TGBot.Send(m.Sender, "Hello World!")
	})

	TGBot.Start()

}
