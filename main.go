package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"opq-plugin-tg2qq/client/opq"

	"opq-plugin-tg2qq/client"

	"github.com/asmcos/requests"
)

var ZanNote = map[int64]int{}

func main() {
	err := client.OPQBot.Start()
	if err != nil {
		log.Println(err.Error())
	}
	defer client.OPQBot.Stop()
	err = client.OPQBot.AddEvent(opq.EventNameOnGroupMessage, func(botQQ int64, packet opq.GroupMsgPack) {
		if packet.FromUserID != client.OPQBot.QQ {
			if packet.Content == "赞我" {
				i, ok := ZanNote[packet.FromUserID]
				if ok {
					if i == time.Now().Day() {
						client.OPQBot.Send(opq.SendMsgPack{
							SendType:   opq.SendTypeTextMsg,
							SendToType: opq.SendToTypeGroup,
							ToUserUid:  packet.FromGroupID,
							Content:    opq.SendTypeTextMsgContent{Content: "今日已赞!"},
						})
						return
					}
				}
				client.OPQBot.Send(opq.SendMsgPack{
					SendType:   opq.SendTypeTextMsg,
					SendToType: opq.SendToTypeGroup,
					ToUserUid:  packet.FromGroupID,
					Content:    opq.SendTypeTextMsgContent{Content: "正在赞请稍后！"},
				})
				success := client.OPQBot.Zan(packet.FromUserID, 50)
				client.OPQBot.Send(opq.SendMsgPack{
					SendType:   opq.SendTypeTextMsg,
					SendToType: opq.SendToTypeGroup,
					ToUserUid:  packet.FromGroupID,
					Content:    opq.SendTypeTextMsgContent{Content: "成功赞了" + strconv.Itoa(success) + "次"},
				})
				ZanNote[packet.FromUserID] = time.Now().Day()
				return
			}
			if packet.Content == "二次元图片" {
				res, err := requests.Get("http://www.dmoe.cc/random.php?return=json")
				if err != nil {
					return
				}
				var pixivPic Pic
				_ = res.Json(&pixivPic)
				client.OPQBot.Send(opq.SendMsgPack{
					SendType:   opq.SendTypePicMsgByUrl,
					SendToType: opq.SendToTypeGroup,
					ToUserUid:  int64(packet.FromGroupID),
					Content:    opq.SendTypePicMsgByUrlContent{Content: "随机", PicUrl: pixivPic.Imgurl},
				})
				return
			}
			if packet.Content == "刷新" && packet.FromUserID == 2435932516 {
				err := client.OPQBot.RefreshKey()
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
		log.Println(botQQ, packet)
	})
	if err != nil {
		log.Println(err.Error())
	}
	err = client.OPQBot.AddEvent(opq.EventNameOnFriendMessage, func(botQQ int64, packet opq.FriendMsgPack) {
		if packet.Content == "赞我" {
			i, ok := ZanNote[packet.FromUin]
			if ok {
				if i == time.Now().Day() {
					client.OPQBot.Send(opq.SendMsgPack{
						SendType:   opq.SendTypeTextMsg,
						SendToType: opq.SendToTypeFriend,
						ToUserUid:  packet.FromUin,
						Content:    opq.SendTypeTextMsgContent{Content: "今日已赞!"},
					})
					return
				}
			}
			client.OPQBot.Send(opq.SendMsgPack{
				SendType:   opq.SendTypeTextMsg,
				SendToType: opq.SendToTypeFriend,
				ToUserUid:  packet.FromUin,
				Content:    opq.SendTypeTextMsgContent{Content: "正在赞请稍后！"},
			})
			success := client.OPQBot.Zan(packet.FromUin, 50)
			client.OPQBot.Send(opq.SendMsgPack{
				SendType:   opq.SendTypeTextMsg,
				SendToType: opq.SendToTypeFriend,
				ToUserUid:  packet.FromUin,
				Content:    opq.SendTypeTextMsgContent{Content: "成功赞了" + strconv.Itoa(success) + "次"},
			})
			ZanNote[packet.FromUin] = time.Now().Day()
			return
		}
		if c := strings.Split(packet.Content, " "); len(c) >= 2 {
			if c[0] == "#查询" {
				log.Println(c[1])
				qq, err := strconv.ParseInt(c[1], 10, 64)
				if err != nil {
					client.OPQBot.Send(opq.SendMsgPack{
						SendType:   opq.SendTypeTextMsg,
						SendToType: opq.SendToTypeFriend,
						ToUserUid:  packet.FromUin,
						Content:    opq.SendTypeTextMsgContent{Content: err.Error()},
					})
				}
				user, err := client.OPQBot.GetUserInfo(qq)
				log.Println(user)
				if err != nil {
					client.OPQBot.Send(opq.SendMsgPack{
						SendType:   opq.SendTypeTextMsg,
						SendToType: opq.SendToTypeFriend,
						ToUserUid:  packet.FromUin,
						Content:    opq.SendTypeTextMsgContent{Content: err.Error()},
					})
				} else {
					var sex string
					if user.Sex == 1 {
						sex = "女"
					} else {
						sex = "男"
					}
					client.OPQBot.Send(opq.SendMsgPack{
						SendType:   opq.SendTypeTextMsg,
						SendToType: opq.SendToTypeFriend,
						ToUserUid:  packet.FromUin,
						Content:    opq.SendTypeTextMsgContent{Content: fmt.Sprintf("用户:%s[%d]%s\n来自:%s\n年龄:%d\n被赞了:%d次\n", user.NickName, user.QQUin, sex, user.Province+user.City, user.Age, user.LikeNums)},
					})
				}
			}
		}
		log.Println(botQQ, packet)
	})
	if err != nil {
		log.Println(err.Error())
	}
	err = client.OPQBot.AddEvent(opq.EventNameOnGroupShut, func(botQQ int64, packet opq.GroupShutPack) {
		log.Println(botQQ, packet)
	})
	if err != nil {
		log.Println(err.Error())
	}
	err = client.OPQBot.AddEvent(opq.EventNameOnConnected, func() {
		log.Println("连接成功！！！")
	})
	if err != nil {
		log.Println(err.Error())
	}
	err = client.OPQBot.AddEvent(opq.EventNameOnDisconnected, func() {
		log.Println("连接断开！！")
	})
	if err != nil {
		log.Println(err.Error())
	}
	err = client.OPQBot.AddEvent(opq.EventNameOnOther, func(botQQ int64, e interface{}) {
		log.Println(e)
	})
	if err != nil {
		log.Println(err.Error())
	}
	//client.OPQBot.Send(opq.SendMsgPack{
	//	SendType:   opq.SendTypePicMsgByUrl,
	//	SendToType: opq.SendToTypeFriend,
	//	ToUserUid:  2435932516,
	//	Content:    opq.SendTypePicMsgByUrlContent{Content: "你好", PicUrl: "https://img-home.csdnimg.cn/images/20201124032511.png"},
	//})
	time.Sleep(1 * time.Hour)
}

type Pic struct {
	Code   string `json:"code"`
	Imgurl string `json:"imgurl"`
	Width  string `json:"width"`
	Height string `json:"height"`
}
