package client

import (
	"encoding/json"
	"fmt"
	"opq-plugin-tg2qq/client/opqbot"
	"opq-plugin-tg2qq/conf"
	"time"

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
		if packet.FromGroupID == conf.ProConf.OPQBot.Group && packet.FromUserID != OPQBot.QQ {
			if packet.MsgType == "TextMsg" {
				TGBot.Notify(MG, tb.Typing)
				TGBot.Send(MG, fmt.Sprintf("[QQ] [%s] %s", packet.FromNickName, packet.Content))
			} else if packet.MsgType == "PicMsg" {
				TGBot.Notify(MG, tb.UploadingPhoto)

				ogc := opqbot.GroupContent{}
				err := json.Unmarshal([]byte(packet.Content), &ogc)
				if err != nil {
					fmt.Println(err)
				}

				for i := 1; i <= len(ogc.GroupPic); i++ {

					photo_caption := fmt.Sprintf("[QQ] [%s]", packet.FromNickName)

					if len(ogc.GroupPic) > 1 {
						photo_caption = fmt.Sprintf("%s (%d/%d) ", photo_caption, i, len(ogc.GroupPic))
					}

					if ogc.Content != nil {
						photo_caption = fmt.Sprintf("%s %s ", photo_caption, ogc.Content)
					}

					photo := &tb.Photo{File: tb.FromURL(ogc.GroupPic[i-1].Url), Caption: photo_caption}

					TGBot.Send(MG, photo)

				}

			}
			logs.Info("-> [OPQ]%+v", packet)
		}
	})
	if err != nil {
		logs.Info(err.Error())
	}
	err = OPQBot.AddEvent(opqbot.EventNameOnFriendMessage, func(botQQ int64, packet opqbot.FriendMsgPack) {
		// if packet.Content == "赞我" {
		// 	i, ok := ZanNote[packet.FromUin]
		// 	if ok {
		// 		if i == time.Now().Day() {
		// 			OPQBot.Send(opqbot.SendMsgPack{
		// 				SendType:   opqbot.SendTypeTextMsg,
		// 				SendToType: opqbot.SendToTypeFriend,
		// 				ToUserUid:  packet.FromUin,
		// 				Content:    opqbot.SendTypeTextMsgContent{Content: "今日已赞!"},
		// 			})
		// 			return
		// 		}
		// 	}
		// 	OPQBot.Send(opqbot.SendMsgPack{
		// 		SendType:   opqbot.SendTypeTextMsg,
		// 		SendToType: opqbot.SendToTypeFriend,
		// 		ToUserUid:  packet.FromUin,
		// 		Content:    opqbot.SendTypeTextMsgContent{Content: "正在赞请稍后！"},
		// 	})
		// 	success := OPQBot.Zan(packet.FromUin, 50)
		// 	OPQBot.Send(opqbot.SendMsgPack{
		// 		SendType:   opqbot.SendTypeTextMsg,
		// 		SendToType: opqbot.SendToTypeFriend,
		// 		ToUserUid:  packet.FromUin,
		// 		Content:    opqbot.SendTypeTextMsgContent{Content: "成功赞了" + strconv.Itoa(success) + "次"},
		// 	})
		// 	logs.Info(packet.FromUin)
		// 	ZanNote[packet.FromUin] = time.Now().Day()
		// 	return
		// }
		// if c := strings.Split(packet.Content, " "); len(c) >= 2 {
		// 	if c[0] == "#查询" {
		// 		logs.Info(c[1])
		// 		qq, err := strconv.ParseInt(c[1], 10, 64)
		// 		if err != nil {
		// 			OPQBot.Send(opqbot.SendMsgPack{
		// 				SendType:   opqbot.SendTypeTextMsg,
		// 				SendToType: opqbot.SendToTypeFriend,
		// 				ToUserUid:  packet.FromUin,
		// 				Content:    opqbot.SendTypeTextMsgContent{Content: err.Error()},
		// 			})
		// 		}
		// 		user, err := OPQBot.GetUserInfo(qq)
		// 		logs.Info("", user)
		// 		if err != nil {
		// 			OPQBot.Send(opqbot.SendMsgPack{
		// 				SendType:   opqbot.SendTypeTextMsg,
		// 				SendToType: opqbot.SendToTypeFriend,
		// 				ToUserUid:  packet.FromUin,
		// 				Content:    opqbot.SendTypeTextMsgContent{Content: err.Error()},
		// 			})
		// 		} else {
		// 			var sex string
		// 			if user.Sex == 1 {
		// 				sex = "女"
		// 			} else {
		// 				sex = "男"
		// 			}
		// 			OPQBot.Send(opqbot.SendMsgPack{
		// 				SendType:   opqbot.SendTypeTextMsg,
		// 				SendToType: opqbot.SendToTypeFriend,
		// 				ToUserUid:  packet.FromUin,
		// 				Content:    opqbot.SendTypeTextMsgContent{Content: fmt.Sprintf("用户:%s[%d]%s\n来自:%s\n年龄:%d\n被赞了:%d次\n", user.NickName, user.QQUin, sex, user.Province+user.City, user.Age, user.LikeNums)},
		// 			})
		// 		}
		// 	}
		// }
		// logs.Info("", botQQ, packet)
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
			OPQBot.Send(opqbot.SendMsgPack{
				SendType:   opqbot.SendTypeTextMsg,
				SendToType: opqbot.SendToTypeGroup,
				ToUserUid:  conf.ProConf.OPQBot.Group,
				Content:    opqbot.SendTypeTextMsgContent{Content: fmt.Sprintf("[opq-plugin-tg2qq] starting...\n\nTG Chat ID: %s", conf.ProConf.TGBot.ChatID)},
			})
		}
	})
	if err != nil {
		logs.Error("[OPQ] 连接失败")
	}
	err = OPQBot.AddEvent(opqbot.EventNameOnDisconnected, func() {
		logs.Warn("[OPQ] 连接断开！！")
	})
	if err != nil {
		logs.Info(err.Error())
	}
	err = OPQBot.AddEvent(opqbot.EventNameOnOther, func(botQQ int64, e interface{}) {
		logs.Error(err.Error())
	})
	if err != nil {
		logs.Info(err.Error())
	}
	time.Sleep(1 * time.Hour)
}
