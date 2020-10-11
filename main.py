#!/usr/bin/env python
# -*- coding: utf-8 -*-

import telebot

telebot.apihelper.proxy = {'https':'socks5://127.0.0.1:1080'}

bot = telebot.TeleBot("TONKEN")

@bot.message_handler(commands=['start', 'help'])
def send_welcome(message):
	bot.reply_to(message, "Howdy, how are you doing?")

@bot.message_handler(func=lambda message: True)
def echo_all(message):
	bot.reply_to(message, message.text)

bot.polling()