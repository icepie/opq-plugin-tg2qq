package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"

	"opq-plugin-tg2qq/client/opqbot"
	"opq-plugin-tg2qq/conf"
	"opq-plugin-tg2qq/util/proxy"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	TGBot *tb.Bot
	MG    = &MyGroup{ChatID: conf.ProConf.TGBot.ChatID}
)

type MyGroup struct {
	ChatID string
}

// Recipient returns personal group chatID
func (mg MyGroup) Recipient() string {
	return mg.ChatID
}

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

	var err error
	TGBot, err = tb.NewBot(TGSet)
	if err != nil {
		logs.Emergency("[TGBot] Connet %s", err.Error())
	} else {
		logs.Info("[TGBot] Online: %+v", *TGBot.Me)
		TGBot.Send(MG, fmt.Sprintf("[opq-plugin-tg2qq] starting...\n\nQQ Group Num: %d", conf.ProConf.OPQBot.Group))
	}

	// Group text msg handler
	TGBot.Handle(tb.OnText, func(m *tb.Message) {
		if strconv.Itoa(int(m.Chat.ID)) == conf.ProConf.TGBot.ChatID {
			if m.Sender.ID != TGBot.Me.ID {
				OPQBot.Send(opqbot.SendMsgPack{
					SendType:   opqbot.SendTypeTextMsg,
					SendToType: opqbot.SendToTypeGroup,
					ToUserUid:  conf.ProConf.OPQBot.Group,
					Content:    opqbot.SendTypeTextMsgContent{Content: fmt.Sprintf("[TG] [%s] %s", m.Sender.Username, m.Text)},
				})
			}
			logs.Info("-> [TGbot] %+v", m.Chat)
		}
	})

	TGBot.Handle(tb.OnPhoto, func(m *tb.Message) {
		if strconv.Itoa(int(m.Chat.ID)) == conf.ProConf.TGBot.ChatID {
			if m.Sender.ID != TGBot.Me.ID {

				fileURL, err := TGBot.FileURLByID(m.Photo.FileID)
				if err != nil {
					logs.Error(err)
					return
				}

				// cacheDir := util.GetAppPath() + "/.cache"
				// exist, err := util.PathExists(cacheDir)
				// if err != nil {
				// 	logs.Error(err)
				// 	return
				// }

				// if !exist {
				// 	err := os.Mkdir(cacheDir, os.ModePerm)
				// 	if err != nil {
				// 		logs.Error(err)
				// 	}
				// }

				var resp *http.Response
				if conf.ProConf.TGBot.Proxy.Enable {
					// Get the data use proxy cilent
					resp, err = TGSet.Client.Get(fileURL)
					if err != nil {
						logs.Error(err)
					}
				} else {
					resp, err = http.Get(fileURL)
					if err != nil {
						logs.Error(err)
					}
				}

				defer resp.Body.Close()

				// outFilePath := cacheDir + "/" + m.Photo.FileID
				// // 创建一个文件用于保存
				// out, err := os.Create(outFilePath)
				// if err != nil {
				// 	logs.Error(err)
				// }
				// defer out.Close()

				// // 然后将响应流和文件流对接起来
				// _, err = io.Copy(out, resp.Body)
				// if err != nil {
				// 	panic(err)
				// }

				// 读取获取的[]byte数据
				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					logs.Error(err)
					return
				}

				imageBase64 := base64.StdEncoding.EncodeToString(data)
				//fmt.Println("base64", imageBase64)

				// 使用旧版opq api

				opqbody := opqbot.SendPicMsgPackV1{
					ToUser:      conf.ProConf.OPQBot.Group,
					SendMsgType: "PicMsg",
					SendToType:  opqbot.SendToTypeGroup,
					Content:     fmt.Sprintf("[TG] [%s] %s", m.Sender.Username, m.Caption),
					// GroupID:     0,
					// AtUser:      0,
					PicBase64Buf: imageBase64,
				}

				b, _ := json.Marshal(opqbody)

				opqresp, err := http.Post(fmt.Sprintf("%s%s?qq=%d&funcname=SendMsg&timeout=10", conf.ProConf.OPQBot.Url, "/v1/LuaApiCaller", conf.ProConf.OPQBot.QQ),
					"application/json",
					bytes.NewBuffer(b))
				if err != nil {
					fmt.Println(err)
				}
				defer opqresp.Body.Close()
				// body, _ := ioutil.ReadAll(opqresp.Body)

			}
			logs.Info("-> [TGbot] %+v", m)
		}
	})

	// Group voice msg handler
	TGBot.Handle(tb.OnVoice, func(m *tb.Message) {
		if strconv.Itoa(int(m.Chat.ID)) == conf.ProConf.TGBot.ChatID {
			if m.Sender.ID != TGBot.Me.ID {
				OPQBot.Send(opqbot.SendMsgPack{
					SendType:   opqbot.SendTypeTextMsg,
					SendToType: opqbot.SendToTypeGroup,
					ToUserUid:  conf.ProConf.OPQBot.Group,
					// Content:    opqbot.SendTypeTextMsgContent{Content: fmt.Sprintf("[TG] [%s] {VoiceUrl: %s}", m.Sender.Username, fileURL)},
					Content: opqbot.SendTypeTextMsgContent{Content: fmt.Sprintf("[TG] [%s] [🗣️]", m.Sender.Username)},
				})
			}
			logs.Info("-> [TGbot] %+v", m.Voice)
		}
	})

	TGBot.Handle(tb.OnAudio, func(m *tb.Message) {
		if strconv.Itoa(int(m.Chat.ID)) == conf.ProConf.TGBot.ChatID {
			logs.Info("-> [TGbot] %+v", m.Chat)

			if m.Sender.ID != TGBot.Me.ID {

				fileURL, err := TGBot.FileURLByID(m.Audio.FileID)
				if err != nil {
					logs.Error(err)
					return
				}

				OPQBot.Send(opqbot.SendMsgPack{
					SendType:   opqbot.SendTypeTextMsg,
					SendToType: opqbot.SendToTypeGroup,
					ToUserUid:  conf.ProConf.OPQBot.Group,
					// Content:    opqbot.SendTypeTextMsgContent{Content: fmt.Sprintf("[TG] [%s] {VoiceUrl: %s}", m.Sender.Username, fileURL)},
					Content: opqbot.SendTypeTextMsgContent{Content: fmt.Sprintf("[TG] [%s] [🎵|%s]\n%s", m.Sender.Username, m.Audio.FileName, m.Caption)},
				})

				OPQBot.Send(opqbot.SendMsgPack{
					SendType:   opqbot.SendTypeVoiceByUrl,
					SendToType: opqbot.SendToTypeGroup,
					ToUserUid:  conf.ProConf.OPQBot.Group,
					Content:    opqbot.SendTypeVoiceByUrlContent{VoiceUrl: fileURL},
				})
			}
			logs.Info("-> [TGbot] %+v", m.Chat)
		}
	})

	TGBot.Start()
}
