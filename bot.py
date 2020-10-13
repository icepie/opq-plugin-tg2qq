#!/usr/bin/env python
# -*- coding: utf-8 -*-
import threading
import telebot
## tgbot debug mode
import logging

logger = telebot.logger
telebot.logger.setLevel(logging.DEBUG) # Outputs debug messages to console.

from botoy import Action, AsyncBotoy, Botoy, EventMsg, FriendMsg, GroupMsg, AsyncAction

qq = 111111111
qqbot = Botoy(qq=qq)
action = Action(qq)

telebot.apihelper.proxy = {'https':'socks5://127.0.0.1:1080'}
tgbot = telebot.TeleBot("11111111:xxxxxxxxxxxxxx")

qq_group_id = 1111111123
tg_chat_id = '-1111231111'

def tg_thread():
    tgbot.polling()

def qq_thread():
    thread = threading.Thread(target=qqbot.run)
    thread.start()

def all_thread():
    qq_thread()
    tg_thread()


@qqbot.on_group_msg
def group(ctx: GroupMsg):
    if ctx.FromGroupId == qq_group_id and ctx.FromUserId != qq:
        tgbot.send_message(tg_chat_id, "[QQ]" + " " + ctx.FromNickName + ": " + ctx.Content)
        
@qqbot.on_event
def event(ctx: EventMsg):
    pass

#@tgbot.message_handler(commands=['start', 'help'])
#def send_welcome(message):
#	tgbot.reply_to(message, "Howdy, how are you doing?")

@tgbot.message_handler(func=lambda message: True)
def echo_all(message):
    action.sendGroupText(qq_group_id, "[TG]" + " " + message.from_user.username + ": " + message.text)


@qqbot.when_connected
def _():
    print('OK!')

@qqbot.when_connected(every_time=True)
def _():
    print('OKK!!')


if __name__ == "__main__":
    all_thread()