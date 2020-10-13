
from botoy import Action, Botoy, EventMsg, FriendMsg, GroupMsg

qq = 88888888
bot = Botoy(qq=qq, use_plugins=True)
action = Action(qq)


@bot.on_friend_msg
def friend(ctx: FriendMsg):
    if ctx.Content == 'test':
        action.sendFriendText(ctx.FromUin, 'ok')


@bot.on_group_msg
def group(ctx: GroupMsg):
    if ctx.Content == 'test':
        action.sendGroupText(ctx.FromGroupId, 'ok')


@bot.on_event
def event(ctx: EventMsg):
    pass


if __name__ == "__main__":
    bot.run()
