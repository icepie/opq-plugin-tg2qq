package conf

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/astaxie/beego/logs"

	"github.com/spf13/viper"
)

// OPQBotConfig OPQ Bot 配置
type OPQBotConfig struct {
	Url      string
	QQ       int64
	FilterQQ []int64
	Group    int64
}

type ProxyConfig struct {
	Enable bool
	Url    string
}

// TGBotConfig Telegram Bot 配置
type TGBotConfig struct {
	Token    string
	ChatID   string
	FilterID []string
	Proxy    ProxyConfig
}

// Config 基础配置
type Config struct {
	OPQBot OPQBotConfig
	TGBot  TGBotConfig
}

// ProConf 新建实例
var ProConf = new(Config)

// init 初始化函数
func init() {

	// 取项目地址
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	var pathsep string
	if runtime.GOOS == "windows" {
		pathsep = "\\"
	} else {
		pathsep = "/"
	}

	cfpath := path + pathsep + "conf.yaml"

	viper.AddConfigPath(path)
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	if err != nil {
		logs.Trace("Please edit the " + cfpath + "，then restart app")
		os.Exit(1)
	}
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(ProConf); err != nil {
		logs.Error(err.Error())
	}

}
