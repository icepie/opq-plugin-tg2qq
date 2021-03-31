# opq-plugin-tg2qq

> 当 `TG` 与 `QQ` 碰撞会产生什么样的火花呢?

OPQ Bot 💗 Telegram Bot

## 使用

### 获取

```bash
$ git clone https://github.com/icepie/opq-plugin-tg2qq
$ cd opq-plugin-tg2qq
$ go build
```

### 配置

```yaml
opqbot:
  url: http://127.0.0.1:8888
  qq: 0
  group: 0
  filterqq:
    - 0
    - 0
tgbot:
  token: ""
  chatid: ""
  filterid:
    - ""
    - ""
  proxy:
    enable: true
    url: sock5://127.0.0.1:1080
```

### 运行

```bash
./opq-plugin-tg2qq
```

## 说明


## 最后

> 主要使用到的库如下

- [telebot](https://github.com/tucnak/telebot)

- [OPQBot](https://github.com/mcoo/OPQBot)

