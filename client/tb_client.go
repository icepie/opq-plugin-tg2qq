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
	"github.com/wxnacy/wgo/arrays"

	"opq-plugin-tg2qq/client/opqbot"
	"opq-plugin-tg2qq/conf"
	"opq-plugin-tg2qq/util"
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
		//TGBot.Send(MG, fmt.Sprintf("[opq-plugin-tg2qq] starting...\n\nQQ Group Num: %d", conf.ProConf.OPQBot.Group))
	}

	// Group text msg handler
	TGBot.Handle(tb.OnText, func(m *tb.Message) {
		if strconv.Itoa(int(m.Chat.ID)) == conf.ProConf.TGBot.ChatID {
			if m.Sender.ID != TGBot.Me.ID && arrays.ContainsString(conf.ProConf.TGBot.FilterID, m.Sender.Recipient()) == -1 {

				username := m.Sender.Username

				if username == "" {
					if m.Sender.LastName != "" {
						username = fmt.Sprintf("%s %s", m.Sender.FirstName, m.Sender.LastName)
					} else {
						username = fmt.Sprintf("%s", m.Sender.FirstName)
					}
				}

				if m.IsForwarded() {
					OPQBot.Send(opqbot.SendMsgPack{
						SendType:   opqbot.SendTypeTextMsg,
						SendToType: opqbot.SendToTypeGroup,
						ToUserUid:  conf.ProConf.OPQBot.Group,
						Content:    opqbot.SendTypeTextMsgContent{Content: fmt.Sprintf("[TG] %s - Forwarded : %s", username, m.Text)},
					})
					logs.Info("-> [TGbot] %+v", m)
				} else if m.IsReply() {

					replyToUsername := m.ReplyTo.Sender.Username

					if replyToUsername == "" {
						if m.Sender.LastName != "" {
							replyToUsername = fmt.Sprintf("%s %s", m.ReplyTo.Sender.FirstName, m.ReplyTo.Sender.LastName)
						} else {
							replyToUsername = fmt.Sprintf("%s", m.ReplyTo.Sender.FirstName)
						}
					}

					OPQBot.Send(opqbot.SendMsgPack{
						SendType:   opqbot.SendTypeTextMsg,
						SendToType: opqbot.SendToTypeGroup,
						ToUserUid:  conf.ProConf.OPQBot.Group,
						Content:    opqbot.SendTypeTextMsgContent{Content: fmt.Sprintf("[TG] %s -> %s : %s", username, replyToUsername, m.Text)},
					})
					logs.Info("-> [TGbot] %+v", m)
				} else {
					OPQBot.Send(opqbot.SendMsgPack{
						SendType:   opqbot.SendTypeTextMsg,
						SendToType: opqbot.SendToTypeGroup,
						ToUserUid:  conf.ProConf.OPQBot.Group,
						Content:    opqbot.SendTypeTextMsgContent{Content: fmt.Sprintf("[TG] %s : %s", username, m.Text)},
					})
				}

			}
			logs.Info("-> [TGbot] %+v", m.Chat)
		}
	})

	TGBot.Handle(tb.OnPhoto, func(m *tb.Message) {
		if strconv.Itoa(int(m.Chat.ID)) == conf.ProConf.TGBot.ChatID {
			if m.Sender.ID != TGBot.Me.ID && arrays.ContainsString(conf.ProConf.TGBot.FilterID, m.Sender.Recipient()) == -1 {

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

				username := m.Sender.Username

				if username == "" {
					if m.Sender.LastName != "" {
						username = fmt.Sprintf("%s %s", m.Sender.FirstName, m.Sender.LastName)
					} else {
						username = fmt.Sprintf("%s", m.Sender.FirstName)
					}
				}

				// 使用旧版opq api

				opqbody := opqbot.SendPicMsgPackV1{
					ToUser:      conf.ProConf.OPQBot.Group,
					SendMsgType: "PicMsg",
					SendToType:  opqbot.SendToTypeGroup,
					// GroupID:     0,
					// AtUser:      0,
					PicBase64Buf: imageBase64,
				}

				content := fmt.Sprintf("[TG] %s", username)

				if m.IsReply() {
					replyToUsername := m.ReplyTo.Sender.Username

					if replyToUsername == "" {
						if m.Sender.LastName != "" {
							replyToUsername = fmt.Sprintf("%s %s", m.ReplyTo.Sender.FirstName, m.ReplyTo.Sender.LastName)
						} else {
							replyToUsername = fmt.Sprintf("%s", m.ReplyTo.Sender.FirstName)
						}
					}

					content = fmt.Sprintf("[TG] %s -> %s", username, replyToUsername)
				}

				if m.IsForwarded() {
					content = fmt.Sprintf("%s - Forwarded", content)
				}

				if m.Caption != "" {
					content = fmt.Sprintf("%s : %s", content, m.Caption)
				}

				opqbody.Content = content

				b, _ := json.Marshal(opqbody)

				opqresp, err := http.Post(fmt.Sprintf("%s%s?qq=%d&funcname=SendMsg&timeout=10", conf.ProConf.OPQBot.Url, "/v1/LuaApiCaller", conf.ProConf.OPQBot.QQ),
					"application/json",
					bytes.NewBuffer(b))
				if err != nil {
					fmt.Println(err)
				}
				defer opqresp.Body.Close()
				// body, _ := ioutil.ReadAll(opqresp.Body)

				// v2 base64 发送图片
				// OPQBot.Send(opqbot.SendMsgPack{
				// 	SendType:   opqbot.SendTypePicMsgByBase64,
				// 	SendToType: opqbot.SendToTypeGroup,
				// 	ToUserUid:  conf.ProConf.OPQBot.QQ,
				// 	Content:    opqbot.SendTypePicMsgByBase64Content{Content: content, Base64: imageBase64},
				// })

			}
			logs.Info("-> [TGbot] %+v", m)
		}
	})

	// Group voice msg handler
	TGBot.Handle(tb.OnVoice, func(m *tb.Message) {
		if strconv.Itoa(int(m.Chat.ID)) == conf.ProConf.TGBot.ChatID {
			if m.Sender.ID != TGBot.Me.ID && arrays.ContainsString(conf.ProConf.TGBot.FilterID, m.Sender.Recipient()) == -1 {

				username := m.Sender.Username

				if username == "" {
					if m.Sender.LastName != "" {
						username = fmt.Sprintf("%s %s", m.Sender.FirstName, m.Sender.LastName)
					} else {
						username = fmt.Sprintf("%s", m.Sender.FirstName)
					}
				}

				content := fmt.Sprintf("[TG] %s - Voice", username)

				if m.IsReply() {
					replyToUsername := m.ReplyTo.Sender.Username

					if replyToUsername == "" {
						if m.Sender.LastName != "" {
							replyToUsername = fmt.Sprintf("%s %s", m.ReplyTo.Sender.FirstName, m.ReplyTo.Sender.LastName)
						} else {
							replyToUsername = fmt.Sprintf("%s", m.ReplyTo.Sender.FirstName)
						}
					}

					content = fmt.Sprintf("[TG] %s -> %s - Voice", username, replyToUsername)
				}

				if m.IsForwarded() {
					content = fmt.Sprintf("%s - Forwarded", content)
				}

				if m.Caption != "" {
					content = fmt.Sprintf("%s :\n%s", content, m.Caption)
				}

				mp := opqbot.SendMsgPack{
					SendType:   opqbot.SendTypeTextMsg,
					SendToType: opqbot.SendToTypeGroup,
					ToUserUid:  conf.ProConf.OPQBot.Group,
					Content:    opqbot.SendTypeTextMsgContent{Content: content},
				}

				OPQBot.Send(mp)
			}
			logs.Info("-> [TGbot] %+v", m.Voice)
		}
	})

	TGBot.Handle(tb.OnAudio, func(m *tb.Message) {
		if strconv.Itoa(int(m.Chat.ID)) == conf.ProConf.TGBot.ChatID {
			if m.Sender.ID != TGBot.Me.ID && arrays.ContainsString(conf.ProConf.TGBot.FilterID, m.Sender.Recipient()) == -1 {

				// fileURL, err := TGBot.FileURLByID(m.Audio.FileID)
				// if err != nil {
				// 	logs.Error(err)
				// 	return
				// }

				username := m.Sender.Username

				if username == "" {
					if m.Sender.LastName != "" {
						username = fmt.Sprintf("%s %s", m.Sender.FirstName, m.Sender.LastName)
					} else {
						username = fmt.Sprintf("%s", m.Sender.FirstName)
					}
				}

				content := fmt.Sprintf("[TG] %s - Audio - %s", username, m.Audio.FileName)

				if m.IsReply() {
					replyToUsername := m.ReplyTo.Sender.Username

					if replyToUsername == "" {
						if m.Sender.LastName != "" {
							replyToUsername = fmt.Sprintf("%s %s", m.ReplyTo.Sender.FirstName, m.ReplyTo.Sender.LastName)
						} else {
							replyToUsername = fmt.Sprintf("%s", m.ReplyTo.Sender.FirstName)
						}
					}

					content = fmt.Sprintf("[TG] %s -> %s - Audio - %s", username, replyToUsername, m.Audio.FileName)
				}

				if m.IsForwarded() {
					content = fmt.Sprintf("%s - Forwarded", content)
				}

				if m.Caption != "" {
					content = fmt.Sprintf("%s :\n%s", content, m.Caption)
				}

				mp := opqbot.SendMsgPack{
					SendType:   opqbot.SendTypeTextMsg,
					SendToType: opqbot.SendToTypeGroup,
					ToUserUid:  conf.ProConf.OPQBot.Group,
					Content:    opqbot.SendTypeTextMsgContent{Content: content},
				}

				OPQBot.Send(mp)

				// OPQBot.Send(opqbot.SendMsgPack{
				// 	SendType:   opqbot.SendTypeVoiceByUrl,
				// 	SendToType: opqbot.SendToTypeGroup,
				// 	ToUserUid:  conf.ProConf.OPQBot.Group,
				// 	Content:    opqbot.SendTypeVoiceByUrlContent{VoiceUrl: fileURL},
				// })
			}
			logs.Info("-> [TGbot] %+v", m.Chat)
		}
	})

	TGBot.Handle(tb.OnVideo, func(m *tb.Message) {
		if strconv.Itoa(int(m.Chat.ID)) == conf.ProConf.TGBot.ChatID {
			if m.Sender.ID != TGBot.Me.ID && arrays.ContainsString(conf.ProConf.TGBot.FilterID, m.Sender.Recipient()) == -1 {

				username := m.Sender.Username

				if username == "" {
					if m.Sender.LastName != "" {
						username = fmt.Sprintf("%s %s", m.Sender.FirstName, m.Sender.LastName)
					} else {
						username = fmt.Sprintf("%s", m.Sender.FirstName)
					}
				}

				content := fmt.Sprintf("[TG] %s - Video - %s", username, m.Video.FileName)

				if m.IsReply() {
					replyToUsername := m.ReplyTo.Sender.Username

					if replyToUsername == "" {
						if m.Sender.LastName != "" {
							replyToUsername = fmt.Sprintf("%s %s", m.ReplyTo.Sender.FirstName, m.ReplyTo.Sender.LastName)
						} else {
							replyToUsername = fmt.Sprintf("%s", m.ReplyTo.Sender.FirstName)
						}
					}

					content = fmt.Sprintf("[TG] %s -> %s - Video - %s", username, replyToUsername, m.Video.FileName)
				}

				if m.IsForwarded() {
					content = fmt.Sprintf("%s - Forwarded", content)
				}

				if m.Caption != "" {
					content = fmt.Sprintf("%s :\n%s", content, m.Caption)
				}

				mp := opqbot.SendMsgPack{
					SendType:   opqbot.SendTypeTextMsg,
					SendToType: opqbot.SendToTypeGroup,
					ToUserUid:  conf.ProConf.OPQBot.Group,
					Content:    opqbot.SendTypeTextMsgContent{Content: content},
				}

				OPQBot.Send(mp)

			}
			logs.Info("-> [TGbot] %+v", m.Chat)
		}
	})

	TGBot.Handle(tb.OnDocument, func(m *tb.Message) {
		if strconv.Itoa(int(m.Chat.ID)) == conf.ProConf.TGBot.ChatID {
			if m.Sender.ID != TGBot.Me.ID && arrays.ContainsString(conf.ProConf.TGBot.FilterID, m.Sender.Recipient()) == -1 {

				username := m.Sender.Username

				if username == "" {
					if m.Sender.LastName != "" {
						username = fmt.Sprintf("%s %s", m.Sender.FirstName, m.Sender.LastName)
					} else {
						username = fmt.Sprintf("%s", m.Sender.FirstName)
					}
				}

				content := fmt.Sprintf("[TG] %s - Document - %s - %s", username, m.Document.FileName, util.ByteSize(uint64(m.Document.FileSize)))

				if m.IsForwarded() {
					content = fmt.Sprintf("%s - Forwarded", content)
				}

				if m.IsReply() {
					replyToUsername := m.ReplyTo.Sender.Username

					if replyToUsername == "" {
						if m.Sender.LastName != "" {
							replyToUsername = fmt.Sprintf("%s %s", m.ReplyTo.Sender.FirstName, m.ReplyTo.Sender.LastName)
						} else {
							replyToUsername = fmt.Sprintf("%s", m.ReplyTo.Sender.FirstName)
						}
					}

					content = fmt.Sprintf("[TG] %s -> %s - Document - %s", username, replyToUsername, m.Document.FileName)
				}

				if m.Caption != "" {
					content = fmt.Sprintf("%s :\n%s", content, m.Caption)
				}

				mp := opqbot.SendMsgPack{
					SendType:   opqbot.SendTypeTextMsg,
					SendToType: opqbot.SendToTypeGroup,
					ToUserUid:  conf.ProConf.OPQBot.Group,
					Content:    opqbot.SendTypeTextMsgContent{Content: content},
				}

				OPQBot.Send(mp)
			}
			logs.Info("-> [TGbot] %+v", m.Chat)
		}
	})

	TGBot.Handle(tb.OnSticker, func(m *tb.Message) {
		if strconv.Itoa(int(m.Chat.ID)) == conf.ProConf.TGBot.ChatID {
			if m.Sender.ID != TGBot.Me.ID && arrays.ContainsString(conf.ProConf.TGBot.FilterID, m.Sender.Recipient()) == -1 {

				username := m.Sender.Username

				if username == "" {
					if m.Sender.LastName != "" {
						username = fmt.Sprintf("%s %s", m.Sender.FirstName, m.Sender.LastName)
					} else {
						username = fmt.Sprintf("%s", m.Sender.FirstName)
					}
				}

				content := fmt.Sprintf("[TG] %s - Sticker - %s", username, m.Sticker.Emoji)

				if m.IsReply() {
					replyToUsername := m.ReplyTo.Sender.Username

					if replyToUsername == "" {
						if m.Sender.LastName != "" {
							replyToUsername = fmt.Sprintf("%s %s", m.ReplyTo.Sender.FirstName, m.ReplyTo.Sender.LastName)
						} else {
							replyToUsername = fmt.Sprintf("%s", m.ReplyTo.Sender.FirstName)
						}
					}

					content = fmt.Sprintf("[TG] %s -> %s - Sticker - %s", username, replyToUsername, m.Sticker.Emoji)
				}

				if m.IsForwarded() {
					content = fmt.Sprintf("%s - Forwarded", content)
				}

				mp := opqbot.SendMsgPack{
					SendType:   opqbot.SendTypeTextMsg,
					SendToType: opqbot.SendToTypeGroup,
					ToUserUid:  conf.ProConf.OPQBot.Group,
					Content:    opqbot.SendTypeTextMsgContent{Content: content},
				}

				OPQBot.Send(mp)
			}
			logs.Info("-> [TGbot] %+v", m.Chat)
		}
	})

	TGBot.Handle(tb.OnAnimation, func(m *tb.Message) {
		if strconv.Itoa(int(m.Chat.ID)) == conf.ProConf.TGBot.ChatID {
			if m.Sender.ID != TGBot.Me.ID && arrays.ContainsString(conf.ProConf.TGBot.FilterID, m.Sender.Recipient()) == -1 {

				username := m.Sender.Username

				if username == "" {
					if m.Sender.LastName != "" {
						username = fmt.Sprintf("%s %s", m.Sender.FirstName, m.Sender.LastName)
					} else {
						username = fmt.Sprintf("%s", m.Sender.FirstName)
					}
				}

				content := fmt.Sprintf("[TG] %s - Animation", username)

				if m.IsReply() {

					replyToUsername := m.ReplyTo.Sender.Username

					if replyToUsername == "" {
						if m.Sender.LastName != "" {
							replyToUsername = fmt.Sprintf("%s %s", m.ReplyTo.Sender.FirstName, m.ReplyTo.Sender.LastName)
						} else {
							replyToUsername = fmt.Sprintf("%s", m.ReplyTo.Sender.FirstName)
						}
					}

					content = fmt.Sprintf("[TG] %s -> %s - Animation - %s", username, replyToUsername, m.Sticker.Emoji)
				}

				if m.IsForwarded() {
					content = fmt.Sprintf("%s - Forwarded", content)
				}

				mp := opqbot.SendMsgPack{
					SendType:   opqbot.SendTypeTextMsg,
					SendToType: opqbot.SendToTypeGroup,
					ToUserUid:  conf.ProConf.OPQBot.Group,
					Content:    opqbot.SendTypeTextMsgContent{Content: content},
				}

				OPQBot.Send(mp)

			}
			logs.Info("-> [TGbot] %+v", m.Chat)
		}
	})

	TGBot.Handle(tb.OnLocation, func(m *tb.Message) {
		if strconv.Itoa(int(m.Chat.ID)) == conf.ProConf.TGBot.ChatID {
			if m.Sender.ID != TGBot.Me.ID && arrays.ContainsString(conf.ProConf.TGBot.FilterID, m.Sender.Recipient()) == -1 {

				username := m.Sender.Username

				if username == "" {
					if m.Sender.LastName != "" {
						username = fmt.Sprintf("%s %s", m.Sender.FirstName, m.Sender.LastName)
					} else {
						username = fmt.Sprintf("%s", m.Sender.FirstName)
					}
				}

				content := fmt.Sprintf("[TG] %s - Location \n\t Lat: %f \n\t Lng: %f", username, m.Location.Lat, m.Location.Lng)

				if m.IsForwarded() {
					content = fmt.Sprintf("[TG] %s - Location - Forwarded \n\t Lng: %f \n\t Lat: %f", username, m.Location.Lat, m.Location.Lng)
				} else if m.IsReply() {
					replyToUsername := m.ReplyTo.Sender.Username

					if replyToUsername == "" {
						if m.Sender.LastName != "" {
							replyToUsername = fmt.Sprintf("%s %s", m.ReplyTo.Sender.FirstName, m.ReplyTo.Sender.LastName)
						} else {
							replyToUsername = fmt.Sprintf("%s", m.ReplyTo.Sender.FirstName)
						}
					}
					content = fmt.Sprintf("[TG] %s -> %s - Location \n\t Lat: %f \n\t Lng: %f", username, replyToUsername, m.Location.Lat, m.Location.Lng)
				}

				// if m.Caption != "" {
				// 	content = fmt.Sprintf("%s\n\n", content)
				// }

				mp := opqbot.SendMsgPack{
					SendType:   opqbot.SendTypeTextMsg,
					SendToType: opqbot.SendToTypeGroup,
					ToUserUid:  conf.ProConf.OPQBot.Group,
					Content:    opqbot.SendTypeTextMsgContent{Content: content},
				}

				OPQBot.Send(mp)

			}
			logs.Info("-> [TGbot] %+v", m.Chat)
		}
	})

	TGBot.Handle(tb.OnNewGroupPhoto, func(m *tb.Message) {
		if strconv.Itoa(int(m.Chat.ID)) == conf.ProConf.TGBot.ChatID {
			if m.Sender.ID != TGBot.Me.ID && arrays.ContainsString(conf.ProConf.TGBot.FilterID, m.Sender.Recipient()) == -1 {

				fileURL, err := TGBot.FileURLByID(m.NewGroupPhoto.FileID)
				if err != nil {
					logs.Error(err)
					return
				}

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
					Content:     "[TG] New Group Photo",
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
			logs.Info("-> [TGbot] %+v", m.Chat)
		}
	})

	TGBot.Handle(tb.OnNewGroupTitle, func(m *tb.Message) {
		if strconv.Itoa(int(m.Chat.ID)) == conf.ProConf.TGBot.ChatID {
			if m.Sender.ID != TGBot.Me.ID && arrays.ContainsString(conf.ProConf.TGBot.FilterID, m.Sender.Recipient()) == -1 {

				OPQBot.Send(opqbot.SendMsgPack{
					SendType:   opqbot.SendTypeTextMsg,
					SendToType: opqbot.SendToTypeGroup,
					ToUserUid:  conf.ProConf.OPQBot.Group,
					Content:    opqbot.SendTypeTextMsgContent{Content: fmt.Sprintf("[TG] New Group Title : %s ", m.NewGroupTitle)},
				})

			}
			logs.Info("-> [TGbot] %+v", m.Chat)
		}
	})

	TGBot.Start()
}
