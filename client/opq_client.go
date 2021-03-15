package client

import (
	"fmt"
	"opq-plugin-tg2qq/client/opqbot"
	"opq-plugin-tg2qq/conf"
	"opq-plugin-tg2qq/util/log"
	"strconv"
	"strings"
	"time"
)

var (
	ZanNote = map[int64]int{}
)

func OPQBotInit() {

	OPQBot := opqbot.NewBotManager(1962847213, conf.ProConf.OPQBot.Url)

	err := OPQBot.Start()
	if err != nil {
		log.OPQLog.Info("[OPQ] 连接失败")
	}
	defer OPQBot.Stop()
	err = OPQBot.AddEvent(opqbot.EventNameOnGroupMessage, func(botQQ int64, packet opqbot.GroupMsgPack) {
		if packet.FromUserID != OPQBot.QQ {
			if packet.Content == "赞我" {
				i, ok := ZanNote[packet.FromUserID]
				if ok {
					if i == time.Now().Day() {
						OPQBot.Send(opqbot.SendMsgPack{
							SendType:   opqbot.SendTypeTextMsg,
							SendToType: opqbot.SendToTypeGroup,
							ToUserUid:  packet.FromGroupID,
							Content:    opqbot.SendTypeTextMsgContent{Content: "今日已赞!"},
						})
						return
					}
				}
				OPQBot.Send(opqbot.SendMsgPack{
					SendType:   opqbot.SendTypeTextMsg,
					SendToType: opqbot.SendToTypeGroup,
					ToUserUid:  packet.FromGroupID,
					Content:    opqbot.SendTypeTextMsgContent{Content: "正在赞请稍后！"},
				})
				success := OPQBot.Zan(packet.FromUserID, 50)
				OPQBot.Send(opqbot.SendMsgPack{
					SendType:   opqbot.SendTypeTextMsg,
					SendToType: opqbot.SendToTypeGroup,
					ToUserUid:  packet.FromGroupID,
					Content:    opqbot.SendTypeTextMsgContent{Content: "成功赞了" + strconv.Itoa(success) + "次"},
				})
				ZanNote[packet.FromUserID] = time.Now().Day()
				return
			}
			if packet.Content == "刷新" && packet.FromUserID == 2435932516 {
				err := OPQBot.RefreshKey()
				if err != nil {
					log.OPQLog.Info(err.Error())
				}
			}
		}
		log.OPQLog.Info("%t", botQQ, packet)
	})
	if err != nil {
		log.OPQLog.Info(err.Error())
	}
	err = OPQBot.AddEvent(opqbot.EventNameOnFriendMessage, func(botQQ int64, packet opqbot.FriendMsgPack) {
		if packet.Content == "赞我" {
			i, ok := ZanNote[packet.FromUin]
			if ok {
				if i == time.Now().Day() {
					OPQBot.Send(opqbot.SendMsgPack{
						SendType:   opqbot.SendTypeTextMsg,
						SendToType: opqbot.SendToTypeFriend,
						ToUserUid:  packet.FromUin,
						Content:    opqbot.SendTypeTextMsgContent{Content: "今日已赞!"},
					})
					return
				}
			}
			OPQBot.Send(opqbot.SendMsgPack{
				SendType:   opqbot.SendTypeTextMsg,
				SendToType: opqbot.SendToTypeFriend,
				ToUserUid:  packet.FromUin,
				Content:    opqbot.SendTypeTextMsgContent{Content: "正在赞请稍后！"},
			})
			success := OPQBot.Zan(packet.FromUin, 50)
			OPQBot.Send(opqbot.SendMsgPack{
				SendType:   opqbot.SendTypeTextMsg,
				SendToType: opqbot.SendToTypeFriend,
				ToUserUid:  packet.FromUin,
				Content:    opqbot.SendTypeTextMsgContent{Content: "成功赞了" + strconv.Itoa(success) + "次"},
			})
			ZanNote[packet.FromUin] = time.Now().Day()
			return
		}
		if c := strings.Split(packet.Content, " "); len(c) >= 2 {
			if c[0] == "#查询" {
				log.OPQLog.Info(c[1])
				qq, err := strconv.ParseInt(c[1], 10, 64)
				if err != nil {
					OPQBot.Send(opqbot.SendMsgPack{
						SendType:   opqbot.SendTypeTextMsg,
						SendToType: opqbot.SendToTypeFriend,
						ToUserUid:  packet.FromUin,
						Content:    opqbot.SendTypeTextMsgContent{Content: err.Error()},
					})
				}
				user, err := OPQBot.GetUserInfo(qq)
				log.OPQLog.Info("", user)
				if err != nil {
					OPQBot.Send(opqbot.SendMsgPack{
						SendType:   opqbot.SendTypeTextMsg,
						SendToType: opqbot.SendToTypeFriend,
						ToUserUid:  packet.FromUin,
						Content:    opqbot.SendTypeTextMsgContent{Content: err.Error()},
					})
				} else {
					var sex string
					if user.Sex == 1 {
						sex = "女"
					} else {
						sex = "男"
					}
					OPQBot.Send(opqbot.SendMsgPack{
						SendType:   opqbot.SendTypeTextMsg,
						SendToType: opqbot.SendToTypeFriend,
						ToUserUid:  packet.FromUin,
						Content:    opqbot.SendTypeTextMsgContent{Content: fmt.Sprintf("用户:%s[%d]%s\n来自:%s\n年龄:%d\n被赞了:%d次\n", user.NickName, user.QQUin, sex, user.Province+user.City, user.Age, user.LikeNums)},
					})
				}
			}
		}
		log.OPQLog.Info("", botQQ, packet)
	})
	if err != nil {
		log.OPQLog.Info(err.Error())
	}
	err = OPQBot.AddEvent(opqbot.EventNameOnGroupShut, func(botQQ int64, packet opqbot.GroupShutPack) {
		log.OPQLog.Info("", botQQ, packet)
	})
	if err != nil {
		log.OPQLog.Info(err.Error())
	}
	err = OPQBot.AddEvent(opqbot.EventNameOnConnected, func() {
		log.OPQLog.Info("连接成功！！！")

		info, _ := OPQBot.GetUserInfo(OPQBot.QQ)

		log.OPQLog.Info("", info)
	})
	if err != nil {
		log.OPQLog.Info("连接失败")
	}
	err = OPQBot.AddEvent(opqbot.EventNameOnDisconnected, func() {
		log.OPQLog.Info("连接断开！！")
	})
	if err != nil {
		log.OPQLog.Info(err.Error())
	}
	err = OPQBot.AddEvent(opqbot.EventNameOnOther, func(botQQ int64, e interface{}) {
		log.OPQLog.Error(err.Error())
	})
	if err != nil {
		log.OPQLog.Info(err.Error())
	}
	time.Sleep(1 * time.Hour)
}
