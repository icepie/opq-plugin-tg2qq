#!/usr/bin/env python
# -*- coding: utf-8 -*-
import threading
import telebot
import urllib
import demjson
import requests
import base64

## tgbot debug mode
#import logging

#logger = telebot.logger
#telebot.logger.setLevel(logging.DEBUG) # Outputs debug messages to console.

from botoy import Action, AsyncBotoy, Botoy, EventMsg, FriendMsg, GroupMsg, AsyncAction

qq_num = 000000000
qqbot = Botoy(qq=qq_num)
action = Action(qq_num)

tg_token = "000000000:XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
tgbot = telebot.TeleBot(tg_token)

qq_group_id = 000000000
tg_chat_id = -000000000000

proxy_config = 'socks5://127.0.0.1:1080'

telebot.apihelper.proxy = {'https':proxy_config}

proxies = {
  'http': proxy_config,
  'https': proxy_config,
}


def tg_thread():
    try:
        tgbot.polling(none_stop=True)
    except:
        print("TG Bot connect fail")

def qq_thread():
    thread = threading.Thread(target=qqbot.run)
    thread.start()

def all_thread():
    qq_thread()
    tg_thread()

def base_64(path):
    try:
        with open(path, 'rb') as f:
            code = base64.b64encode(f.read()).decode()  # 读取文件内容，转换为base64编码
            return code
    except:
        pass
        return

@qqbot.on_group_msg
def group(ctx: GroupMsg):
    if ctx.FromGroupId == qq_group_id and ctx.FromUserId != qq_num:
        if ctx.MsgType == 'TextMsg':
            tgbot.send_chat_action(tg_chat_id, 'typing')
            tgbot.send_message(tg_chat_id, "[QQ]" + " " + ctx.FromNickName + ": " + ctx.Content)
        elif ctx.MsgType == 'PicMsg':
            ctxjs = demjson.decode(ctx.Content)
            pic_text = ""
            if ctxjs['Content']:
                pic_text = ": " + ctxjs['Content']
            tgbot.send_chat_action(tg_chat_id, 'upload_photo')
            tgbot.send_photo(tg_chat_id, ctxjs['GroupPic'][0]['Url'], caption = "[QQ]" + " " + ctx.FromNickName + pic_text)

@qqbot.on_event
def event(ctx: EventMsg):
    pass

#@tgbot.message_handler(commands=['start', 'help'])
#def send_welcome(message):
#	tgbot.reply_to(message, "Howdy, how are you doing?")

@tgbot.message_handler(content_types=["text"])
def echo_all(message):
    print(message)
    action.sendGroupText(qq_group_id, "[TG]" + " " + message.from_user.username + ": " + message.text)

@tgbot.message_handler(content_types=["photo"])
def tg_handle_photo(message):
    print(message)
    pic_text = ""
    if message.caption:
        pic_text = ": " + message.caption
    fileId=message.photo[0].file_id
    file_info = tgbot.get_file(fileId)
    pUrl = 'https://api.telegram.org/file/bot{0}/{1}'.format(tg_token, file_info.file_path)
    #print(pUrl)
    picfile = requests.get(pUrl, proxies=proxies)
    picfile_name = ".cache/" + fileId + ".png"
    open(picfile_name, 'wb').write(picfile.content)
    action.sendGroupPic(qq_group_id, picBase64Buf=base_64(picfile_name), content= "[TG]" + " " + message.from_user.username + pic_text, atUser= 0)

@qqbot.when_connected()
def _():
    print('OK!')

@qqbot.when_connected(every_time=True)
def _():
    print('OKK!!')


if __name__ == "__main__":
    all_thread()