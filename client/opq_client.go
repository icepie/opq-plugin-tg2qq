package client

import (
	"encoding/json"
	"fmt"
	"opq-plugin-tg2qq/client/opqbot"
	"opq-plugin-tg2qq/conf"
	"opq-plugin-tg2qq/util"
	"strings"

	"time"

	"github.com/wxnacy/wgo/arrays"

	"github.com/astaxie/beego/logs"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	OPQBot opqbot.BotManager
)

func OPQBotInit() {

	OPQBot = opqbot.NewBotManager(conf.ProConf.OPQBot.QQ, conf.ProConf.OPQBot.Url)

	err := OPQBot.Start()
	if err != nil {
		logs.Info("[OPQ] Fail to connet")
	}
	defer OPQBot.Stop()
	err = OPQBot.AddEvent(opqbot.EventNameOnGroupMessage, func(botQQ int64, packet opqbot.GroupMsgPack) {
		if packet.FromGroupID == conf.ProConf.OPQBot.Group && packet.FromUserID != OPQBot.QQ && arrays.ContainsInt(conf.ProConf.OPQBot.FilterQQ, packet.FromUserID) == -1 {
			if packet.MsgType == "TextMsg" {
				TGBot.Notify(MG, tb.Typing)
				TGBot.Send(MG, fmt.Sprintf("[QQ] %s : %s", packet.FromNickName, packet.Content))
			} else if packet.MsgType == "PicMsg" {
				TGBot.Notify(MG, tb.UploadingPhoto)

				ogc := opqbot.GroupContent{}
				err := json.Unmarshal([]byte(packet.Content), &ogc)
				if err != nil {
					fmt.Println(err)
				}

				for i := 1; i <= len(ogc.GroupPic); i++ {

					photo_caption := fmt.Sprintf("[QQ] %s ", packet.FromNickName)

					if len(ogc.GroupPic) > 1 {
						photo_caption = fmt.Sprintf("%s (%d/%d) ", photo_caption, i, len(ogc.GroupPic))
					}

					if ogc.Content != nil {
						photo_caption = fmt.Sprintf("%s : %s ", photo_caption, ogc.Content)
					}

					photo := &tb.Photo{File: tb.FromURL(ogc.GroupPic[i-1].Url), Caption: photo_caption}

					TGBot.Send(MG, photo)

				}
			} else if packet.MsgType == "AtMsg" {

				TGBot.Notify(MG, tb.Typing)

				gam := opqbot.GroupAtMsgContent{}
				err := json.Unmarshal([]byte(packet.Content), &gam)
				if err != nil {
					fmt.Println(err)
				}

				if gam.Tips == "[回复]" {
					user, err := OPQBot.GetUserInfo(gam.UserID[0])
					if err != nil {
						fmt.Println(err)
					}

					content := strings.TrimPrefix(gam.Content, fmt.Sprintf("@%s ", user.NickName))

					TGBot.Send(MG, fmt.Sprintf("[QQ] %s -> %s : %s", packet.FromNickName, user.NickName, content))
				} else {
					TGBot.Send(MG, fmt.Sprintf("[QQ] %s : %s", packet.FromNickName, gam.Content))
				}

			} else if packet.MsgType == "VoiceMsg" {
				TGBot.Notify(MG, tb.Typing)
				TGBot.Send(MG, fmt.Sprintf("[QQ] %s - Voice", packet.FromNickName))
			} else if packet.MsgType == "VideoMsg" {
				TGBot.Notify(MG, tb.Typing)

				// gvm := opqbot.VideoMsgContent{}
				// err := json.Unmarshal([]byte(packet.Content), &gvm)
				// if err != nil {
				// 	fmt.Println(err)
				// }

				// videodata := opqbot.VideoData{
				// 	GroupID:  packet.FromGroupID,
				// 	VideoUrl: gvm.VideoUrl,
				// 	VideoMd5: gvm.VideoMd5,
				// }

				// pb, _ := json.Marshal(videodata)

				// opqresp, err := http.Post(fmt.Sprintf("%s%s?qq=%d&funcname=PttCenterSvr.ShortVideoDownReq&timeout=10", conf.ProConf.OPQBot.Url, "/v1/LuaApiCaller", conf.ProConf.OPQBot.QQ),
				// 	"application/json",
				// 	bytes.NewBuffer(pb))
				// if err != nil {
				// 	fmt.Println(err)
				// }

				// b, err := ioutil.ReadAll(opqresp.Body)
				// if err != nil {
				// }

				// defer opqresp.Body.Close()

				// fmt.Println(string(b))

				// rte := opqbot.VideoDataRet{}
				// jsonErr := json.Unmarshal(b, &rte)
				// if jsonErr != nil {
				// 	log.Fatal(jsonErr)
				// }

				// vid := &tb.Video{File: tb.FromURL(rte.VideoUrl), Caption: fmt.Sprintf("[QQ] %s - Video ", packet.FromNickName)}

				// TGBot.Send(MG, vid)

				TGBot.Send(MG, fmt.Sprintf("[QQ] %s - Video ", packet.FromNickName))

			} else if packet.MsgType == "GroupFileMsg" {
				TGBot.Notify(MG, tb.Typing)

				gfm := opqbot.GroupFileMsgContent{}
				err := json.Unmarshal([]byte(packet.Content), &gfm)
				if err != nil {
					fmt.Println(err)
				}

				TGBot.Send(MG, fmt.Sprintf("[QQ] %s - File - %s - %s", packet.FromNickName, gfm.FileName, util.ByteSize(gfm.FileSize)))
			} else if packet.MsgType == "JsonMsg" {
				TGBot.Notify(MG, tb.Typing)

				TGBot.Send(MG, fmt.Sprintf("[QQ] %s - Json", packet.FromNickName))

			} else if packet.MsgType == "XmlMsg" {
				TGBot.Notify(MG, tb.Typing)

				TGBot.Send(MG, fmt.Sprintf("[QQ] %s - XML", packet.FromNickName))

				// Location := &tb.Location{Lat: 34.611679, Lng: 112.429459}
				// TGBot.Send(MG, Location)

			} else if packet.MsgType == "RedBagMsg" {
				TGBot.Notify(MG, tb.Typing)

				RedBaginfoStr := fmt.Sprint(packet.RedBaginfo)

				redBagInfoText := strings.TrimSuffix(strings.TrimPrefix(strings.Replace(RedBaginfoStr, ": ", ":", -1), "map["), "]")

				s := strings.Split(redBagInfoText, " ")

				redBagInfo := make(map[string]string)

				for _, v := range s {
					// fmt.Println(v)

					kv := strings.Split(v, ":")

					// i

					redBagInfo[kv[0]] = kv[1]
				}

				TGBot.Send(MG, fmt.Sprintf("[QQ] %s - RedBag - %s", packet.FromNickName, redBagInfo["Tittle"]))

			}

			logs.Info("-> [OPQ]%+v", packet)
		}
	})
	if err != nil {
		logs.Info(err.Error())
	}
	err = OPQBot.AddEvent(opqbot.EventNameOnFriendMessage, func(botQQ int64, packet opqbot.FriendMsgPack) {
		logs.Info("", botQQ, packet)
	})
	if err != nil {
		logs.Info(err.Error())
	}
	err = OPQBot.AddEvent(opqbot.EventNameOnGroupShut, func(botQQ int64, packet opqbot.GroupShutPack) {
		logs.Info("", botQQ, packet)
	})
	if err != nil {
		logs.Info(err.Error())
	}
	err = OPQBot.AddEvent(opqbot.EventNameOnConnected, func() {
		logs.Info("[OPQ] 连接成功！！！")

		info, err := OPQBot.GetUserInfo(OPQBot.QQ)
		if err != nil {
			logs.Error("[OPQ] Fail to get opq bot info from QQ %d", conf.ProConf.OPQBot.QQ)
		} else {
			logs.Info("[OPQ] Online: %+v", info)
			// OPQBot.Send(opqbot.SendMsgPack{
			// 	SendType:   opqbot.SendTypeTextMsg,
			// 	SendToType: opqbot.SendToTypeGroup,
			// 	ToUserUid:  conf.ProConf.OPQBot.Group,
			// 	Content:    opqbot.SendTypeTextMsgContent{Content: fmt.Sprintf("[opq-plugin-tg2qq] starting...\n\nTG Chat ID: %s", conf.ProConf.TGBot.ChatID)},
			// })
		}
	})
	if err != nil {
		logs.Error("[OPQ] 连接失败")
	}
	err = OPQBot.AddEvent(opqbot.EventNameOnDisconnected, func() {
		logs.Warn("[OPQ] 连接断开！！重新启动连接...")
		OPQBotInit()
	})
	if err != nil {
		logs.Info(err.Error())
	}
	err = OPQBot.AddEvent(opqbot.EventNameOnOther, func(botQQ int64, e interface{}) {
		logs.Info(e)
	})
	if err != nil {
		logs.Info(err.Error())
	}

	time.Sleep(525600 * time.Hour)

}
